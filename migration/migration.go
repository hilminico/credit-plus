package migration

import (
	"creditPlus/internal/config"
	"creditPlus/internal/domain"
	"log"
)

func main() {
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
