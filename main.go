package main

import (
	"github.com/gin-gonic/gin"
	"github.com/peternabil/go-api/controllers"
	"github.com/peternabil/go-api/intitializers"
)

func init() {
	intitializers.LoadEnvVariables()
	intitializers.LoadDB()
}

func main() {
	r := gin.Default()

	r.POST("/signup", controllers.SignUp)
	r.GET("/users", controllers.UserIndex)
	r.GET("/users/:id", controllers.UserFind)

	r.GET("/transaction", controllers.TransactionIndex)
	r.GET("/transaction/:id", controllers.TransactionFind)
	r.POST("/transaction", controllers.TransactionCreate)
	r.PUT("/transaction/:id", controllers.TransactionEdit)
	r.DELETE("/transaction/:id", controllers.TransactionDelete)

	r.GET("/category", controllers.CategoryIndex)
	r.GET("/category/:id", controllers.CategoryFind)
	r.POST("/category", controllers.CategoryCreate)
	r.PUT("/category/:id", controllers.CategoryEdit)
	r.DELETE("/category", controllers.CategoryDelete)

	r.GET("/priority", controllers.PriorityIndex)
	r.GET("/priority/:id", controllers.PriorityFind)
	r.POST("/priority", controllers.PriorityCreate)
	r.PUT("/priority/:id", controllers.PriorityEdit)
	r.DELETE("/priority/:id", controllers.PriorityDelete)

	r.Run() // listen and serve on 0.0.0.0:8080
}
