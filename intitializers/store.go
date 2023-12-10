package intitializers

import (
	"github.com/google/uuid"
	"github.com/peternabil/go-api/models"
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
	return DB.Where("id = ?", id).First(&transaction).Error
}
func (s MainStore) GetTransactions(id uuid.UUID, transactions *[]models.Transaction) error {
	return DB.Where("user_id = ?", id).Find(&transactions).Error
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