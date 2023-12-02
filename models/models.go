package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Category struct {
	gorm.Model
	ID          uuid.UUID `gorm:"type:uuid;default:gen_random_uuid()"`
	Name        string
	Description string
}

type Priority struct {
	gorm.Model
	ID          uuid.UUID `gorm:"type:uuid;default:gen_random_uuid()"`
	Name        string
	Description string
	Level       int
}

type Transaction struct {
	gorm.Model
	ID          uuid.UUID `gorm:"type:uuid;default:gen_random_uuid()"`
	Title       string
	Category    Category
	Amount      int
	Negative    bool
	Description string
	Priority    Priority
}
