package repository

import (
	"github.com/togzhanzhakhani/projects/internal/models"
	"gorm.io/gorm"
)

type TaskRepository interface {
	GetAllTasks() ([]models.Task, error)
	GetTaskByID(id uint) (*models.Task, error)
	CreateTask(task *models.Task) error
	UpdateTask(task *models.Task) error
	DeleteTask(id uint) error
	SearchTasksByTitle(title string) ([]models.Task, error)
	SearchTasksByStatus(status string) ([]models.Task, error)
	SearchTasksByPriority(priority string) ([]models.Task, error)
	SearchTasksByAssignee(userID uint) ([]models.Task, error)
	SearchTasksByProject(projectID uint) ([]models.Task, error)
	UserExists(userID int) bool
	ProjectExists(ProjectID int) bool
}

type taskRepository struct {
	DB *gorm.DB
}

func NewTaskRepository(db *gorm.DB) TaskRepository {
	return &taskRepository{
		DB: db,
	}
}

func (repo *taskRepository) GetAllTasks() ([]models.Task, error) {
	var tasks []models.Task
	err := repo.DB.Find(&tasks).Error
	return tasks, err
}

func (repo *taskRepository) GetTaskByID(id uint) (*models.Task, error) {
	var task models.Task
	err := repo.DB.First(&task, id).Error
	if err != nil {
		return nil, err
	}
	return &task, nil
}

func (repo *taskRepository) CreateTask(task *models.Task) error {
	return repo.DB.Create(task).Error
}

func (repo *taskRepository) UpdateTask(task *models.Task) error {
	return repo.DB.Save(task).Error
}

func (repo *taskRepository) DeleteTask(id uint) error {
	return repo.DB.Delete(&models.Task{}, id).Error
}

func (repo *taskRepository) SearchTasksByTitle(title string) ([]models.Task, error) {
	var tasks []models.Task
	err := repo.DB.Where("title LIKE ?", "%"+title+"%").Find(&tasks).Error
	return tasks, err
}

func (repo *taskRepository) SearchTasksByStatus(status string) ([]models.Task, error) {
	var tasks []models.Task
	err := repo.DB.Where("status = ?", status).Find(&tasks).Error
	return tasks, err
}

func (repo *taskRepository) SearchTasksByPriority(priority string) ([]models.Task, error) {
	var tasks []models.Task
	err := repo.DB.Where("priority = ?", priority).Find(&tasks).Error
	return tasks, err
}

func (repo *taskRepository) SearchTasksByAssignee(userID uint) ([]models.Task, error) {
	var tasks []models.Task
	err := repo.DB.Where("assignee_id = ?", userID).Find(&tasks).Error
	return tasks, err
}

func (repo *taskRepository) SearchTasksByProject(projectID uint) ([]models.Task, error) {
	var tasks []models.Task
	err := repo.DB.Where("project_id = ?", projectID).Find(&tasks).Error
	return tasks, err
}

func (tr *taskRepository) UserExists(userID int) bool {
    var count int64
    tr.DB.Model(&models.User{}).Where("id = ?", userID).Count(&count)
    return count > 0
}

func (tr *taskRepository) ProjectExists(ProjectID int) bool {
    var count int64
    tr.DB.Model(&models.Project{}).Where("id = ?", ProjectID).Count(&count)
    return count > 0
}