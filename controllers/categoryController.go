package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/peternabil/go-api/models"
)

func (server *Server) CategoryIndex(c *gin.Context) {
	categories := []models.Category{}
	user := c.MustGet("user").(models.User)
	if result := server.store.GetCategories(user.UID, &categories); result != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "no categories for this user"})
		return
	}
	c.JSON(200, gin.H{
		"categories": categories,
	})
}

func (server *Server) CategoryFind(c *gin.Context) {
	cId := c.Param("id")
	category := models.Category{ID: uuid.MustParse(cId)}
	res := server.store.GetCategory(category.ID, &category)
	if res != nil {
		c.Status(404)
		return
	}
	c.JSON(200, gin.H{
		"category": category,
	})
}

func (server *Server) CategoryCreate(c *gin.Context) {
	var body struct {
		Name        string
		Description string
	}
	user := c.MustGet("user").(models.User)
	c.BindJSON(&body)
	category := models.Category{Name: body.Name, Description: body.Description, UserID: user.UID}
	result := server.store.CreateCategory(&category)
	if result != nil {
		c.Status(400)
		return
	}
	c.JSON(200, gin.H{
		"category": category,
	})
}

func (server *Server) CategoryEdit(c *gin.Context) {
	var body struct {
		Name        string
		Description string
	}
	catId := c.Param("id")
	c.BindJSON(&body)
	cat := models.Category{ID: uuid.MustParse(catId)}
	res := server.store.GetCategory(cat.ID, &cat)
	if res != nil {
		c.Status(404)
		return
	}
	cat.Name = body.Name
	cat.Description = body.Description
	server.store.EditCategory(&cat)
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
