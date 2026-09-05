package handlers

import (
	"errors"
	"strings"
)

type DTOBook struct {
	Title  string `json:"title"`
	Author string `json:"author"`
}

func (d *DTOBook) ValidateForCreate() error {
	d.Title = strings.TrimSpace(d.Title)
	d.Author = strings.TrimSpace(d.Author)

	if d.Title == "" {
		return errors.New("title is empty")
	}
	if d.Author == "" {
		return errors.New("author is empty")
	}
	return nil
}
