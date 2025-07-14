package migration

import (
	"creditPlus/internal/config"
	"creditPlus/internal/domain"
	"log"
)

func RunMigration() {
	db := config.GetDB()

	// Auto migrate models
	err := db.AutoMigrate(
		&domain.Customer{},
		&domain.CustomerDetail{},
		&domain.LimitLoans{},
		&domain.Transaction{},
	)
	if err != nil {
		log.Fatal(err)
	}
}
