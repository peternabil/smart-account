package controllers

import (
	"time"

	"github.com/gin-contrib/cors"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/penglongli/gin-metrics/ginmetrics"
	"github.com/peternabil/go-api/store"
	"go.uber.org/zap"
)

type Server struct {
	store  store.Store
	router *gin.Engine
	mw     gin.HandlerFunc // Optional middleware for testing
}

func NewServer(store store.Store, mw gin.HandlerFunc) (*Server, error) {
	server := &Server{
		store: store,
		mw:    mw,
	}
	server.NewRouter()
	return server, nil
}

func (server *Server) NewRouter() {

	r := gin.Default()

	SetupCORS(r)
	SetupLogger(r)
	SetupMetrics(r)

	// auth not required
	nonAuth := r.Group("/smart-account/api")
	nonAuth.POST("/auth/signup", server.SignUp)
	nonAuth.POST("/auth/login", server.Login)
	// auth required
	auth := nonAuth.Group("/v1")
	if server.mw != nil {
		auth.Use(server.mw)
	} else {
		auth.Use(server.Auth())
	}

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
	auth.DELETE("/category/:id", server.CategoryDelete)

	auth.GET("/priority", server.PriorityIndex)
	auth.GET("/priority/:id", server.PriorityFind)
	auth.POST("/priority", server.PriorityCreate)
	auth.PUT("/priority/:id", server.PriorityEdit)
	auth.DELETE("/priority/:id", server.PriorityDelete)

	auth.GET("/daily", server.GetDailyValues)
	auth.GET("/highest-cat", server.GetHighestCategory)
	auth.GET("/highest-prio", server.GetHighestPriority)

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

func SetupCORS(router *gin.Engine) {
	config := cors.DefaultConfig()
	// config.AllowOrigins = []string{"http://google.com"}
	// config.AllowOrigins = []string{"http://google.com", "http://facebook.com"}
	config.AllowOrigins = []string{"*"}
	// config.AllowAllOrigins = true

	// router.Use(cors.New(config))
	router.Use(CORSMiddleware())
}

func SetupMetrics(router *gin.Engine) {
	m := ginmetrics.GetMonitor()
	m.SetMetricPath("/smart-account/api/metrics")
	m.SetSlowTime(10)
	m.SetDuration([]float64{0.1, 0.3, 1.2, 5, 10})
	m.Use(router)
}
