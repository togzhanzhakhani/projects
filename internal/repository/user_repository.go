package repository

import (
	"gorm.io/gorm"
	
	"github.com/togzhanzhakhani/projects/internal/models"
)

type UserRepository struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{DB: db}
}

func (repo *UserRepository) CreateUser(user *models.User) error {
	return repo.DB.Create(user).Error
}

func (repo *UserRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	err := repo.DB.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (repo *UserRepository) GetAllUsers() ([]models.User, error) {
	var users []models.User
	err := repo.DB.Find(&users).Error
	return users, err
}

func (repo *UserRepository) GetUserByID(id uint) (*models.User, error) {
	var user models.User
	err := repo.DB.First(&user, id).Error
	return &user, err
}

func (repo *UserRepository) UpdateUser(user *models.User) error {
	return repo.DB.Save(user).Error
}

func (repo *UserRepository) DeleteUser(id uint) error {
	return repo.DB.Delete(&models.User{}, id).Error
}

func (repo *UserRepository) FindByName(name string) ([]models.User, error) {
	var users []models.User
	err := repo.DB.Where("name LIKE ?", "%"+name+"%").Find(&users).Error
	return users, err
}

func (repo *UserRepository) FindByEmailLike(email string) ([]models.User, error) {
	var users []models.User
	err := repo.DB.Where("email LIKE ?", "%"+email+"%").Find(&users).Error
	return users, err
}

func (repo *UserRepository) GetTasksByUserID(userID uint) ([]models.Task, error) {
	var tasks []models.Task
	err := repo.DB.Where("assignee_id = ?", userID).Find(&tasks).Error
	return tasks, err
}