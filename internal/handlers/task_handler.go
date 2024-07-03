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

type TaskHandler struct {
	TaskRepo repository.TaskRepository
}

func NewTaskHandler(taskRepo repository.TaskRepository) *TaskHandler {
	return &TaskHandler{
		TaskRepo: taskRepo,
	}
}

func (th *TaskHandler) GetAllTasks(c *gin.Context) {
	tasks, err := th.TaskRepo.GetAllTasks()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch tasks"})
		return
	}

	if len(tasks) == 0 {
		c.JSON(http.StatusOK, gin.H{"message": "No tasks found"})
		return
	}

	c.JSON(http.StatusOK, tasks)
}

func (th *TaskHandler) GetTaskByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	task, err := th.TaskRepo.GetTaskByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	c.JSON(http.StatusOK, task)
}

func (th *TaskHandler) processTask(c *gin.Context, id uint, isUpdate bool) {
	var input struct {
		Title       string `json:"title" validate:"required"`
		Description string `json:"description" validate:"required,max=100"`
		Priority    string `json:"priority" validate:"required"`
		Status      string `json:"status" validate:"required"`
		AssigneeID  int    `json:"assignee_id" validate:"required,gt=0"`
		ProjectID   int    `json:"project_id" validate:"required,gt=0"`
		CreatedAt   string `json:"created_at" validate:"required"`
		CompletedAt string `json:"completed_at" validate:"required,gtfield=CreatedAt"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	createdAt, err := time.Parse("2006-01-02", input.CreatedAt)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid created_at date format"})
		return
	}

	completedAt, err := time.Parse("2006-01-02", input.CompletedAt)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end date format"})
		return
	}

	task := models.Task{
		ID:          int(id),
		Title:       input.Title,
		Description: input.Description,
		Priority:    input.Priority,
		Status:      input.Status,
		AssigneeID:  input.AssigneeID,
		ProjectID:   input.ProjectID,
		CreatedAt:   createdAt,
		CompletedAt: completedAt,
	}

	if !validation.ValidateStruct(c, &task) {
		return
	}

	if !th.TaskRepo.UserExists(task.AssigneeID) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Assignee does not exist"})
		return
	}

	if !th.TaskRepo.ProjectExists(task.ProjectID) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Project does not exist"})
		return
	}

	if isUpdate {
		err = th.TaskRepo.UpdateTask(&task)
	} else {
		err = th.TaskRepo.CreateTask(&task)
	}

	if err != nil {
		var errMsg string
		if isUpdate {
			errMsg = "Failed to update task"
		} else {
			errMsg = "Failed to create task"
		}
		log.Printf("Error %s: %v", errMsg, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": errMsg})
		return
	}

	if isUpdate {
		c.JSON(http.StatusOK, task)
	} else {
		c.JSON(http.StatusCreated, task)
	}
}

func (th *TaskHandler) CreateTask(c *gin.Context) {
	th.processTask(c, 0, false)
}

func (th *TaskHandler) UpdateTask(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	th.processTask(c, uint(id), true)
}

func (th *TaskHandler) DeleteTask(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	task, err := th.TaskRepo.GetTaskByID(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch task"})
		return
	}

	if task == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	if err := th.TaskRepo.DeleteTask(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete task"})
		return
	}

	c.Status(http.StatusNoContent)
}

func (th *TaskHandler) SearchTasksByTitle(c *gin.Context) {
	title := c.Query("title")
	tasks, err := th.TaskRepo.SearchTasksByTitle(title)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search tasks by title"})
		return
	}

	if len(tasks) == 0 {
		c.JSON(http.StatusOK, gin.H{"message": "No tasks found with the given title"})
		return
	}

	c.JSON(http.StatusOK, tasks)
}

func (th *TaskHandler) SearchTasksByStatus(c *gin.Context) {
	status := c.Query("status")
	tasks, err := th.TaskRepo.SearchTasksByStatus(status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search tasks by status"})
		return
	}

	if len(tasks) == 0 {
		c.JSON(http.StatusOK, gin.H{"message": "No tasks found"})
		return
	}

	c.JSON(http.StatusOK, tasks)
}

func (th *TaskHandler) SearchTasksByPriority(c *gin.Context) {
	priority := c.Query("priority")
	tasks, err := th.TaskRepo.SearchTasksByPriority(priority)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search tasks by priority"})
		return
	}

	if len(tasks) == 0 {
		c.JSON(http.StatusOK, gin.H{"message": "No tasks found"})
		return
	}

	c.JSON(http.StatusOK, tasks)
}

func (th *TaskHandler) SearchTasksByAssignee(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Query("assignee"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid assignee ID"})
		return
	}

	tasks, err := th.TaskRepo.SearchTasksByAssignee(uint(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search tasks by assignee"})
		return
	}

	if len(tasks) == 0 {
		c.JSON(http.StatusOK, gin.H{"message": "No tasks found"})
		return
	}

	c.JSON(http.StatusOK, tasks)
}

func (th *TaskHandler) SearchTasksByProject(c *gin.Context) {
	projectID, err := strconv.ParseUint(c.Query("project"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	tasks, err := th.TaskRepo.SearchTasksByProject(uint(projectID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search tasks by project"})
		return
	}

	if len(tasks) == 0 {
		c.JSON(http.StatusOK, gin.H{"message": "No tasks found"})
		return
	}

	c.JSON(http.StatusOK, tasks)
}
