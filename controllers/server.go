package controllers

import (
	"time"

	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/penglongli/gin-metrics/ginmetrics"
	"github.com/peternabil/go-api/store"
	"go.uber.org/zap"
)

type Server struct {
	store  store.Store
	router *gin.Engine
}

func NewServer(store store.Store) (*Server, error) {
	server := &Server{
		store: store,
	}
	server.NewRouter()
	return server, nil
}

func (server *Server) NewRouter() {

	r := gin.Default()

	SetupLogger(r)
	SetupMetrics(r)

	// auth not required
	nonAuth := r.Group("/")
	nonAuth.POST("/signup", server.SignUp)
	nonAuth.POST("/login", server.Login)
	// auth required
	auth := r.Group("/api", server.Auth())
	auth.GET("/users", server.UserIndex)
	auth.GET("/users/:id", server.UserFind)

	auth.GET("/transaction", server.TransactionIndex)
	auth.GET("/transaction/:id", server.TransactionFind)
	auth.POST("/transaction", server.TransactionCreate)
	auth.PUT("/transaction/:id", server.TransactionEdit)
	auth.DELETE("/transaction/:id", server.TransactionDelete)

	auth.GET("/category", server.CategoryIndex)
	auth.GET("/category/:id", server.CategoryFind)
	auth.POST("/category", server.CategoryCreate)
	auth.PUT("/category/:id", server.CategoryEdit)
	auth.DELETE("/category", server.CategoryDelete)

	auth.GET("/priority", server.PriorityIndex)
	auth.GET("/priority/:id", server.PriorityFind)
	auth.POST("/priority", server.PriorityCreate)
	auth.PUT("/priority/:id", server.PriorityEdit)
	auth.DELETE("/priority/:id", server.PriorityDelete)

	server.router = r
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func SetupLogger(router *gin.Engine) {
	logger, _ := zap.NewProduction()
	router.Use(ginzap.Ginzap(logger, time.RFC3339, true))
	router.Use(ginzap.RecoveryWithZap(logger, true))

}

func SetupMetrics(router *gin.Engine) {
	m := ginmetrics.GetMonitor()
	m.SetMetricPath("/metrics")
	m.SetSlowTime(10)
	m.SetDuration([]float64{0.1, 0.3, 1.2, 5, 10})
	m.Use(router)
}
