package main

import (
	"fmt"
	"os"

	"github.com/lpernett/godotenv"
	"github.com/peternabil/go-api/controllers"
	"github.com/peternabil/go-api/intitializers"
	"github.com/peternabil/go-api/models"
)

func init() {
	intitializers.LoadEnvVariables()
	intitializers.LoadDB()
}

func migrate() {
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

func main() {
	godotenv.Load()
	var store intitializers.MainStore
	server, _ := controllers.NewServer(store)
	intitializers.LoadDB()
	migrate()
	server.Start(fmt.Sprintf(
		"%s:%s",
		os.Getenv("API_ADDRESS"),
		os.Getenv("API_PORT"),
	))
}
