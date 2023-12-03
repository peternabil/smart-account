package main

import (
	"fmt"

	"github.com/peternabil/go-api/intitializers"
	"github.com/peternabil/go-api/models"
)

func init() {
	intitializers.LoadEnvVariables()
	intitializers.LoadDB()
}

func main() {
	var err error
	err = intitializers.DB.AutoMigrate(&models.Category{})
	if err != nil {
		fmt.Println(err.Error())
	}
	err = intitializers.DB.AutoMigrate(&models.Priority{})
	if err != nil {
		fmt.Println(err.Error())
	}
	err = intitializers.DB.AutoMigrate(&models.User{})
	if err != nil {
		fmt.Println(err.Error())
	}
	err = intitializers.DB.AutoMigrate(&models.Transaction{})
	if err != nil {
		fmt.Println(err.Error())
	}
}
