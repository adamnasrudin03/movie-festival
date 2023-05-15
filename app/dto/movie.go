package dto

type CreateUpdateMovie struct {
	Title       string `json:"title" validate:"required"`
	Description string `json:"description" validate:"required"`
	Duration    uint64 `json:"duration" validate:"required"`
	WatchUrl    string `json:"watch_url" validate:"required"`
	Artists     string `json:"artists" validate:"required"`
	Genres      string `json:"genres" validate:"required"`
}
