package repository

import (
	"github.com/togzhanzhakhani/projects/internal/models"
	"gorm.io/gorm"
)

type ProjectRepository struct {
	DB *gorm.DB
}

func NewProjectRepository(db *gorm.DB) *ProjectRepository {
	return &ProjectRepository{DB: db}
}

func (pr *ProjectRepository) GetAllProjects() ([]models.Project, error) {
	var projects []models.Project
	if err := pr.DB.Find(&projects).Error; err != nil {
		return nil, err
	}
	return projects, nil
}

func (pr *ProjectRepository) CreateProject(project *models.Project) error {
	return pr.DB.Create(project).Error
}

func (pr *ProjectRepository) GetProjectByID(id uint) (*models.Project, error) {
	var project models.Project
	if err := pr.DB.First(&project, id).Error; err != nil {
		return nil, err
	}
	return &project, nil
}

func (pr *ProjectRepository) UpdateProject(project *models.Project) error {
	return pr.DB.Save(project).Error
}

func (pr *ProjectRepository) DeleteProject(id uint) error {
	return pr.DB.Delete(&models.Project{}, id).Error
}

func (pr *ProjectRepository) GetTasksByProjectID(id uint) ([]models.Task, error) {
	var tasks []models.Task
	if err := pr.DB.Where("project_id = ?", id).Find(&tasks).Error; err != nil {
		return nil, err
	}
	return tasks, nil
}

func (pr *ProjectRepository) SearchProjectsByTitle(title string) ([]models.Project, error) {
	var projects []models.Project
	if err := pr.DB.Where("name LIKE ?", "%"+title+"%").Find(&projects).Error; err != nil {
		return nil, err
	}
	return projects, nil
}

func (pr *ProjectRepository) SearchProjectsByManagerID(managerID uint) ([]models.Project, error) {
	var projects []models.Project
	if err := pr.DB.Where("manager_id = ?", managerID).Find(&projects).Error; err != nil {
		return nil, err
	}
	return projects, nil
}

func (pr *ProjectRepository) UserExists(userID int) bool {
    var count int64
    pr.DB.Model(&models.User{}).Where("id = ?", userID).Count(&count)
    return count > 0
}
