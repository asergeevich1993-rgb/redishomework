package server

import (
	"context"
	"net/http"
	"redis/handlers"
)

type HTTPServer struct {
	handlers *handlers.HttpHandlers
	svr      *http.Server
}

func NewServer(hh *handlers.HttpHandlers) *HTTPServer {
	return &HTTPServer{
		handlers: hh,
	}
}

func (hs *HTTPServer) StartServer() {
	rounter := http.NewServeMux()
	rounter.HandleFunc("POST /book", hs.handlers.HandleCreateBook)
	rounter.HandleFunc("GET /book/{id}", hs.handlers.HandleGetBook)
	rounter.HandleFunc("POST /finish", hs.handlers.HandleFinis)

	hs.svr = &http.Server{
		Addr:    ":8080",
		Handler: rounter,
	}
	hs.svr.ListenAndServe()
}
func (hs *HTTPServer) FinishServer(ctx context.Context) error {
	return hs.svr.Shutdown(ctx)
}
