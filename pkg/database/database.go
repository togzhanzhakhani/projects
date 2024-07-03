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

var db *gorm.DB

func SetupDatabase() {
    var err error

    dbURL := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
        os.Getenv("POSTGRES_HOST"), os.Getenv("POSTGRES_PORT"), os.Getenv("POSTGRES_USER"),
        os.Getenv("POSTGRES_PASSWORD"), os.Getenv("POSTGRES_DB"))

    db, err = gorm.Open(postgres.Open(dbURL), &gorm.Config{})
    if err != nil {
        log.Fatal(err)
    }

    err = db.AutoMigrate(&models.User{}, &models.Task{}, &models.Project{})
    if err != nil {
        log.Fatal(err)
    }

    log.Println("Database migrated successfully")
}

func GetDB() *gorm.DB {
    return db
}

