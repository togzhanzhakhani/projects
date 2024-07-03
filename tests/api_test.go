package tests

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/togzhanzhakhani/projects/internal/handlers"
	"github.com/togzhanzhakhani/projects/internal/models"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) CreateUser(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) FindByEmail(email string) (*models.User, error) {
    args := m.Called(email)
    if user, ok := args.Get(0).(*models.User); ok {
        return user, args.Error(1)
    }
    return nil, args.Error(1)
}

func (m *MockUserRepository) GetAllUsers() ([]models.User, error) {
	args := m.Called()
	return args.Get(0).([]models.User), args.Error(1)
}

func (m *MockUserRepository) GetUserByID(id uint) (*models.User, error) {
	args := m.Called(id)
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) UpdateUser(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) DeleteUser(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockUserRepository) FindByName(name string) ([]models.User, error) {
	args := m.Called(name)
	return args.Get(0).([]models.User), args.Error(1)
}

func (m *MockUserRepository) FindByEmailLike(email string) ([]models.User, error) {
	args := m.Called(email)
	return args.Get(0).([]models.User), args.Error(1)
}

func (m *MockUserRepository) GetTasksByUserID(userID uint) ([]models.Task, error) {
	args := m.Called(userID)
	return args.Get(0).([]models.Task), args.Error(1)
}

func setupUserHandler(t *testing.T) (*handlers.UserHandler, *MockUserRepository) {
	mockRepo := new(MockUserRepository)
	handler := handlers.NewUserHandler(mockRepo)
	return handler, mockRepo
}

func formatTime(t time.Time) string {
	return t.Format(time.RFC3339Nano)
}

func TestCreateUser_Success(t *testing.T) {
    handler, mockRepo := setupUserHandler(t)

    mockRepo.On("FindByEmail", "johndoe@example.com").Return(nil, errors.New("not found"))

    userJSON := `{"name":"John Doe","email":"johndoe@example.com","role":"admin"}`
    req, err := http.NewRequest("POST", "/users", bytes.NewBuffer([]byte(userJSON)))
    if err != nil {
        t.Fatal(err)
    }
    req.Header.Set("Content-Type", "application/json")

    mockRepo.On("CreateUser", mock.AnythingOfType("*models.User")).Return(nil)

    rr := httptest.NewRecorder()
    router := gin.Default()
    router.POST("/users", handler.CreateUser)
    router.ServeHTTP(rr, req)

    assert.Equal(t, http.StatusCreated, rr.Code, "статус код не соответствует ожидаемому")
}


func TestGetUserByID(t *testing.T) {
	handler, mockRepo := setupUserHandler(t)

	mockUser := &models.User{
		ID:              1,
		Name:            "John Doe",
		Email:           "johndoe@example.com",
		RegistrationDate: time.Now(),
		Role:            "admin",
	}
	mockRepo.On("GetUserByID", uint(1)).Return(mockUser, nil)

	req, err := http.NewRequest("GET", "/users/1", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router := gin.Default()
	router.GET("/users/:id", handler.GetUserByID)
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code, "статус код не соответствует ожидаемому")
	expected := `{"id":1,"name":"John Doe","email":"johndoe@example.com","registration_date":"` + formatTime(mockUser.RegistrationDate) + `","role":"admin"}`
	assert.JSONEq(t, expected, rr.Body.String(), "тело ответа не соответствует ожидаемому")
}

func TestUpdateUser(t *testing.T) {
	handler, mockRepo := setupUserHandler(t)

	userJSON := `{"name":"Jane Doe","email":"janedoe@example.com","role":"admin"}`
	req, err := http.NewRequest("PUT", "/users/1", bytes.NewBuffer([]byte(userJSON)))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	mockUser := &models.User{
		ID:    1,
		Name:  "John Doe",
		Email: "johndoe@example.com",
		Role:  "admin",
	}
	mockRepo.On("GetUserByID", uint(1)).Return(mockUser, nil)

	mockRepo.On("FindByEmail", "janedoe@example.com").Return(nil, errors.New("not found"))
	mockRepo.On("UpdateUser", mock.AnythingOfType("*models.User")).Return(nil)

	rr := httptest.NewRecorder()
	router := gin.Default()
	router.PUT("/users/:id", handler.UpdateUser)
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code, "статус код не соответствует ожидаемому")
}

