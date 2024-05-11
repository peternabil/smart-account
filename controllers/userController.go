package controllers

import (
	"errors"
	"fmt"
	"net/http"

	emailverifier "github.com/AfterShip/email-verifier"
	"github.com/gin-gonic/gin"
	"github.com/go-passwd/validator"
	"github.com/google/uuid"
	"github.com/peternabil/go-api/models"
	"golang.org/x/crypto/bcrypt"
)

func (server *Server) UserIndex(c *gin.Context) {
	users := []models.User{}
	result := server.store.GetUsers(&users)
	if result != nil {
		c.Status(400)
		return
	}
	c.JSON(200, gin.H{
		"users": users,
	})
}

func (server *Server) UserFind(c *gin.Context) {
	uId := c.Param("id")
	ussId, uuidErr := uuid.Parse(uId)
	if uuidErr != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "invalid uuid"})
		return
	}
	user := models.User{}
	if res := server.store.GetUser(ussId, &user); res != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}
	c.JSON(200, gin.H{
		"user": user,
	})
}

func (server *Server) SignUp(c *gin.Context) {
	var body struct {
		Email     string
		FirstName string
		LastName  string
		Password  string
	}
	reqErr := c.BindJSON(&body)
	if reqErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": reqErr.Error()})
		return
	}
	passwordValid := validator.New(validator.MinLength(6, errors.New("password must be at least 6 chars")), validator.MaxLength(30, errors.New("password must be at most 30 chars")), validator.CommonPassword(errors.New("password cannot be commonly used password")), validator.ContainsAtLeast("abcdefghijklmnopqrstuvwxyz", 5, errors.New("password must contain at least 5 chars")), validator.ContainsAtLeast("_@.()@$#", 1, errors.New("password must contain at least 1 special char _@.()@$#")))
	err := passwordValid.Validate(body.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	verifier := emailverifier.NewVerifier()
	ret, err := verifier.Verify(body.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if !ret.Syntax.Valid {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid email syntax"})
		return
	}
	encryptedPass, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to encrypt password for some reason"})
		return
	}
	fmt.Println(encryptedPass)
	user := models.User{Email: body.Email, Password: string(encryptedPass), FirstName: body.FirstName, LastName: body.LastName}
	result := server.store.SignUp(&user)
	if result != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
	c.JSON(200, gin.H{
		"user": user,
	})
}

func (server *Server) Login(c *gin.Context) {
	var body struct {
		Email    string
		Password string
	}
	reqErr := c.BindJSON(&body)
	fmt.Println(body)
	if reqErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": reqErr.Error()})
		return
	}
	user := models.User{Email: body.Email}
	fmt.Println(user)
	if usError := server.store.FindUser(body.Email, &user); usError != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "email or password is incorrect"})
		return
	}
	fmt.Println(user)
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "email or password is incorrect"})
		return
	}
	token, tokErr := server.store.CreateToken(user)
	if tokErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": tokErr.Error()})
		return
	}
	c.JSON(200, gin.H{"token": token, "user": user})
}
