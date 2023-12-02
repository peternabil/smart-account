package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/peternabil/go-api/intitializers"
	"github.com/peternabil/go-api/models"
)

func TransactionIndex(c *gin.Context) {
	transactions := []models.Transaction{}
	result := intitializers.DB.Find(&transactions)
	if result.Error != nil {
		c.Status(400)
		return
	}
	c.JSON(200, gin.H{
		"transactions": transactions,
	})
}

func TransactionFind(c *gin.Context) {
	tId := c.Param("id")
	transaction := models.Transaction{ID: uuid.MustParse(tId)}
	res := intitializers.DB.First(&transaction)
	if res.Error != nil {
		c.Status(404)
		return
	}
	c.JSON(200, gin.H{
		"transaction": transaction,
	})
}

func TransactionCreate(c *gin.Context) {
	var body struct {
		Title       string
		Category    string
		Amount      int
		Negative    bool
		Description string
		Priority    string
	}
	c.BindJSON(&body)
	cat := models.Category{ID: uuid.MustParse(body.Category)}
	prio := models.Priority{ID: uuid.MustParse(body.Category)}
	category := intitializers.DB.First(&cat)
	if category.Error != nil {
		c.Status(400)
		return
	}
	priority := intitializers.DB.First(&prio)
	if priority.Error != nil {
		c.Status(400)
		return
	}
	transaction := models.Transaction{Title: body.Title, Category: cat, Priority: prio, Amount: body.Amount, Negative: body.Negative, Description: body.Description}
	result := intitializers.DB.Create(&transaction)
	if result.Error != nil {
		c.Status(400)
		return
	}
	c.JSON(200, gin.H{
		"transaction": transaction,
	})
}

func TransactionEdit(c *gin.Context) {
	var body struct {
		Title       string
		Category    string
		Amount      int
		Negative    bool
		Description string
		Priority    string
	}
	tId := c.Param("id")
	c.BindJSON(&body)
	transaction := models.Transaction{ID: uuid.MustParse(tId)}
	res := intitializers.DB.First(&transaction)
	if res.Error != nil {
		c.Status(404)
		return
	}
	cat := models.Category{ID: uuid.MustParse(body.Category)}
	prio := models.Priority{ID: uuid.MustParse(body.Category)}
	category := intitializers.DB.First(&cat)
	if category.Error != nil {
		c.Status(400)
		return
	}
	priority := intitializers.DB.First(&prio)
	if priority.Error != nil {
		c.Status(400)
		return
	}
	transaction.Title = body.Title
	transaction.Category = cat
	transaction.Priority = prio
	transaction.Amount = body.Amount
	transaction.Negative = body.Negative
	transaction.Description = body.Description
	intitializers.DB.Save(&transaction)
	c.JSON(200, gin.H{
		"transaction": transaction,
	})
}

func TransactionDelete(c *gin.Context) {
	tId := c.Param("id")
	transaction := models.Transaction{ID: uuid.MustParse(tId)}
	res := intitializers.DB.Delete(&transaction)
	if res.Error != nil {
		c.Status(400)
		return
	}
	c.JSON(200, gin.H{
		"transaction": transaction,
	})
}
