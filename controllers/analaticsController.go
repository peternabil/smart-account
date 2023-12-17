package controllers

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/peternabil/go-api/models"
)

func setNegative(c *gin.Context, negative *bool) error {
	q := c.Request.URL.Query()
	negativeStr := q.Get("negative")
	var err error
	*negative, err = strconv.ParseBool(negativeStr)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return err
	}
	return nil
}

func setDates(c *gin.Context, startDate, endDate *time.Time) error {
	startDateVal, endDateVal, startDateErr, endDateErr := getDates(c.Request)
	if startDateErr != nil {
		c.JSON(400, gin.H{"error": startDateErr.Error()})
		return startDateErr
	}
	if endDateErr != nil {
		c.JSON(400, gin.H{"error": endDateErr.Error()})
		return endDateErr
	}
	*startDate = startDateVal
	*endDate = endDateVal
	return nil
}

func (server *Server) GetDailyValues(c *gin.Context) {
	var startDate time.Time
	var endDate time.Time
	var negative bool
	err := setDates(c, &startDate, &endDate)
	if err != nil {
		return
	}
	err = setNegative(c, &negative)
	if err != nil {
		return
	}
	user := c.MustGet("user").(models.User)
	transactions := []models.Transaction{}
	server.store.GetTransactionsDateRange(user.UID, &transactions, startDate, endDate, negative)
	// var spending []models.Spending
	// Amount   int
	// Category Category
	// Priority Priority
	// datetime time.Time
	// Negative bool
	// for index, obj := range transactions {
	// 	spending = append(spending, {Amount})
	// }
	c.JSON(200, gin.H{"transactions": transactions})
}
