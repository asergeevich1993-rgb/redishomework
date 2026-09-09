package storage

import (
	"context"
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

func (s *Storage) CreateIndex(ctx context.Context) error {

	indexSQL := `CREATE INDEX IF NOT EXISTS idx_books_author ON libradis(author)`
	if _, err := s.db.Exec(ctx, indexSQL); err != nil {
		return err
	}

	return nil
}
func (s *Storage) GetIndex(ctx context.Context) ([]string, error) {
	sql := `SELECT indexname FROM pg_indexes WHERE tablename='libradis'`
	rows, err := s.db.Query(ctx, sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var indexes []string

	for rows.Next() {
		var index string
		if err := rows.Scan(&index); err != nil {
			return nil, err
		}
		indexes = append(indexes, index)

	}
	return indexes, rows.Err()

}
func (s *Storage) InsertBook(ctx context.Context, title, author string) (int, error) {
	var id int
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return 0, err
	}

	sql := `INSERT INTO libradis (title,author) VALUES ($1,$2) RETURNING id`
	err = tx.QueryRow(ctx, sql, title, author).Scan(&id)
	if err != nil {
		tx.Rollback(ctx)
		return 0, err
	}
	err = tx.Commit(ctx)
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
	s.rdb.Expire(ctx, strid, 20*time.Second)
	return id, nil
}

func (s *Storage) GetBook(ctx context.Context, id int) (Book, string, error) {

	strid := strconv.Itoa(id)
	data, err := s.rdb.HGetAll(ctx, strid).Result()
	if err != nil {
		return Book{}, "", err
	}
	if len(data) > 0 {
		return Book{
			ID:     id,
			Title:  data["title"],
			Author: data["author"],
		}, "redis", nil
	}

	var book Book
	sql := `SELECT id,title,author FROM libradis WHERE id = $1`
	err = s.db.QueryRow(ctx, sql, id).Scan(&book.ID, &book.Title, &book.Author)
	if err != nil {
		return Book{}, "", err
	}
	return book, "postgres", nil

}
func (s *Storage) UpdateBook(ctx context.Context, id int, title, author string) (int, error) {

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return 0, err
	}
	var existid int
	sql := `SELECT id FROM libradis WHERE id=$1 FOR UPDATE`
	err = tx.QueryRow(ctx, sql, id).Scan(&existid)
	if err != nil {
		tx.Rollback(ctx)
		return 0, err
	}
	var intid int
	sql2 := `UPDATE libradis SET title=$1,author=$2 WHERE id=$3 RETURNING id`
	err = tx.QueryRow(ctx, sql2, title, author, id).Scan(&intid)
	if err != nil {
		tx.Rollback(ctx)
		return 0, err
	}
	return intid, tx.Commit(ctx)

}
