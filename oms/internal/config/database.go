package config

import (
	"log"
	"strings"

	"oms/internal/model"

	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// DB is the global database connection
var DB *gorm.DB

// InitDB initializes the database connection
func InitDB(dsn string) (*gorm.DB, error) {
	var db *gorm.DB
	var err error

	// Support both MySQL and SQLite based on DSN prefix
	if strings.HasPrefix(dsn, "sqlite://") || strings.HasPrefix(dsn, "file:") {
		// SQLite DSN: sqlite://oms.db or file:oms.db
		sqlitePath := strings.TrimPrefix(dsn, "sqlite://")
		sqlitePath = strings.TrimPrefix(sqlitePath, "file:")
		if sqlitePath == "" {
			sqlitePath = "oms.db"
		}
		db, err = gorm.Open(sqlite.Open(sqlitePath), &gorm.Config{})
		if err != nil {
			return nil, err
		}
		log.Println("SQLite database connected successfully")
	} else {
		// MySQL DSN
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err != nil {
			return nil, err
		}
		log.Println("MySQL database connected successfully")
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	// Connection pool settings (not applicable for SQLite)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)

	// AutoMigrate all models
	if err := autoMigrate(db); err != nil {
		return nil, err
	}

	DB = db
	return db, nil
}

// autoMigrate runs database migrations for all models
func autoMigrate(db *gorm.DB) error {
	log.Println("Running database migrations...")
	if err := db.AutoMigrate(&model.User{}, &model.Category{}, &model.Product{}, &model.Order{}, &model.OrderItem{}, &model.Staff{}, &model.Address{}, &model.Review{}); err != nil {
		return err
	}
	log.Println("Database migrations completed successfully")
	return nil
}
