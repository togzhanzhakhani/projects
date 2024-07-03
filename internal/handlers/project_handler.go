package handlers

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/togzhanzhakhani/projects/internal/models"
	"github.com/togzhanzhakhani/projects/internal/repository"
	"github.com/togzhanzhakhani/projects/internal/validation"
)

type ProjectHandler struct {
	ProjectRepo *repository.ProjectRepository
}

func NewProjectHandler(pr *repository.ProjectRepository) *ProjectHandler {
	return &ProjectHandler{ProjectRepo: pr}
}

func (ph *ProjectHandler) GetAllProjects(c *gin.Context) {
	projects, err := ph.ProjectRepo.GetAllProjects()
	if err != nil {
		log.Printf("Error retrieving projects: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve projects"})
		return
	}
	c.JSON(http.StatusOK, projects)
}
func (ph *ProjectHandler) processProject(c *gin.Context, id uint, isUpdate bool) {
	var input struct {
		Name        string `json:"name" validate:"required"`
		Description string `json:"description" validate:"required,max=100"`
		StartDate   string `json:"start_date" validate:"required"`
		EndDate     string `json:"end_date" validate:"required,gtfield=StartDate"`
		ManagerID   int    `json:"manager_id" validate:"required,gt=0,manager-exists"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	startDate, err := time.Parse("2006-01-02", input.StartDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start date format"})
		return
	}

	endDate, err := time.Parse("2006-01-02", input.EndDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end date format"})
		return
	}

	project := models.Project{
		Name:        input.Name,
		Description: input.Description,
		StartDate:   startDate,
		EndDate:     endDate,
		ManagerID:   input.ManagerID,
	}

	if !validation.ValidateStruct(c, &project) {
		return
	}

	if !ph.ProjectRepo.UserExists(project.ManagerID) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Manager does not exist"})
		return
	}

	if isUpdate {
		project.ID = int(id)
		err = ph.ProjectRepo.UpdateProject(&project)
	} else {
		err = ph.ProjectRepo.CreateProject(&project)
	}

	if err != nil {
		var errMsg string
		if isUpdate {
			errMsg = "Failed to update project"
		} else {
			errMsg = "Failed to create project"
		}
		log.Printf("Error %s: %v", errMsg, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": errMsg})
		return
	}

	if isUpdate {
		c.JSON(http.StatusOK, project)
	} else {
		c.JSON(http.StatusCreated, project)
	}
}

func (ph *ProjectHandler) UpdateProject(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	ph.processProject(c, uint(id), true)
}

func (ph *ProjectHandler) CreateProject(c *gin.Context) {
	ph.processProject(c, 0, false)
}

func (ph *ProjectHandler) GetProjectByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	project, err := ph.ProjectRepo.GetProjectByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
		return
	}

	c.JSON(http.StatusOK, project)
}

func (ph *ProjectHandler) DeleteProject(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	project, err := ph.ProjectRepo.GetProjectByID(uint(id))
	if err != nil {
		log.Printf("Error retrieving project: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve project"})
		return
	}

	if project == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
		return
	}

	if err := ph.ProjectRepo.DeleteProject(uint(id)); err != nil {
		log.Printf("Error deleting project: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete project"})
		return
	}

	c.Status(http.StatusNoContent)
}

func (ph *ProjectHandler) GetTasksByProjectID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	tasks, err := ph.ProjectRepo.GetTasksByProjectID(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve tasks for project"})
		return
	}

	if len(tasks) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No tasks found"})
		return
	}

	c.JSON(http.StatusOK, tasks)
}

func (ph *ProjectHandler) SearchProjectsByTitle(c *gin.Context) {
	title := c.Query("title")
	if title == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing title parameter"})
		return
	}

	projects, err := ph.ProjectRepo.SearchProjectsByTitle(title)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search projects by title"})
		return
	}

	if len(projects) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No projects found with the given title"})
		return
	}

	c.JSON(http.StatusOK, projects)
}

func (ph *ProjectHandler) SearchProjectsByManagerID(c *gin.Context) {
	managerID, err := strconv.ParseUint(c.Query("manager"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid manager ID"})
		return
	}

	projects, err := ph.ProjectRepo.SearchProjectsByManagerID(uint(managerID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search projects by manager ID"})
		return
	}

	if len(projects) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No projects found"})
		return
	}

	c.JSON(http.StatusOK, projects)
}
