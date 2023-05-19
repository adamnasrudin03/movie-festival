package entity

// Genre represents the model for an Genre
type Genre struct {
	ID   uint64 `gorm:"primaryKey" json:"id"`
	Name string `gorm:"not null;uniqueIndex" json:"name" `
	GORMModel
}

type GenreMovies struct {
	GenreID uint64 `gorm:"primaryKey" json:"genre_id"`
	MovieID uint64 `gorm:"primaryKey" json:"movie_id"`
	GORMModel
}

func (GenreMovies) TableName() string {
	return "genre_movies"
}
