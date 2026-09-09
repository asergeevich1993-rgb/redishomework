package main

import (
	"context"
	"os"
	"redis/handlers"
	"redis/server"
	"redis/storage"
	"time"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	shutdownctx, shutdowncancel := context.WithTimeout(ctx, 5*time.Second)
	defer shutdowncancel()
	dsn := os.Getenv("CONN_STRING")
	stor, err := storage.NewConnectDB(ctx, dsn, "REDIS_STRING")
	if err != nil {
		panic(err)
	}
	stor.CreateIndex(ctx)
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
