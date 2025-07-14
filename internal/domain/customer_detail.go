package domain

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type CustomerDetail struct {
	gorm.Model
	ID               uint64 `gorm:"primarykey"`
	UniqueIdentifier string `gorm:"type:char(36);"`
	CustomerId       string `gorm:"type:char(36);"`
	NIK              string
	FullName         string
	LegalName        string
	PlaceOfBirth     string
	DateBirth        time.Time
	Salary           uint64
	ImageKtp         string
	ImageSelf        string
	CreatedAt        time.Time `gorm:"autoCreateTime"`
	UpdatedAt        time.Time `gorm:"autoUpdateTime"`
}

// BeforeCreate will set a UUID rather than numeric ID.
func (c *CustomerDetail) BeforeCreate(tx *gorm.DB) (err error) {
	c.UniqueIdentifier = uuid.NewString()

	return
}

type CustomerDetailRequest struct {
	NIK          string `json:"nik" validate:"required"`
	FullName     string `json:"full_name" validate:"required"`
	LegalName    string `json:"legal_name" validate:"required"`
	PlaceOfBirth string `json:"place_of_birth" validate:"required"`
	DateBirth    string `json:"date_birth" validate:"required"`
	Salary       uint64 `json:"salary" validate:"required"`
	ImageKtp     string `json:"image_ktp" validate:"required"`
	ImageSelf    string `json:"image_self" validate:"required"`
}
