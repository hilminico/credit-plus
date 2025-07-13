package domain

import (
	"errors"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"regexp"
	"time"
)

type Customer struct {
	gorm.Model
	ID               uint64 `gorm:"primarykey"`
	UniqueIdentifier string `gorm:"type:char(36);"`
	Email            string `gorm:"not null;unique"`
	Password         string
	CreatedAt        time.Time `gorm:"autoCreateTime"`
	UpdatedAt        time.Time `gorm:"autoUpdateTime"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,min=8,max=72,excludesall= "`
}

type LoginResponse struct {
	Email        string `json:"email"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// BeforeCreate will set a UUID rather than numeric ID.
func (c *Customer) BeforeCreate(tx *gorm.DB) (err error) {
	c.UniqueIdentifier = uuid.NewString()

	if !isValidEmail(c.Email) {
		return gorm.ErrInvalidData
	}

	// Hash the password
	password, err := bcrypt.GenerateFromPassword([]byte(c.Password), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("could not hash password")
	}

	// assign to model
	c.Password = string(password)

	return
}

func isValidEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}
