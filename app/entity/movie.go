package entity

// Movie represents the model for an Movie
type Movie struct {
	ID           uint64 `gorm:"primaryKey" json:"id"`
	Title        string `gorm:"not null" json:"title" `
	Description  string `gorm:"not null" json:"description"`
	Duration     uint64 `gorm:"not null" json:"duration"`
	DurationType string `gorm:"not null;default:'minutes'" json:"duration_type"`
	WatchUrl     string `gorm:"not null" json:"watch_url"`
	Artists      string `gorm:"not null" json:"artists"`
	Genres       string `gorm:"not null" json:"genres"`
	Viewers      uint64 `gorm:"not null;default:0" json:"viewers"`
	GORMModel
}
