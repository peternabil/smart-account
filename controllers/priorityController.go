package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/peternabil/go-api/models"
)

func (server *Server) PriorityIndex(c *gin.Context) {
	priorities := []models.Priority{}
	user := server.store.GetUserFromToken(c)
	if result := server.store.GetPriorities(user.UID, &priorities); result != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "no priorities for this user"})
		return
	}
	c.JSON(200, gin.H{
		"priorities": priorities,
	})
}

func (server *Server) PriorityFind(c *gin.Context) {
	pId := c.Param("id")
	priority := models.Priority{ID: uuid.MustParse(pId)}
	res := server.store.GetPriority(priority.ID, &priority)
	if res != nil {
		c.Status(404)
		return
	}
	c.JSON(200, gin.H{
		"priority": priority,
	})
}

func (server *Server) PriorityCreate(c *gin.Context) {
	var body struct {
		Name        string
		Description string
		Level       int
	}
	user := server.store.GetUserFromToken(c)
	err := c.BindJSON(&body)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	priority := models.Priority{Name: body.Name, Description: body.Description, UserID: user.UID, Level: body.Level}
	result := server.store.CreatePriority(&priority)
	if result != nil {
		c.Status(400)
		return
	}
	c.JSON(200, gin.H{
		"priority": priority,
	})
}

func (server *Server) PriorityEdit(c *gin.Context) {
	var body struct {
		Name        string
		Description string
		Level       int
	}
	pId := c.Param("id")
	err := c.BindJSON(&body)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	prio := models.Priority{ID: uuid.MustParse(pId)}
	res := server.store.GetPriority(prio.ID, &prio)
	if res != nil {
		c.Status(404)
		return
	}
	prio.Name = body.Name
	prio.Description = body.Description
	prio.Level = body.Level
	server.store.EditPriority(&prio)
	c.JSON(200, gin.H{
		"priority": prio,
	})
}

func (server *Server) PriorityDelete(c *gin.Context) {
	pId := c.Param("id")
	priority := models.Priority{ID: uuid.MustParse(pId)}
	res := server.store.DeletePriority(&priority)
	if res != nil {
		c.Status(400)
		return
	}
	c.JSON(200, gin.H{
		"priority": priority,
	})
}
