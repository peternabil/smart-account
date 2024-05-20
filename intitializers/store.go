package intitializers

import (
	"errors"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/peternabil/go-api/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// MainStore represents the main store implementation.
type MainStore struct {
	DB *gorm.DB
}

// NewMainStore creates a new instance of the main store.
func NewMainStore(db *gorm.DB) *MainStore {
	return &MainStore{
		DB: db,
	}
}

func (s MainStore) ReadToken(tokenStr string) (models.User, error) {
	user := models.User{}
	if tokenStr == "" {
		return user, errors.New("you must be logged in to perform this request")
	}
	reqToken := tokenStr
	splitToken := strings.Split(reqToken, "Bearer ")
	if len(splitToken) == 1 {
		return user, errors.New("you must be logged in to perform this request")
	}
	reqToken = splitToken[1]
	claims := &models.Claims{}
	token, err := jwt.ParseWithClaims(reqToken, claims, func(token *jwt.Token) (any, error) {
		return []byte(os.Getenv("JWT_SECRET_KEY")), nil
	})
	if err != nil || !token.Valid {
		return user, errors.New("you must be logged in to perform this request")
	}
	email := claims.Email
	user.Email = email
	res, err := s.FindUser(email)
	user = res
	if err != nil {
		return user, errors.New("you must be logged in to perform this request")
	}
	return user, nil
}

func (s MainStore) CreateToken(user models.User) (string, error) {
	sampleSecretKey := []byte(os.Getenv("JWT_SECRET_KEY"))
	claims := &models.Claims{
		Email: user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(72 * time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(sampleSecretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
func (s MainStore) GetUserFromToken(c *gin.Context) models.User {
	return c.MustGet("user").(models.User)
}
func (s MainStore) CreateTransaction(transaction *models.Transaction) (models.Transaction, error) {
	err := DB.Create(transaction).Error
	return *transaction, err
}
func (s MainStore) EditTransaction(transaction *models.Transaction) (models.Transaction, error) {
	err := DB.Save(&transaction).Error
	return *transaction, err
}
func (s MainStore) DeleteTransaction(transaction *models.Transaction) error {
	return DB.Delete(&transaction).Error
}
func (s MainStore) GetTransaction(id uuid.UUID, transaction *models.Transaction) (models.Transaction, error) {
	err := DB.Preload(clause.Associations).First(&transaction).Error
	return *transaction, err
}
func (s MainStore) GetTransactions(id uuid.UUID, page, pageSize int, count *int64) ([]models.Transaction, error) {
	transactions := []models.Transaction{}
	DB.Preload(clause.Associations).Where("user_id = ?", id).Find(&transactions).Count(count)
	err := DB.Preload(clause.Associations).Order("created_at desc").Limit(pageSize).Offset((page-1)*pageSize).Where("user_id = ?", id).Find(&transactions).Error
	return transactions, err
}
func (s MainStore) CreateCategory(category *models.Category) (models.Category, error) {
	err := DB.Create(&category).Error
	return *category, err
}
func (s MainStore) EditCategory(category *models.Category) (models.Category, error) {
	err := DB.Save(&category).Error
	return *category, err
}
func (s MainStore) DeleteCategory(category *models.Category) error {
	return DB.Delete(&category).Error
}
func (s MainStore) GetCategory(id uuid.UUID, category *models.Category) (models.Category, error) {
	err := DB.Preload(clause.Associations).Where("user_id = ?", id).First(&category).Error
	return *category, err
}
func (s MainStore) GetCategories(id uuid.UUID) ([]models.Category, error) {
	categories := []models.Category{}
	err := DB.Preload(clause.Associations).Where("user_id = ?", id).Find(&categories).Error
	return categories, err
}

func (s MainStore) CreatePriority(priority *models.Priority) (models.Priority, error) {
	err := DB.Create(&priority).Error
	return *priority, err
}

func (s MainStore) EditPriority(priority *models.Priority) (models.Priority, error) {
	err := DB.Save(&priority).Error
	return *priority, err
}
func (s MainStore) DeletePriority(priority *models.Priority) error {
	return DB.Delete(&priority).Error
}
func (s MainStore) GetPriority(id uuid.UUID, priority *models.Priority) (models.Priority, error) {
	err := DB.Preload(clause.Associations).Where("user_id = ?", id).First(&priority).Error
	return *priority, err
}
func (s MainStore) GetPriorities(id uuid.UUID) ([]models.Priority, error) {
	priorities := []models.Priority{}
	err := DB.Preload(clause.Associations).Where("user_id = ?", id).Find(&priorities).Error
	return priorities, err
}

func (s MainStore) GetUsers() ([]models.User, error) {
	users := []models.User{}
	err := DB.Find(&users).Error
	return users, err
}
func (s MainStore) GetUser(user *models.User) (models.User, error) {
	err := DB.First(&user).Error
	return *user, err
}
func (s MainStore) SignUp(user *models.User) (models.User, error) {
	err := DB.Create(&user).Error
	return *user, err
}
func (s MainStore) FindUser(email string) (models.User, error) {
	user := models.User{}
	err := DB.Where("email = ?", email).First(&user).Error
	return user, err
}

func (s MainStore) GetTransactionsDateRangeGroupByDay(id uuid.UUID, startDate, endDate time.Time, negative bool) ([]models.Spending, error) {
	spendings := []models.Spending{}
	err := DB.Table("transactions").Select("date(created_at) as date, sum(amount) as total, negative as Negative").Where("user_id = ?", id).Group("date(created_at), negative").Having("date(created_at) BETWEEN ? AND ? AND negative = ?", startDate, endDate, negative).Order("date(created_at) ASC").Scan(&spendings).Error
	return spendings, err
}

func (s MainStore) GetHighestSpendingCategory(id uuid.UUID, startDate, endDate time.Time, negative bool) ([]models.SpendingCategory, error) {
	spendings := []models.SpendingCategory{}
	err := DB.Preload(clause.Associations).Table("transactions t , categories c").Select("sum(amount) as total, category_id, c.name as CName").Where("t.user_id = ? AND t.created_at BETWEEN ? AND ? AND negative = ? AND c.id = category_id", id, startDate, endDate, negative).Group("category_id, c.name").Order("total desc").Scan(&spendings).Error
	return spendings, err
}

func (s MainStore) GetHighestSpendingPriority(id uuid.UUID, startDate, endDate time.Time, negative bool) ([]models.SpendingPriority, error) {
	spendings := []models.SpendingPriority{}
	err := DB.Preload(clause.Associations).Table("transactions t , priorities p").Select("sum(amount) as total, priority_id, p.name as PName, p.level as Level").Where("t.user_id = ? AND t.created_at BETWEEN ? AND ? AND negative = ? AND p.id = priority_id", id, startDate, endDate, negative).Group("priority_id, p.name, p.level").Order("total desc").Scan(&spendings).Error
	return spendings, err
}

func (s MainStore) TotalSpending(id uuid.UUID, startDate, endDate time.Time, negative bool) ([]models.SpendingPriority, error) {
	spendings := []models.SpendingPriority{}
	err := DB.Preload(clause.Associations).Table("transactions t , priorities p").Select("sum(amount) as total, priority_id, p.name as PName, p.level as Level").Where("t.user_id = ? AND t.created_at BETWEEN ? AND ? AND negative = ? AND p.id = priority_id", id, startDate, endDate, negative).Group("priority_id, p.name, p.level").Order("total desc").Scan(&spendings).Error
	return spendings, err
}
