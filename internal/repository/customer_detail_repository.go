package repository

import (
	"context"
	"creditPlus/internal/domain"
	"errors"
	"gorm.io/gorm"
)

type CustomerDetailRepository struct {
	db *gorm.DB
}

func NewCustomerDetailRepository(db *gorm.DB) *CustomerDetailRepository {
	return &CustomerDetailRepository{db: db}
}

func (r *CustomerDetailRepository) FindByCustomerID(ctx context.Context, id string) (*domain.CustomerDetail, error) {
	resultChan := make(chan *domain.CustomerDetail)
	errChan := make(chan error, 1)

	go func() {
		var customerDetail domain.CustomerDetail
		err := r.db.WithContext(ctx).Where("customer_id = ?", id).Find(&customerDetail).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			resultChan <- &domain.CustomerDetail{} // Send empty model
		} else if err != nil {
			errChan <- err
			return
		}
		resultChan <- &customerDetail
	}()

	select {
	case customerDetail := <-resultChan:
		return customerDetail, nil
	case err := <-errChan:
		return nil, err
	}
}

func (r *CustomerDetailRepository) Create(ctx context.Context, data *domain.CustomerDetail) error {
	resultChan := make(chan error, 1)

	go func() {
		defer close(resultChan)
		resultChan <- r.db.WithContext(ctx).Create(data).Error
	}()

	return <-resultChan
}

func (r *CustomerDetailRepository) Update(ctx context.Context, user *domain.CustomerDetail) error {
	resultChan := make(chan error, 1)

	go func() {
		defer close(resultChan)
		resultChan <- r.db.WithContext(ctx).Save(user).Error
	}()

	return <-resultChan
}
