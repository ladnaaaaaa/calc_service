package database

import (
	"log"
	"os"
	"testing"

	"github.com/ladnaaaaaa/calc_service/internal/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Init() {
	var err error

	// Настраиваем логгер в зависимости от окружения
	var logLevel logger.LogLevel
	if os.Getenv("GIN_MODE") == "test" {
		logLevel = logger.Silent
	} else {
		logLevel = logger.Info
	}

	config := &gorm.Config{
		Logger:      logger.Default.LogMode(logLevel),
		PrepareStmt: true, // Включаем подготовленные выражения
	}

	DB, err = gorm.Open(sqlite.Open("calc_service.db"), config)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Настраиваем пул соединений
	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatal("Failed to get database instance:", err)
	}

	// Устанавливаем оптимальные параметры пула соединений
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(0) // Не ограничиваем время жизни соединений

	// Auto migrate the schema
	err = DB.AutoMigrate(
		&models.User{},
		&models.Expression{},
		&models.Task{},
	)
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}
}

// InitTestDB инициализирует базу данных для тестов
func InitTestDB(t *testing.T) {
	var err error

	// Используем in-memory базу данных для тестов
	DB, err = gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{
		Logger:      logger.Default.LogMode(logger.Silent),
		PrepareStmt: true,
	})
	if err != nil {
		t.Fatal("Failed to connect to test database:", err)
	}

	// Auto migrate the schema
	err = DB.AutoMigrate(
		&models.User{},
		&models.Expression{},
		&models.Task{},
	)
	if err != nil {
		t.Fatal("Failed to migrate test database:", err)
	}
}

// ClearDB очищает все таблицы в базе данных
func ClearDB() {
	DB.Exec("DELETE FROM users")
	DB.Exec("DELETE FROM expressions")
	DB.Exec("DELETE FROM tasks")
}
