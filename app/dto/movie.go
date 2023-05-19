package dto

import (
	"time"
)

type CreateUpdateMovie struct {
	Title        string                    `json:"title" validate:"required"`
	Description  string                    `json:"description" validate:"required"`
	Duration     uint64                    `json:"duration" validate:"required"`
	DurationType string                    `json:"duration_type"`
	WatchUrl     string                    `json:"watch_url" validate:"required"`
	Artists      string                    `json:"artists" validate:"required"`
	Genres       []CreateUpdateGenreMovies `json:"genres"`
}
type CreateUpdateGenreMovies struct {
	ID uint64 ` json:"id"  validate:"required,numeric,min=1"`
}

type MovieRes struct {
	ID           uint64     `json:"id"`
	Title        string     `json:"title" `
	Description  string     `json:"description"`
	Duration     uint64     `json:"duration"`
	DurationType string     `json:"duration_type"`
	WatchUrl     string     `json:"watch_url"`
	Artists      string     `json:"artists"`
	Viewers      uint64     `json:"viewers"`
	Genres       string     `json:"genres"`
	GenreDetails []GenreRes `json:"genre_detail,omitempty"`
	CreatedAt    time.Time  `json:"created_at,omitempty"`
	UpdatedAt    time.Time  `json:"updated_at,omitempty"`
}

type GenreRes struct {
	ID   uint64 `json:"id,omitempty"`
	Name string `json:"name,omitempty" `
}

func (MovieRes) TableName() string {
	return "movies"
}

type FileRes struct {
	Path string `json:"path" `
}