func TestDeleteUser(t *testing.T) {
	handler, mockRepo := setupUserHandler(t)

	mockRepo.On("GetUserByID", uint(1)).Return(&models.User{ID: 1}, nil)

	mockRepo.On("DeleteUser", uint(1)).Return(nil)

	req, err := http.NewRequest("DELETE", "/users/1", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router := gin.Default()
	router.DELETE("/users/:id", handler.DeleteUser)
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNoContent, rr.Code, "статус код не соответствует ожидаемому")
}

func TestGetAllUsers(t *testing.T) {
	handler, mockRepo := setupUserHandler(t)

	mockUsers := []models.User{
		{ID: 1, Name: "John Doe", Email: "johndoe@example.com", RegistrationDate: time.Now(), Role: "admin"},
		{ID: 2, Name: "Jane Smith", Email: "janesmith@example.com", RegistrationDate: time.Now(), Role: "user"},
	}
	mockRepo.On("GetAllUsers").Return(mockUsers, nil)

	req, err := http.NewRequest("GET", "/users", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router := gin.Default()
	router.GET("/users", handler.GetAllUsers)
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code, "статус код не соответствует ожидаемому")
	expected := `[{"id":1,"name":"John Doe","email":"johndoe@example.com","registration_date":"` + formatTime(mockUsers[0].RegistrationDate) + `","role":"admin"},{"id":2,"name":"Jane Smith","email":"janesmith@example.com","registration_date":"` + formatTime(mockUsers[1].RegistrationDate) + `","role":"user"}]`
	assert.JSONEq(t, expected, rr.Body.String(), "тело ответа не соответствует ожидаемому")
}

func TestSearchUsersByName(t *testing.T) {
	handler, mockRepo := setupUserHandler(t)

	mockUsers := []models.User{
		{ID: 1, Name: "John Doe", Email: "johndoe@example.com", RegistrationDate: time.Now(), Role: "admin"},
		{ID: 2, Name: "Jane Smith", Email: "janesmith@example.com", RegistrationDate: time.Now(), Role: "user"},
	}
	mockRepo.On("FindByName", "John").Return(mockUsers, nil)

	req, err := http.NewRequest("GET", "/users/search?name=John", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router := gin.Default()
	router.GET("/users/search", handler.SearchUsersByName)
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code, "статус код не соответствует ожидаемому")
	expected := `[{"id":1,"name":"John Doe","email":"johndoe@example.com","registration_date":"` + formatTime(mockUsers[0].RegistrationDate) + `","role":"admin"},{"id":2,"name":"Jane Smith","email":"janesmith@example.com","registration_date":"` + formatTime(mockUsers[1].RegistrationDate) + `","role":"user"}]`
	assert.JSONEq(t, expected, rr.Body.String(), "тело ответа не соответствует ожидаемому")
}

func TestSearchUsersByEmail(t *testing.T) {
	handler, mockRepo := setupUserHandler(t)

	mockUsers := []models.User{
		{ID: 1, Name: "John Doe", Email: "johndoe@example.com", RegistrationDate: time.Now(), Role: "admin"},
		{ID: 2, Name: "Jane Smith", Email: "janesmith@example.com", RegistrationDate: time.Now(), Role: "user"},
	}
	mockRepo.On("FindByEmailLike", "example").Return(mockUsers, nil)

	req, err := http.NewRequest("GET", "/users/search?email=example", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router := gin.Default()
	router.GET("/users/search", handler.SearchUsersByEmail)
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code, "статус код не соответствует ожидаемому")
	expected := `[{"id":1,"name":"John Doe","email":"johndoe@example.com","registration_date":"` + formatTime(mockUsers[0].RegistrationDate) + `","role":"admin"},{"id":2,"name":"Jane Smith","email":"janesmith@example.com","registration_date":"` + formatTime(mockUsers[1].RegistrationDate) + `","role":"user"}]`
	assert.JSONEq(t, expected, rr.Body.String(), "тело ответа не соответствует ожидаемому")
}
