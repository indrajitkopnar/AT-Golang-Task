package config

import (
	"fmt"
	"log"
	"os"

	"go-chat-app/models"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// DB instance
var DB *gorm.DB

// ConnectDatabase initializes the database connection
func ConnectDatabase() error {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: No .env file found, using default values")
	}

	// Load database credentials from .env
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	// Open PostgreSQL connection
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Println("❌ Failed to connect to the database:", err)
		return err
	}

	// Store database instance globally
	DB = db

	// Auto-Migrate tables
	db.AutoMigrate(&models.User{}, &models.Message{}) // ✅ Now includes messages

	fmt.Println("✅ Database connected & User table migrated!")
	return nil
}
