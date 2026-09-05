package storage

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/redis/go-redis/v9"
)

type Storage struct {
	db  *pgx.Conn
	rdb *redis.Client
}

func NewConnectDB(ctx context.Context, dsn string, redisAddr string) (*Storage, error) {
	conn, err := pgx.Connect(ctx, dsn)
	if err != nil {
		return nil, err
	}
	rdb := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		PoolSize: 20})

	return &Storage{
		db:  conn,
		rdb: rdb,
	}, nil
}

func (s *Storage) CreateTable(ctx context.Context) error {
	sql := `CREATE TABLE IF NOT EXISTS libradis(
 id SERIAL PRIMARY KEY,
 title VARCHAR(200) NOT NULL,
 author VARCHAR(100) NOT NULL)`

	tag, err := s.db.Exec(ctx, sql)
	if err != nil {
		panic(err)
	}
	fmt.Println("База создана : ", tag)
	return nil
}
func (s *Storage) InsertBook(ctx context.Context, title, author string) (int, error) {
	var id int
	sql := `INSERT INTO libradis (title,author) VALUES ($1,$2) RETURNING id`
	err := s.db.QueryRow(ctx, sql, title, author).Scan(&id)
	if err != nil {
		return 0, err
	}
	strid := strconv.Itoa(id)
	err = s.rdb.HSet(ctx, strid, map[string]interface{}{
		"author": author,
		"title":  title,
	}).Err()
	if err != nil {
		return 0, err
	}
	s.rdb.Expire(ctx, strid, 5*time.Second)
	return id, nil
}

func (s *Storage) GetBook(ctx context.Context, id int) (Book, string, error) {

	strid := strconv.Itoa(id)
	data, err := s.rdb.HGetAll(ctx, strid).Result()
	if err == nil {
		title := data["title"]
		author := data["author"]
		return Book{
			ID:     id,
			Title:  title,
			Author: author,
		}, " Redis ", nil
	}
	if err != redis.Nil {
		return Book{}, "", err
	}
	var book Book
	sql := `SELECT id,title,author FROM libradis WHERE id = $1`
	err = s.db.QueryRow(ctx, sql, id).Scan(&book.ID, &book.Title, &book.Author)
	if err != nil {
		return Book{}, "", err
	}
	return book, "postgres", nil

}
