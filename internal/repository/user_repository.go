package repository

import (
	"gorm.io/gorm"
	
	"github.com/togzhanzhakhani/projects/internal/models"
)

type UserRepository interface {
	CreateUser(user *models.User) error
	FindByEmail(email string) (*models.User, error)
	GetAllUsers() ([]models.User, error)
	GetUserByID(id uint) (*models.User, error)
	UpdateUser(user *models.User) error
	DeleteUser(id uint) error
	FindByName(name string) ([]models.User, error)
	FindByEmailLike(email string) ([]models.User, error)
	GetTasksByUserID(userID uint) ([]models.Task, error)
}

type userRepository struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) *userRepository {
	return &userRepository{DB: db}
}

func (repo *userRepository) CreateUser(user *models.User) error {
	return repo.DB.Create(user).Error
}

func (repo *userRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	err := repo.DB.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (repo *userRepository) GetAllUsers() ([]models.User, error) {
	var users []models.User
	err := repo.DB.Find(&users).Error
	return users, err
}

func (repo *userRepository) GetUserByID(id uint) (*models.User, error) {
	var user models.User
	err := repo.DB.First(&user, id).Error
	return &user, err
}

func (repo *userRepository) UpdateUser(user *models.User) error {
	return repo.DB.Save(user).Error
}

func (repo *userRepository) DeleteUser(id uint) error {
	return repo.DB.Delete(&models.User{}, id).Error
}

func (repo *userRepository) FindByName(name string) ([]models.User, error) {
	var users []models.User
	err := repo.DB.Where("name LIKE ?", "%"+name+"%").Find(&users).Error
	return users, err
}

func (repo *userRepository) FindByEmailLike(email string) ([]models.User, error) {
	var users []models.User
	err := repo.DB.Where("email LIKE ?", "%"+email+"%").Find(&users).Error
	return users, err
}

func (repo *userRepository) GetTasksByUserID(userID uint) ([]models.Task, error) {
	var tasks []models.Task
	err := repo.DB.Where("assignee_id = ?", userID).Find(&tasks).Error
	return tasks, err
}