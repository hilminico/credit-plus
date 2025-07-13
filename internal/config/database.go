package config

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"os"
	"sync"
	"time"
)

var (
	db     *gorm.DB
	dbOnce sync.Once
)

func GetDB() *gorm.DB {
	dbOnce.Do(func() {
		dsn := fmt.Sprintf(
			"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			os.Getenv("DB_USER"),
			os.Getenv("DB_PASSWORD"),
			os.Getenv("DB_HOST"),
			os.Getenv("DB_PORT"),
			os.Getenv("DB_NAME"),
		)

		var err error
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
			NowFunc: func() time.Time {
				return time.Now().UTC()
			},
		})
		if err != nil {
			log.Fatalf("Failed to connect to database: %v", err)
		}

		sqlDB, err := db.DB()
		if err != nil {
			log.Fatalf("Failed to get SQL DB: %v", err)
		}

		// Set connection pool settings
		sqlDB.SetMaxIdleConns(10)
		sqlDB.SetMaxOpenConns(100)
		sqlDB.SetConnMaxLifetime(time.Hour)

		fmt.Println("MySQL database connection established")
	})
	return db
}
