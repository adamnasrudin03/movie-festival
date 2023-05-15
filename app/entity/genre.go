package entity

// Genre represents the model for an Genre
type Genre struct {
	ID   uint64 `gorm:"primaryKey" json:"id"`
	Name string `gorm:"not null;uniqueIndex" json:"name" `
	GORMModel
}
