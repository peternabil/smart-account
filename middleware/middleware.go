package middleware

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/peternabil/go-api/intitializers"
	"github.com/peternabil/go-api/models"
)

func Auth() gin.HandlerFunc {
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
		if usError := intitializers.DB.Where("email = ?", email).First(&user).Error; usError != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "no user with this email"})
			c.Abort()
			return
		}
		err = intitializers.DB.Find(&user).Error
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "you must be logged in to perform this request"})
			c.Abort()
			return
		}
		c.Set("user", user)
	}
}
