package store

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/peternabil/go-api/models"
)

type Store interface {
	ReadToken(tokenStr string) (models.User, error)
	CreateToken(user models.User) (string, error)
	GetUserFromToken(c *gin.Context) models.User
	CreateTransaction(transaction *models.Transaction) (models.Transaction, error)
	EditTransaction(transaction *models.Transaction) (models.Transaction, error)
	DeleteTransaction(transaction *models.Transaction) error
	GetTransaction(transaction *models.Transaction) (models.Transaction, error)
	GetTransactions(id uuid.UUID, page, pageSize int, count *int64) ([]models.Transaction, error)

	CreateCategory(category *models.Category) (models.Category, error)
	EditCategory(category *models.Category) (models.Category, error)
	DeleteCategory(category *models.Category) error
	GetCategory(id uuid.UUID, category *models.Category) (models.Category, error)
	GetCategories(id uuid.UUID) ([]models.Category, error)

	CreatePriority(priority *models.Priority) (models.Priority, error)
	EditPriority(priority *models.Priority) (models.Priority, error)
	DeletePriority(priority *models.Priority) error
	GetPriority(id uuid.UUID, priority *models.Priority) (models.Priority, error)
	GetPriorities(id uuid.UUID) ([]models.Priority, error)

	SignUp(user *models.User) (models.User, error)
	GetUser(user *models.User) (models.User, error)
	GetUsers() ([]models.User, error)
	FindUser(email string) (models.User, error)

	GetTransactionsDateRangeGroupByDay(id uuid.UUID, startDate, endDate time.Time, negative bool) ([]models.Spending, error)
	GetHighestSpendingCategory(id uuid.UUID, startDate, endDate time.Time, negative bool) ([]models.SpendingCategory, error)
	GetHighestSpendingPriority(id uuid.UUID, startDate, endDate time.Time, negative bool) ([]models.SpendingPriority, error)
	TotalSpending(id uuid.UUID, startDate, endDate time.Time, negative bool) ([]models.SpendingPriority, error)
}
