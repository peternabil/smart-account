package store

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/peternabil/go-api/models"
)

type Store interface {
	ParseToken(c *gin.Context)
	CreateTransaction(transaction *models.Transaction) error
	EditTransaction(transaction *models.Transaction) error
	DeleteTransaction(transaction *models.Transaction) error
	GetTransaction(id uuid.UUID, transaction *models.Transaction) error
	GetTransactions(id uuid.UUID, transactions *[]models.Transaction, page, pageSize int, count *int64) error

	CreateCategory(category *models.Category) error
	EditCategory(category *models.Category) error
	DeleteCategory(category *models.Category) error
	GetCategory(id uuid.UUID, category *models.Category) error
	GetCategories(id uuid.UUID, categories *[]models.Category) error

	CreatePriority(priority *models.Priority) error
	EditPriority(priority *models.Priority) error
	DeletePriority(priority *models.Priority) error
	GetPriority(id uuid.UUID, priority *models.Priority) error
	GetPriorities(id uuid.UUID, priorities *[]models.Priority) error

	SignUp(user *models.User) error
	GetUser(id uuid.UUID, user *models.User) error
	GetUsers(users *[]models.User) error
	FindUser(email string, user *models.User) error

	GetTransactionsDateRangeGroupByDay(id uuid.UUID, spendings *[]models.Spending, startDate, endDate time.Time, negative bool) error
	GetHighestSpendingCategory(id uuid.UUID, spendings *[]models.SpendingCategory, startDate, endDate time.Time, negative bool) error
	GetHighestSpendingPriority(id uuid.UUID, spendings *[]models.SpendingPriority, startDate, endDate time.Time, negative bool) error
}
