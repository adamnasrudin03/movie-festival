package seeders

import (
	"adamnasrudin03/movie-festival/app/entity"

	"gorm.io/gorm"
)

func InitUser(db *gorm.DB) {
	tx := db.Begin()
	var users = []entity.User{}
	tx.Select("id").Where("role = ? ", "ADMIN").Find(&users)
	if len(users) == 0 {
		user := entity.User{
			Name:     "Admin",
			Password: "password123",
			Email:    "admin@gmail.com",
			Role:     "ADMIN",
		}
		tx.Create(&user)
	}

	tx.Commit()
}
