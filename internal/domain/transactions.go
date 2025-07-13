package domain

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Transaction struct {
	gorm.Model
	ID                uint64 `gorm:"primarykey"`
	UniqueIdentifier  string `gorm:"type:char(36);"`
	LimitLoanId       string `gorm:"type:char(36);"`
	ContractNumber    string `gorm:"not null;"`
	OTR               uint64
	AdminFee          uint64
	InstallmentAMount uint64
	AmountOfInterest  float64
	CreatedAt         time.Time `gorm:"autoCreateTime"`
	UpdatedAt         time.Time `gorm:"autoUpdateTime"`
}

// BeforeCreate will set a UUID rather than numeric ID.
func (t *Transaction) BeforeCreate(tx *gorm.DB) (err error) {
	t.UniqueIdentifier = uuid.NewString()

	return
}
