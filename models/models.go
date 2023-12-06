package models

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	UID       uuid.UUID `gorm:"type:uuid;default:gen_random_uuid()"`
	Email     string    `gorm:"unique"`
	FirstName string
	LastName  string
	Password  string
}

type Category struct {
	gorm.Model
	ID          uuid.UUID `gorm:"type:uuid;default:gen_random_uuid()"`
	Name        string
	Description string
	UserID      uuid.UUID
}

type Priority struct {
	gorm.Model
	ID          uuid.UUID `gorm:"type:uuid;default:gen_random_uuid()"`
	Name        string
	Description string
	Level       int
	UserID      uuid.UUID
}

type Transaction struct {
	gorm.Model
	ID          uuid.UUID `gorm:"type:uuid;default:gen_random_uuid()"`
	Title       string
	CategoryID  uuid.UUID
	Amount      int
	Negative    bool
	Description string
	PriorityID  uuid.UUID
	UserID      uuid.UUID
}

type Claims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}
