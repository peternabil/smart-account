package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/peternabil/go-api/models"
)

func (server *Server) PriorityIndex(c *gin.Context) {
	user := server.store.GetUserFromToken(c)
	priorities, err := server.store.GetPriorities(user.UID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "no priorities for this user"})
		return
	}
	c.JSON(200, gin.H{
		"priorities": priorities,
	})
}

func (server *Server) PriorityFind(c *gin.Context) {
	pId := c.Param("id")
	user := server.store.GetUserFromToken(c)
	priority := models.Priority{ID: uuid.MustParse(pId)}
	res, err := server.store.GetPriority(user.UID, &priority)
	if err != nil {
		c.Status(404)
		return
	}
	c.JSON(200, gin.H{
		"priority": res,
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
	result, err := server.store.CreatePriority(&priority)
	if err != nil {
		c.Status(400)
		return
	}
	c.JSON(200, gin.H{
		"priority": result,
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
	user := server.store.GetUserFromToken(c)
	prio := models.Priority{ID: uuid.MustParse(pId)}
	res, err := server.store.GetPriority(user.UID, &prio)
	if err != nil {
		c.Status(404)
		return
	}
	prio = res
	prio.Name = body.Name
	prio.Description = body.Description
	prio.Level = body.Level
	res, err = server.store.EditPriority(&prio)
	if err != nil {
		c.Status(500)
		return
	}
	c.JSON(200, gin.H{
		"priority": res,
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
