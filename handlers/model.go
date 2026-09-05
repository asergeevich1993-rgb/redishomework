package handlers

import "redis/storage"

type Model struct {
	Book storage.Book `json:"book"`
	Str  string       `json:"str"`
}
