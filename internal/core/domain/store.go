package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Store represents a SaaS tenant (shop/store).
type Store struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Name      string    `json:"name" gorm:"type:varchar(255);not null"`
	PlanType  string    `json:"plan_type" gorm:"type:varchar(50);default:'free'"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

func (s *Store) BeforeCreate(tx *gorm.DB) error {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	return nil
}


// SubscriptionHistory tracks SaaS plan renewal history.
type SubscriptionHistory struct {
	ID            uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	StoreID       uuid.UUID `json:"store_id" gorm:"type:uuid;not null"`
	Store         Store     `json:"-" gorm:"foreignKey:StoreID;constraint:OnDelete:CASCADE"`
	PlanName      string    `json:"plan_name" gorm:"type:varchar(50);not null"`
	AmountPaid    float64   `json:"amount_paid" gorm:"type:decimal(10,2);not null"`
	StartDate     time.Time `json:"start_date" gorm:"type:date;not null"`
	EndDate       time.Time `json:"end_date" gorm:"type:date;not null"`
	PaymentStatus string    `json:"payment_status" gorm:"type:varchar(50);default:'success'"`
	CreatedAt     time.Time `json:"created_at" gorm:"autoCreateTime"`
}

func (s *SubscriptionHistory) BeforeCreate(tx *gorm.DB) error {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	return nil
}	