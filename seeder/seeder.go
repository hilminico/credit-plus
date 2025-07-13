package seeder

import (
	"creditPlus/internal/config"
	"creditPlus/internal/domain"
	"fmt"
	"gorm.io/gorm"
	"log"
)

type Seeder struct {
	db *gorm.DB
}

func main() {
	seeder := NewSeeder()
	seeder.Run()
}

func NewSeeder() *Seeder {
	return &Seeder{
		db: config.GetDB(),
	}
}

func (s *Seeder) Run() {
	fmt.Println("Running seeders...")
	s.SeedCustomer()
}

func (s *Seeder) SeedCustomer() {
	// Check if users already exist
	var count int64
	s.db.Model(&domain.Customer{}).Count(&count)
	if count > 0 {
		log.Println("Customer already seeded, skipping...")
		return
	}

	customers := []domain.Customer{
		{
			Email:    "test@gmail.com",
			Password: "test123",
		},
	}

	result := s.db.Create(&customers)
	if result.Error != nil {
		log.Fatalf("Failed to seed customers: %v", result.Error)
	}

	log.Printf("Seeded %d customers successfully\n", len(customers))
}
