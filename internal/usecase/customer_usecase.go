package usecase

import (
	"context"
	"creditPlus/helper/token"
	"creditPlus/internal/domain"
	"creditPlus/internal/repository"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"os"
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
	err = bcrypt.CompareHashAndPassword([]byte(customer.Password), []byte(LoginRequest.Password))
	if err != nil {
		return nil, errors.New("customer.password-miss-match")
	}

	// generate token access
	jwtToken, err := token.CreateJWTToken(customer.UniqueIdentifier, customer.Email, os.Getenv("SECRET_KEY"))
	if err != nil {
		return nil, err
	}

	return &domain.LoginResponse{
		Email:    customer.Email,
		JwtToken: jwtToken,
	}, nil
}
