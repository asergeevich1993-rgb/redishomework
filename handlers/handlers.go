package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"redis/storage"
	"strconv"
)

type HttpHandlers struct {
	storage *storage.Storage
	cancel  context.CancelFunc
	ctx     context.Context
}

func NewHandler(stor *storage.Storage, cancel context.CancelFunc, ctx context.Context) *HttpHandlers {
	return &HttpHandlers{
		storage: stor,
		cancel:  cancel,
		ctx:     ctx,
	}
}
func (hh *HttpHandlers) HandleGetIndex(w http.ResponseWriter, r *http.Request) {
	indexes, err := hh.storage.GetIndex(hh.ctx)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(indexes)

}

func (hh *HttpHandlers) HandleCreateBook(w http.ResponseWriter, r *http.Request) {
	var dto DTOBook
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := dto.ValidateForCreate(); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	id, err := hh.storage.InsertBook(hh.ctx, dto.Title, dto.Author)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	json.NewEncoder(w).Encode(map[string]int{"Created id: ": id})
}
func (hh *HttpHandlers) HandleGetBook(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	intid, err := strconv.Atoi(id)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	books, strg, err := hh.storage.GetBook(hh.ctx, intid)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(Model{
		Book: books,
		Str:  strg,
	})
}
func (hh *HttpHandlers) HandleFinis(w http.ResponseWriter, r *http.Request) {
	hh.cancel()
	w.WriteHeader(200)
}

func writeError(w http.ResponseWriter, status int, err string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(err)

}
