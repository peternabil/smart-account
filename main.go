package main

import (
	"github.com/gin-gonic/gin"
	"github.com/peternabil/go-api/controllers"
	"github.com/peternabil/go-api/intitializers"
	"github.com/peternabil/go-api/middleware"
)

func init() {
	intitializers.LoadEnvVariables()
	intitializers.LoadDB()
}

func main() {
	r := gin.Default()

	// auth not required
	nonAuth := r.Group("/")

	nonAuth.POST("/signup", controllers.SignUp)
	nonAuth.POST("/login", controllers.Login)

	// auth required
	auth := r.Group("/api", middleware.Auth())

	auth.GET("/users", controllers.UserIndex)
	auth.GET("/users/:id", controllers.UserFind)

	auth.GET("/transaction", controllers.TransactionIndex)
	auth.GET("/transaction/:id", controllers.TransactionFind)
	auth.POST("/transaction", controllers.TransactionCreate)
	auth.PUT("/transaction/:id", controllers.TransactionEdit)
	auth.DELETE("/transaction/:id", controllers.TransactionDelete)

	auth.GET("/category", controllers.CategoryIndex)
	auth.GET("/category/:id", controllers.CategoryFind)
	auth.POST("/category", controllers.CategoryCreate)
	auth.PUT("/category/:id", controllers.CategoryEdit)
	auth.DELETE("/category", controllers.CategoryDelete)

	auth.GET("/priority", controllers.PriorityIndex)
	auth.GET("/priority/:id", controllers.PriorityFind)
	auth.POST("/priority", controllers.PriorityCreate)
	auth.PUT("/priority/:id", controllers.PriorityEdit)
	auth.DELETE("/priority/:id", controllers.PriorityDelete)

	r.Run() // listen and serve on 0.0.0.0:8080
}
