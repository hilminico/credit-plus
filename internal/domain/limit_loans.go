package domain

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type LimitLoans struct {
	gorm.Model
	ID               uint64 `gorm:"primarykey"`
	UniqueIdentifier string `gorm:"type:char(36);"`
	CustomerId       string `gorm:"type:char(36);"`
	Year             uint8
	Month            uint8
	Limit            uint64
	CreatedAt        time.Time `gorm:"autoCreateTime"`
	UpdatedAt        time.Time `gorm:"autoUpdateTime"`
}

// BeforeCreate will set a UUID rather than numeric ID.
func (l *LimitLoans) BeforeCreate(tx *gorm.DB) (err error) {
	l.UniqueIdentifier = uuid.NewString()

	return
}
