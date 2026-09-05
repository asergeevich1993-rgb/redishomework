package main

import (
	"context"
	"redis/handlers"
	"redis/server"
	"redis/storage"
	"time"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	shutdownctx, shutdowncancel := context.WithTimeout(ctx, 5*time.Second)
	defer shutdowncancel()
	dsn := "postgres://postgres:Svs1512!@localhost:5432/testdb"
	stor, err := storage.NewConnectDB(ctx, dsn, "localhost:6379")
	if err != nil {
		panic(err)
	}
	stor.CreateTable(ctx)
	handler := handlers.NewHandler(stor, cancel, ctx)
	svr := server.NewServer(handler)

	go func() {
		<-ctx.Done()
		if err := svr.FinishServer(shutdownctx); err != nil {
			panic(err)
		}

	}()

	svr.StartServer()

}
