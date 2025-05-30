package database

import (
	"log"
	"rwa-backend/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

// InitDatabase initializes the SQLite database connection and runs migrations
func InitDatabase() {
	var err error

	// Connect to SQLite database
	DB, err = gorm.Open(sqlite.Open("assets.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Auto migrate the schema
	err = DB.AutoMigrate(&models.Asset{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	log.Println("✅ Database connected and migrated successfully")
}

// GetDB returns the database instance
func GetDB() *gorm.DB {
	return DB
}
