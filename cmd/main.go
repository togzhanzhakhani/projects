package main

import (
	"log"
	"github.com/gin-gonic/gin"
	"os"
	"net/http"

	"github.com/togzhanzhakhani/projects/internal/handlers"
	"github.com/togzhanzhakhani/projects/pkg/database"
	"github.com/togzhanzhakhani/projects/internal/repository"
)

func main() {
	database.SetupDatabase()
	db := database.GetDB()
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Error getting sqlDB from gorm DB: %v", err)
	}
	defer sqlDB.Close()

	router := gin.Default()
	
	userRepo := repository.NewUserRepository(db)
	taskRepo := repository.NewTaskRepository(db)
	projectRepo := repository.NewProjectRepository(db)
	
	userHandler := handlers.NewUserHandler(userRepo)
	taskHandler := handlers.NewTaskHandler(taskRepo)
	projectHandler := handlers.NewProjectHandler(projectRepo)
	
	userRoutes := router.Group("/users")
	{
		userRoutes.GET("/", userHandler.GetAllUsers)
		userRoutes.POST("/", userHandler.CreateUser)
		userRoutes.GET("/:id", userHandler.GetUserByID)
		userRoutes.PUT("/:id", userHandler.UpdateUser)
		userRoutes.DELETE("/:id", userHandler.DeleteUser)
		userRoutes.GET("/:id/tasks", userHandler.GetTasksByUserID)
		userRoutes.GET("/search", func(c *gin.Context) {
			if name := c.Query("name"); name != "" {
				userHandler.SearchUsersByName(c)
			} else if email := c.Query("email"); email != "" {
				userHandler.SearchUsersByEmail(c)
			} else {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Query parameter 'name' or 'email' is required"})
			}
		})
	}
	
	taskRoutes := router.Group("/tasks")
	{
		taskRoutes.GET("/", taskHandler.GetAllTasks)
		taskRoutes.POST("/", taskHandler.CreateTask)
		taskRoutes.GET("/:id", taskHandler.GetTaskByID)
		taskRoutes.PUT("/:id", taskHandler.UpdateTask)
		taskRoutes.DELETE("/:id", func(c *gin.Context) {
			if title := c.Query("title"); title != "" {
				taskHandler.SearchTasksByTitle(c)
			} else if status := c.Query("status"); status != "" {
				taskHandler.SearchTasksByStatus(c)
			} else if priority := c.Query("priority"); priority != "" {
				taskHandler.SearchTasksByPriority(c)
			} else if assignee := c.Query("assignee"); assignee != "" {
				taskHandler.SearchTasksByAssignee(c)
			} else if project := c.Query("project"); project != "" {
				taskHandler.SearchTasksByProject(c)
			} else {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Missing query parameter"})
			}
		})
	}

	projectRoutes := router.Group("/projects")
	{
		projectRoutes.GET("/", projectHandler.GetAllProjects)
		projectRoutes.POST("/", projectHandler.CreateProject)
		projectRoutes.GET("/:id", projectHandler.GetProjectByID)
		projectRoutes.PUT("/:id", projectHandler.UpdateProject)
		projectRoutes.DELETE("/:id", projectHandler.DeleteProject)
		projectRoutes.GET("/:id/tasks", projectHandler.GetTasksByProjectID)
		projectRoutes.GET("/search", func(c *gin.Context) {
			if title := c.Query("title"); title != "" {
				projectHandler.SearchProjectsByTitle(c)
			} else if manager := c.Query("manager"); manager != "" {
				projectHandler.SearchProjectsByManagerID(c)
			} else {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Missing query parameter"})
			}
		})
	}
	
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	err = router.Run(":" + port)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
