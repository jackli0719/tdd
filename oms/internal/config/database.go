package config

import (
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// DB is the global database connection
var DB *gorm.DB

// InitDB initializes the database connection
func InitDB(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	// Connection pool settings
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)

	DB = db
	log.Println("Database connected successfully")
	return db, nil
}
