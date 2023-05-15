package dto

type CreateUpdateGenre struct {
	Name string `json:"name" validate:"required"`
}
