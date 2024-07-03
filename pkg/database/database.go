// database/database.go

package database

import (
    "gorm.io/gorm"
	"gorm.io/driver/postgres"
    "log"
    "os"
    "fmt"
	"github.com/togzhanzhakhani/projects/internal/models"

    _ "github.com/lib/pq"
)

var db *gorm.DB // Переменная db будет хранить объект GORM DB

// Функция для инициализации базы данных и автомиграции
func SetupDatabase() {
    var err error

    // Создаем строку подключения к базе данных
    dbURL := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
        os.Getenv("POSTGRES_HOST"), os.Getenv("POSTGRES_PORT"), os.Getenv("POSTGRES_USER"),
        os.Getenv("POSTGRES_PASSWORD"), os.Getenv("POSTGRES_DB"))

    // Открываем соединение с базой данных
    db, err = gorm.Open(postgres.Open(dbURL), &gorm.Config{})
    if err != nil {
        log.Fatal(err)
    }

    // Автоматически создаем таблицы и миграции
    err = db.AutoMigrate(&models.User{}, &models.Task{}, &models.Project{})
    if err != nil {
        log.Fatal(err)
    }

    log.Println("Database migrated successfully")
}

// Функция для получения объекта DB GORM
func GetDB() *gorm.DB {
    return db
}

