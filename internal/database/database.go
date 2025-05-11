package database

import (
	"github.com/ladnaaaaaa/calc_service/internal/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"os"
	"path/filepath"
)

var DB *gorm.DB

func InitDB() {
	// Создаем директорию для базы данных, если она не существует
	dbDir := "data"
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		log.Fatal("Failed to create database directory:", err)
	}

	// Путь к файлу базы данных
	dbPath := filepath.Join(dbDir, "calculator.db")

	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Auto migrate the schema
	err = db.AutoMigrate(&models.User{}, &models.Expression{}, &models.Task{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	DB = db
}
