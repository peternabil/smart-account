package controllers

import (
	"fmt"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
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
	user := server.store.GetUserFromToken(c)
	spendings, err := server.store.GetTransactionsDateRangeGroupByDay(user.UID, startDate, endDate, negative)
	if err != nil {
		c.Status(500)
		return
	}
	fmt.Println(spendings)
	c.JSON(200, gin.H{"spending": spendings})
}

func (server *Server) GetHighestCategory(c *gin.Context) {
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
	user := server.store.GetUserFromToken(c)
	spendings, err := server.store.GetHighestSpendingCategory(user.UID, startDate, endDate, negative)
	if err != nil {
		c.Status(500)
		return
	}
	c.JSON(200, gin.H{"spending": spendings})
}
func (server *Server) GetHighestPriority(c *gin.Context) {
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
	user := server.store.GetUserFromToken(c)
	spendings, err := server.store.GetHighestSpendingPriority(user.UID, startDate, endDate, negative)
	if err != nil {
		c.Status(500)
		return
	}
	c.JSON(200, gin.H{"spending": spendings})
}
