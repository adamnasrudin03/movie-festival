package entity

// Log represents the model for an log
type Log struct {
	ID          uint64       `gorm:"primaryKey" json:"id"`
	Name        string       `gorm:"not null" json:"name"`
	Action      string       `gorm:"not null" json:"action"`
	TableNameID uint64       `json:"table_name_id"`
	TableName   string       `json:"table_name"`
	UserID      uint64       `json:"user_id"`
	User        UserRelation `gorm:"ForeignKey:UserID" json:"user"`
	GORMModel
}

type UserRelation struct {
	ID    uint64 `json:"id"`
	Name  string `json:"name"`
	Role  string `json:"role"`
	Email string `json:"email"`
}

func (UserRelation) TableName() string {
	return "users"
}
