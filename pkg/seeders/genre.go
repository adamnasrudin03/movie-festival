package seeders

import (
	"adamnasrudin03/movie-festival/app/entity"

	"gorm.io/gorm"
)

func InitGenre(db *gorm.DB) {
	tx := db.Begin()
	var genres = []entity.Genre{}
	tx.Select("id").Find(&genres)
	if len(genres) == 0 {
		genres = []entity.Genre{
			{
				Name: "Action",
			},
			{
				Name: "Adventure",
			},
			{
				Name: "Comedy",
			},
			{
				Name: "Crime",
			},
			{
				Name: "Drama",
			},
			{
				Name: "Fantasy",
			},
			{
				Name: "Horror",
			},
			{
				Name: "Mystery",
			},
			{
				Name: "Romance",
			},
			{
				Name: "Sci-fi",
			},
			{
				Name: "Thriller",
			},
		}
		tx.Create(genres)
	}

	tx.Commit()
}
