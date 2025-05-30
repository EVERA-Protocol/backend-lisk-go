package database

import (
	"log"
	"rwa-backend/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	_ "modernc.org/sqlite" // Pure Go SQLite driver
)

var DB *gorm.DB

// InitDatabase initializes the SQLite database connection and runs migrations
func InitDatabase() {
	var err error

	// Connect to SQLite database using modernc.org/sqlite (pure Go driver)
	DB, err = gorm.Open(sqlite.Dialector{
		DriverName: "sqlite",
		DSN:        "assets.db",
	}, &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Auto migrate the schema
	err = DB.AutoMigrate(&models.Asset{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	log.Println("✅ Database connected and migrated successfully with pure Go SQLite driver")
}

// GetDB returns the database instance
func GetDB() *gorm.DB {
	return DB
}
