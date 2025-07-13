package repository

import (
	"context"
	"creditPlus/internal/domain"
	"gorm.io/gorm"
)

type CustomerRepository struct {
	db *gorm.DB
}

func NewCustomerRepository(db *gorm.DB) *CustomerRepository {
	return &CustomerRepository{db: db}
}

func (r *CustomerRepository) FindByEmail(ctx context.Context, email string) (*domain.Customer, error) {
	resultChan := make(chan *domain.Customer)
	errChan := make(chan error, 1)

	go func() {
		var customer domain.Customer
		err := r.db.WithContext(ctx).Where("email = ?", email).First(&customer).Error
		if err != nil {
			errChan <- err
			return
		}
		resultChan <- &customer
	}()

	select {
	case customer := <-resultChan:
		return customer, nil
	case err := <-errChan:
		return nil, err
	}
}
