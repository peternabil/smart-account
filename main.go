package main

import (
	"fmt"
	"os"

	"github.com/lpernett/godotenv"
	"github.com/peternabil/go-api/controllers"
	"github.com/peternabil/go-api/intitializers"
)

func init() {
	intitializers.LoadEnvVariables()
	intitializers.LoadDB()
}

func main() {
	godotenv.Load()
	var store intitializers.MainStore
	server, _ := controllers.NewServer(store)
	intitializers.LoadDB()
	server.Start(fmt.Sprintf(
		"%s:%s",
		os.Getenv("API_ADDRESS"),
		os.Getenv("API_PORT"),
	))
}
