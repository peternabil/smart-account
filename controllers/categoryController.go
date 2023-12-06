package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/peternabil/go-api/intitializers"
	"github.com/peternabil/go-api/models"
)

func CategoryIndex(c *gin.Context) {
	var categories []models.Category
	user := c.MustGet("user").(models.User)
	if result := intitializers.DB.Where("user_id = ?", user.UID).Find(&categories).Error; result != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "no categories for this user"})
		return
	}
	c.JSON(200, gin.H{
		"categories": categories,
	})
}

func CategoryFind(c *gin.Context) {
	cId := c.Param("id")
	category := models.Category{ID: uuid.MustParse(cId)}
	res := intitializers.DB.First(&category)
	if res.Error != nil {
		c.Status(404)
		return
	}
	c.JSON(200, gin.H{
		"category": category,
	})
}

func CategoryCreate(c *gin.Context) {
	var body struct {
		Name        string
		Description string
	}
	user := c.MustGet("user").(models.User)
	c.BindJSON(&body)
	category := models.Category{Name: body.Name, Description: body.Description, UserID: user.UID}
	result := intitializers.DB.Create(&category)
	if result.Error != nil {
		c.Status(400)
		return
	}
	c.JSON(200, gin.H{
		"category": category,
	})
}

func CategoryEdit(c *gin.Context) {
	var body struct {
		Name        string
		Description string
	}
	catId := c.Param("id")
	c.BindJSON(&body)
	cat := models.Category{ID: uuid.MustParse(catId)}
	res := intitializers.DB.Find(&cat)
	if res.Error != nil {
		c.Status(404)
		return
	}
	cat.Name = body.Name
	cat.Description = body.Description
	intitializers.DB.Save(&cat)
	c.JSON(200, gin.H{
		"category": cat,
	})
}

func CategoryDelete(c *gin.Context) {
	cId := c.Param("id")
	category := models.Category{ID: uuid.MustParse(cId)}
	res := intitializers.DB.Delete(&category)
	if res.Error != nil {
		c.Status(400)
		return
	}
	c.JSON(200, gin.H{
		"category": category,
	})
}
