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

const maxAttempts = 10

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

		for attempt := 1; attempt <= maxAttempts; attempt++ {
			db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
				NowFunc: func() time.Time {
					return time.Now().UTC()
				},
			})

			if err == nil {
				sqlDB, err := db.DB()
				if err == nil {
					err = sqlDB.Ping()
					if err == nil {
						log.Println("Successfully connected to MySQL!")
						// Set connection pool settings
						sqlDB.SetMaxIdleConns(10)
						sqlDB.SetMaxOpenConns(100)
						sqlDB.SetConnMaxLifetime(time.Hour)

						break
					}
				}
			}

			log.Printf("Attempt %d/%d: Connection failed (%v), retrying in 5 seconds...", attempt, maxAttempts, err)
			time.Sleep(5 * time.Second)
		}
	})
	return db
}
