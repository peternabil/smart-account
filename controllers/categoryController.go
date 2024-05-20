package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/peternabil/go-api/models"
)

func (server *Server) CategoryIndex(c *gin.Context) {
	user := server.store.GetUserFromToken(c)
	cats, err := server.store.GetCategories(user.UID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "no categories for this user"})
		return
	}
	c.JSON(200, gin.H{
		"categories": cats,
	})
}

func (server *Server) CategoryFind(c *gin.Context) {
	user := server.store.GetUserFromToken(c)
	cId := c.Param("id")
	category := models.Category{ID: uuid.MustParse(cId)}
	cat, err := server.store.GetCategory(user.UID, &category)
	if err != nil {
		c.Status(404)
		return
	}
	c.JSON(200, gin.H{
		"category": cat,
	})
}

func (server *Server) CategoryCreate(c *gin.Context) {
	var body struct {
		Name        string
		Description string
	}
	user := server.store.GetUserFromToken(c)
	err := c.BindJSON(&body)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	category := models.Category{Name: body.Name, Description: body.Description, UserID: user.UID}
	cat, err := server.store.CreateCategory(&category)
	if err != nil {
		c.Status(400)
		return
	}
	c.JSON(200, gin.H{
		"category": cat,
	})
}

func (server *Server) CategoryEdit(c *gin.Context) {
	var body struct {
		Name        string
		Description string
	}
	catId := c.Param("id")
	err := c.BindJSON(&body)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	user := server.store.GetUserFromToken(c)
	cat := models.Category{ID: uuid.MustParse(catId)}
	res, err := server.store.GetCategory(user.UID, &cat)
	if err != nil {
		c.Status(404)
		return
	}
	cat = res
	cat.Name = body.Name
	cat.Description = body.Description
	cat, err = server.store.EditCategory(&cat)
	if err != nil {
		c.Status(500)
		return
	}
	c.JSON(200, gin.H{
		"category": cat,
	})
}

func (server *Server) CategoryDelete(c *gin.Context) {
	cId := c.Param("id")
	category := models.Category{ID: uuid.MustParse(cId)}
	res := server.store.DeleteCategory(&category)
	if res != nil {
		c.Status(400)
		return
	}
	c.JSON(200, gin.H{
		"category": category,
	})
}
