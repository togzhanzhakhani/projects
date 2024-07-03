package handlers

import (
	"log"
	"net/http"
	"github.com/gin-gonic/gin"
	"strconv"

	"github.com/togzhanzhakhani/projects/internal/repository"
	"github.com/togzhanzhakhani/projects/internal/models"
	"github.com/togzhanzhakhani/projects/internal/validation"
)

type UserHandler struct {
	UserRepo repository.UserRepository
}

func NewUserHandler(userRepo repository.UserRepository) *UserHandler {
	return &UserHandler{UserRepo: userRepo}
}

func (uh *UserHandler) CreateUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if !validation.ValidateStruct(c, &user) {
		return
	}

	if _, err := uh.UserRepo.FindByEmail(user.Email); err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email already exists"})
		return
	}

	if err := uh.UserRepo.CreateUser(&user); err != nil {
		log.Printf("Error creating user: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, user)
}

func (uh *UserHandler) GetAllUsers(c *gin.Context) {
	users, err := uh.UserRepo.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve users"})
		return
	}
	c.JSON(http.StatusOK, users)
}

func (uh *UserHandler) GetUserByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	user, err := uh.UserRepo.GetUserByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (uh *UserHandler) UpdateUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if !validation.ValidateStruct(c, &user) {
		return
	}

	existingUser, err := uh.UserRepo.GetUserByID(uint(id))
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user"})
        return
    }

	if user.Email != existingUser.Email {
		if _, err := uh.UserRepo.FindByEmail(user.Email); err == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Email already exists"})
			return
		}
	}

	user.ID = uint(id)
	user.RegistrationDate = existingUser.RegistrationDate

	if err := uh.UserRepo.UpdateUser(&user); err != nil {
		log.Printf("Error updating user: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (uh *UserHandler) DeleteUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	if _, err := uh.UserRepo.GetUserByID(uint(id)); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if err := uh.UserRepo.DeleteUser(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}

	c.Status(http.StatusNoContent)
}

func (uh *UserHandler) GetTasksByUserID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	tasks, err := uh.UserRepo.GetTasksByUserID(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve tasks"})
		return
	}

	c.JSON(http.StatusOK, tasks)
}

func (uh *UserHandler) SearchUsersByName(c *gin.Context) {
	name := c.Query("name")
	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Name query parameter is required"})
		return
	}

	users, err := uh.UserRepo.FindByName(name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search users"})
		return
	}

	if len(users) == 0 {
        c.JSON(http.StatusNotFound, gin.H{"message": "No users found"})
        return
    }

	c.JSON(http.StatusOK, users)
}

func (uh *UserHandler) SearchUsersByEmail(c *gin.Context) {
	email := c.Query("email")
	if email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email query parameter is required"})
		return
	}

	users, err := uh.UserRepo.FindByEmailLike(email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search users"})
		return
	}

	if len(users) == 0 {
        c.JSON(http.StatusNotFound, gin.H{"message": "No users found"})
        return
    }

	c.JSON(http.StatusOK, users)
}