package entity

import (
	"adamnasrudin03/movie-festival/pkg/helpers"
	"errors"

	"gorm.io/gorm"
)

// User represents the model for an user
type User struct {
	ID       uint64 `gorm:"primaryKey" json:"id"`
	Name     string `gorm:"not null" json:"name" `
	Role     string `gorm:"not null;default:'USER'" json:"role"`
	Email    string `gorm:"not null;uniqueIndex" json:"email" `
	Password string `gorm:"not null" json:"password,omitempty"`
	GORMModel
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	hashedPass, err := helpers.HashPassword(u.Password)
	if err != nil {
		return
	}
	u.Password = hashedPass
	if (u.Role != "ADMIN" && u.Role != "USER") && u.Role != "" {
		err = errors.New("role is invalid, must be 'ADMIN' or 'USER'")
		return
	}

	return
}
