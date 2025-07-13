package usecase

import (
	"context"
	"creditPlus/internal/domain"
	"creditPlus/internal/repository"
)

type CustomerService struct {
	customerRepo *repository.CustomerRepository
}

func NewCustomerService(
	customerRepo *repository.CustomerRepository,
) *CustomerService {
	return &CustomerService{
		customerRepo: customerRepo,
	}
}

func (s *CustomerService) Login(ctx context.Context, LoginRequest *domain.LoginRequest) (*domain.LoginResponse, error) {
	customer, err := s.customerRepo.FindByEmail(ctx, LoginRequest.Email)
	if err != nil {
		return nil, err
	}

	// check password

	// generate token access

	return &domain.LoginResponse{
		Email:        customer.Email,
		AccessToken:  "token",
		RefreshToken: "token",
	}, nil
}
