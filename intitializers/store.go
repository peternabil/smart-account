package intitializers

import (
	"time"

	"github.com/google/uuid"
	"github.com/peternabil/go-api/models"
	"gorm.io/gorm/clause"
)

type MainStore struct {
	name string
}

func (s MainStore) CreateTransaction(transaction *models.Transaction) error {
	return DB.Create(transaction).Error
}
func (s MainStore) EditTransaction(transaction *models.Transaction) error {
	return DB.Save(&transaction).Error
}
func (s MainStore) DeleteTransaction(transaction *models.Transaction) error {
	return DB.Delete(&transaction).Error
}
func (s MainStore) GetTransaction(id uuid.UUID, transaction *models.Transaction) error {
	return DB.Preload(clause.Associations).Where("id = ?", id).First(&transaction).Error
}
func (s MainStore) GetTransactions(id uuid.UUID, transactions *[]models.Transaction, page, pageSize int, count *int64) error {
	DB.Preload(clause.Associations).Where("user_id = ?", id).Find(&transactions).Count(count)

	return DB.Preload(clause.Associations).Order("created_at desc").Limit(pageSize).Offset((page-1)*pageSize).Where("user_id = ?", id).Find(&transactions).Error
}
func (s MainStore) CreateCategory(category *models.Category) error {
	return DB.Create(&category).Error
}
func (s MainStore) EditCategory(category *models.Category) error {
	return DB.Save(&category).Error
}
func (s MainStore) DeleteCategory(category *models.Category) error {
	return DB.Delete(&category).Error
}
func (s MainStore) GetCategory(id uuid.UUID, category *models.Category) error {
	return DB.First(&category).Error
}
func (s MainStore) GetCategories(id uuid.UUID, categories *[]models.Category) error {
	return DB.Find(&categories).Error
}

func (s MainStore) CreatePriority(priority *models.Priority) error {
	return DB.Create(&priority).Error
}

func (s MainStore) EditPriority(priority *models.Priority) error {
	return DB.Save(&priority).Error
}
func (s MainStore) DeletePriority(priority *models.Priority) error {
	return DB.Delete(&priority).Error
}
func (s MainStore) GetPriority(id uuid.UUID, priority *models.Priority) error {
	return DB.First(&priority).Error
}
func (s MainStore) GetPriorities(id uuid.UUID, priorities *[]models.Priority) error {
	return DB.Find(&priorities).Error
}

func (s MainStore) GetUsers(users *[]models.User) error {
	return DB.Find(&users).Error
}
func (s MainStore) GetUser(id uuid.UUID, user *models.User) error {
	return DB.Where("uid = ?", id).First(&user).Error
}
func (s MainStore) SignUp(user *models.User) error {
	return DB.Create(&user).Error
}
func (s MainStore) FindUser(email string, user *models.User) error {
	return DB.Where("email = ?", email).First(&user).Error
}

func (s MainStore) GetTransactionsDateRangeGroupByDay(id uuid.UUID, spendings *[]models.Spending, startDate, endDate time.Time, negative bool) error {
	return DB.Table("transactions").Select("date(created_at) as date, sum(amount) as total, negative as Negative").Where("user_id = ?", id).Group("date(created_at), negative").Having("date(created_at) BETWEEN ? AND ? AND negative = ?", startDate, endDate, negative).Order("date(created_at) ASC").Scan(spendings).Error
}

func (s MainStore) GetHighestSpendingCategory(id uuid.UUID, spendings *[]models.SpendingCategory, startDate, endDate time.Time, negative bool) error {
	return DB.Preload(clause.Associations).Table("transactions t , categories c").Select("sum(amount) as total, category_id, c.name as CName").Where("t.user_id = ? AND t.created_at BETWEEN ? AND ? AND negative = ? AND c.id = category_id", id, startDate, endDate, negative).Group("category_id, c.name").Order("total desc").Scan(spendings).Error
}

func (s MainStore) GetHighestSpendingPriority(id uuid.UUID, spendings *[]models.SpendingPriority, startDate, endDate time.Time, negative bool) error {
	return DB.Preload(clause.Associations).Table("transactions t , priorities p").Select("sum(amount) as total, priority_id, p.name as PName, p.level as Level").Where("t.user_id = ? AND t.created_at BETWEEN ? AND ? AND negative = ? AND p.id = priority_id", id, startDate, endDate, negative).Group("priority_id, p.name, p.level").Order("total desc").Scan(spendings).Error
}

func (s MainStore) TotalSpending(id uuid.UUID, spendings *[]models.SpendingPriority, startDate, endDate time.Time, negative bool) error {
	return DB.Preload(clause.Associations).Table("transactions t , priorities p").Select("sum(amount) as total, priority_id, p.name as PName, p.level as Level").Where("t.user_id = ? AND t.created_at BETWEEN ? AND ? AND negative = ? AND p.id = priority_id", id, startDate, endDate, negative).Group("priority_id, p.name, p.level").Order("total desc").Scan(spendings).Error
}
