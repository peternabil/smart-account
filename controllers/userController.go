package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/peternabil/go-api/intitializers"
	"github.com/peternabil/go-api/models"
	"golang.org/x/crypto/bcrypt"
)

func UserIndex(c *gin.Context) {
	users := []models.User{}
	result := intitializers.DB.Find(&users)
	if result.Error != nil {
		c.Status(400)
		return
	}
	c.JSON(200, gin.H{
		"users": users,
	})
}

func UserFind(c *gin.Context) {
	uId := c.Param("id")
	user := models.User{UID: uuid.MustParse(uId)}
	res := intitializers.DB.First(&user)
	if res.Error != nil {
		c.Status(404)
		return
	}
	c.JSON(200, gin.H{
		"user": user,
	})
}

func SignUp(c *gin.Context) {
	var body struct {
		Email    string
		Password string
	}
	reqErr := c.BindJSON(&body)
	if reqErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": reqErr.Error()})
		return
	}
	encryptedPass, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to encrypt password for some reason"})
		return
	}
	fmt.Println(encryptedPass)
	user := models.User{Email: body.Email, Password: string(encryptedPass)}
	result := intitializers.DB.Create(&user)
	if result.Error != nil {
		c.Status(400)
		return
	}
	c.JSON(200, gin.H{
		"user": user,
	})
}
