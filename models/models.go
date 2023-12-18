package models

import (
	"time"

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
	Category    Category
	Priority    Priority
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

type SpendingCategory struct {
	Category   Category
	CategoryID uuid.UUID
	Date       time.Time
	Total      int64
	Negative   bool
}

type SpendingPriority struct {
	Priority   Priority
	PriorityID uuid.UUID
	Date       time.Time
	Total      int64
	Negative   bool
}

type Spending struct {
	Date     time.Time
	Total    int64
	Negative bool
}
