package controllers

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/peternabil/go-api/models"
)

func (server Server) Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		user := models.User{}
		if c.Request.Header.Get("Authorization") == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "you must be logged in to perform this request"})
			c.Abort()
			return
		}
		reqToken := c.Request.Header.Get("Authorization")
		splitToken := strings.Split(reqToken, "Bearer ")
		if len(splitToken) == 1 {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "you must be logged in to perform this request"})
			c.Abort()
			return
		}
		reqToken = splitToken[1]
		claims := &models.Claims{}
		token, err := jwt.ParseWithClaims(reqToken, claims, func(token *jwt.Token) (any, error) {
			return []byte(os.Getenv("JWT_SECRET_KEY")), nil
		})
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "you must be logged in to perform this request"})
			c.Abort()
			return
		}
		email := claims.Email
		user.Email = email
		if usError := server.store.FindUser(email, &user); usError != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "no user with this email"})
			c.Abort()
			return
		}
		c.Set("user", user)
	}
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
