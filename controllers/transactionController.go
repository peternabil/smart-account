package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/peternabil/go-api/models"
)

func (server *Server) TransactionIndex(c *gin.Context) {
	user := server.store.GetUserFromToken(c)
	page, pageSize := getPaginationArgs(c.Request)
	var count int64
	transactions, err := server.store.GetTransactions(user.UID, page, pageSize, &count)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "no transactions for this user"})
		return
	}
	c.JSON(200, gin.H{
		"transactions": transactions,
		"page":         page,
		"page_size":    pageSize,
		// "total_pages":,
		"count": count,
	})
}

func (server *Server) TransactionFind(c *gin.Context) {
	tId := c.Param("id")
	utId, uuidErr := uuid.Parse(tId)
	if uuidErr != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "invalid uuid"})
		return
	}
	transaction := models.Transaction{ID: utId, UserID: server.store.GetUserFromToken(c).UID}
	res, err := server.store.GetTransaction(&transaction)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "transaction not found"})
		return
	}
	c.JSON(200, gin.H{
		"transaction": res,
	})
}

func (server *Server) TransactionCreate(c *gin.Context) {
	var body struct {
		Title       string
		Category    string
		Amount      int
		Negative    bool
		Description string
		Priority    string
	}
	err := c.BindJSON(&body)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	user := server.store.GetUserFromToken(c)
	cat := models.Category{ID: uuid.MustParse(body.Category)}
	prio := models.Priority{ID: uuid.MustParse(body.Priority)}
	category, err := server.store.GetCategory(user.UID, &cat)
	if err != nil {
		c.JSON(400, gin.H{"error": "category not found"})
		return
	}
	cat = category
	priority, err := server.store.GetPriority(user.UID, &prio)
	if err != nil {
		c.JSON(400, gin.H{"error": "priority not found"})
		return
	}
	prio = priority
	transaction := models.Transaction{Title: body.Title, CategoryID: cat.ID, PriorityID: prio.ID, Amount: body.Amount, Negative: body.Negative, Description: body.Description, UserID: user.UID}
	result, err := server.store.CreateTransaction(&transaction)
	if err != nil {
		c.JSON(500, gin.H{"error": "Could not create transaction"})
		return
	}
	c.JSON(200, gin.H{
		"transaction": result,
	})
}

func (server *Server) TransactionEdit(c *gin.Context) {
	var body struct {
		Title       string
		Category    string
		Amount      int
		Negative    bool
		Description string
		Priority    string
	}
	tId := c.Param("id")
	err := c.BindJSON(&body)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	user := server.store.GetUserFromToken(c)
	transaction := models.Transaction{ID: uuid.MustParse(tId), UserID: user.UID}
	res, err := server.store.GetTransaction(&transaction)
	transaction = res
	if err != nil {
		c.Status(404)
		return
	}
	cat := models.Category{ID: uuid.MustParse(body.Category)}
	prio := models.Priority{ID: uuid.MustParse(body.Priority)}
	category, err := server.store.GetCategory(cat.ID, &cat)
	if err != nil {
		c.Status(400)
		return
	}
	cat = category
	priority, err := server.store.GetPriority(prio.ID, &prio)
	if err != nil {
		c.Status(400)
		return
	}
	prio = priority
	transaction.Title = body.Title
	transaction.CategoryID = cat.ID
	transaction.PriorityID = prio.ID
	transaction.Amount = body.Amount
	transaction.Negative = body.Negative
	transaction.Description = body.Description
	transaction, err = server.store.EditTransaction(&transaction)
	if err != nil {
		c.Status(500)
		return
	}
	c.JSON(200, gin.H{
		"transaction": transaction,
	})
}

func (server *Server) TransactionDelete(c *gin.Context) {
	tId := c.Param("id")
	transaction := models.Transaction{ID: uuid.MustParse(tId)}
	res := server.store.DeleteTransaction(&transaction)
	if res != nil {
		c.Status(400)
		return
	}
	c.JSON(200, gin.H{
		"transaction": transaction,
	})
}
