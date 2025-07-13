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
