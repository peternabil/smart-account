package controllers

import (
	"net/http"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/lpernett/godotenv"
	"github.com/peternabil/go-api/store"
	"github.com/stretchr/testify/require"
)

func newTestServer(t *testing.T, store store.Store) *Server {
	server, err := NewServer(store, nil)
	require.NoError(t, err)

	server1 := &http.Server{
		Addr:    ":8080",
		Handler: server.router,
	}

	go server.Start(":8080")
	err = server1.Close()
	require.NoError(t, err)

	return server
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	godotenv.Load()
	os.Exit(m.Run())
}
