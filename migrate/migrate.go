package main

import (
	"github.com/peternabil/go-api/intitializers"
	"github.com/peternabil/go-api/models"
)

func init() {
	intitializers.LoadEnvVariables()
	intitializers.LoadDB()
}

func main() {
	intitializers.DB.AutoMigrate(&models.Category{})
	intitializers.DB.AutoMigrate(&models.Transaction{})
	intitializers.DB.AutoMigrate(&models.Priority{})
}
