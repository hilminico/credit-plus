package usecase

import (
	"context"
	"creditPlus/internal/domain"
	"creditPlus/internal/repository"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"os"
	"time"
)

type CustomerService struct {
	customerRepo       *repository.CustomerRepository
	customerDetailRepo *repository.CustomerDetailRepository
}

func NewCustomerService(
	customerRepo *repository.CustomerRepository,
	customerDetailRepo *repository.CustomerDetailRepository,
) *CustomerService {
	return &CustomerService{
		customerRepo:       customerRepo,
		customerDetailRepo: customerDetailRepo,
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

	// Create token
	claims := &domain.CustomerClaims{
		UniqueIdentifier: customer.UniqueIdentifier,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
		},
	}

	tokenJwt := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	token, err := tokenJwt.SignedString([]byte(os.Getenv("SECRET_KEY")))
	if err != nil {
		return nil, errors.New("general.token_failed")
	}

	return &domain.LoginResponse{
		Email:    customer.Email,
		JwtToken: token,
	}, nil
}

func (s *CustomerService) GetCustomer(ctx context.Context, id string) (*domain.Customer, error) {
	customer, err := s.customerRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return customer, nil
}

func (s *CustomerService) UpdateCustomerDetail(ctx context.Context, id string, updateData *domain.CustomerDetailRequest) (*domain.CustomerDetail, error) {
	customerDetail, err := s.customerDetailRepo.FindByCustomerID(ctx, id)
	if err != nil {
		return nil, err
	}

	if customerDetail.ID == 0 {
		err := s.customerDetailRepo.Create(ctx, &domain.CustomerDetail{
			CustomerId:   id,
			NIK:          updateData.NIK,
			FullName:     updateData.FullName,
			LegalName:    updateData.LegalName,
			PlaceOfBirth: updateData.PlaceOfBirth,
			DateBirth:    updateData.DateBirth,
			Salary:       updateData.Salary,
			ImageKtp:     updateData.ImageKtp,
			ImageSelf:    updateData.ImageSelf,
		})
		if err != nil {
			return nil, err
		}

	} else {
		customerDetail.NIK = updateData.NIK
		customerDetail.FullName = updateData.FullName
		customerDetail.LegalName = updateData.LegalName
		customerDetail.PlaceOfBirth = updateData.PlaceOfBirth
		customerDetail.DateBirth = updateData.DateBirth
		customerDetail.Salary = updateData.Salary
		customerDetail.ImageKtp = updateData.ImageKtp
		customerDetail.ImageSelf = updateData.ImageSelf

		err = s.customerDetailRepo.Update(ctx, customerDetail)
		if err != nil {
			return nil, err
		}
	}

	customerDetail, err = s.customerDetailRepo.FindByCustomerID(ctx, id)
	if err != nil {
		return nil, err
	}

	return customerDetail, nil
}
