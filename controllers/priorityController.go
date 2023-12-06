package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/peternabil/go-api/intitializers"
	"github.com/peternabil/go-api/models"
)

func PriorityIndex(c *gin.Context) {
	priorities := []models.Priority{}
	user := c.MustGet("user").(models.User)
	if result := intitializers.DB.Where("user_id = ?", user.UID).Find(&priorities).Error; result != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "no priorities for this user"})
		return
	}
	c.JSON(200, gin.H{
		"priorities": priorities,
	})
}

func PriorityFind(c *gin.Context) {
	pId := c.Param("id")
	priority := models.Priority{ID: uuid.MustParse(pId)}
	res := intitializers.DB.First(&priority)
	if res.Error != nil {
		c.Status(404)
		return
	}
	c.JSON(200, gin.H{
		"priority": priority,
	})
}

func PriorityCreate(c *gin.Context) {
	var body struct {
		Name        string
		Description string
		Level       int
	}
	user := c.MustGet("user").(models.User)
	c.BindJSON(&body)
	priority := models.Priority{Name: body.Name, Description: body.Description, UserID: user.UID, Level: body.Level}
	result := intitializers.DB.Create(&priority)
	if result.Error != nil {
		c.Status(400)
		return
	}
	c.JSON(200, gin.H{
		"priority": priority,
	})
}

func PriorityEdit(c *gin.Context) {
	var body struct {
		Name        string
		Description string
		Level       int
	}
	pId := c.Param("id")
	c.BindJSON(&body)
	prio := models.Priority{ID: uuid.MustParse(pId)}
	res := intitializers.DB.Find(&prio)
	if res.Error != nil {
		c.Status(404)
		return
	}
	prio.Name = body.Name
	prio.Description = body.Description
	prio.Level = body.Level
	intitializers.DB.Save(&prio)
	c.JSON(200, gin.H{
		"priority": prio,
	})
}

func PriorityDelete(c *gin.Context) {
	pId := c.Param("id")
	priority := models.Priority{ID: uuid.MustParse(pId)}
	res := intitializers.DB.Delete(&priority)
	if res.Error != nil {
		c.Status(400)
		return
	}
	c.JSON(200, gin.H{
		"priority": priority,
	})
}
