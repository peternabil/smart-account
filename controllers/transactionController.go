package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/peternabil/go-api/intitializers"
	"github.com/peternabil/go-api/models"
)

func TransactionIndex(c *gin.Context) {
	transactions := []models.Transaction{}
	user := c.MustGet("user").(models.User)
	if result := intitializers.DB.Where("user_id = ?", user.UID).Find(&transactions).Error; result != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "no transactions for this user"})
		return
	}
	c.JSON(200, gin.H{
		"transactions": transactions,
	})
}

func TransactionFind(c *gin.Context) {
	tId := c.Param("id")
	utId, uuidErr := uuid.Parse(tId)
	if uuidErr != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "invalid uuid"})
		return
	}
	transaction := models.Transaction{}
	if res := intitializers.DB.Where("id = ?", utId).First(&transaction).Error; res != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "transaction not found"})
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
	user := c.MustGet("user").(models.User)
	cat := models.Category{ID: uuid.MustParse(body.Category)}
	prio := models.Priority{ID: uuid.MustParse(body.Priority)}
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
	transaction := models.Transaction{Title: body.Title, CategoryID: cat.ID, PriorityID: prio.ID, Amount: body.Amount, Negative: body.Negative, Description: body.Description, UserID: user.UID}
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
	prio := models.Priority{ID: uuid.MustParse(body.Priority)}
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
	transaction.CategoryID = cat.ID
	transaction.PriorityID = prio.ID
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
