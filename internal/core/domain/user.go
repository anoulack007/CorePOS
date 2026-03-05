package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// User represents a staff member or cashier belonging to a store.
type User struct {
	ID           uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	StoreID      uuid.UUID `json:"store_id" gorm:"type:uuid;not null"`
	Store        Store     `json:"-" gorm:"foreignKey:StoreID;constraint:OnDelete:CASCADE"`
	Username     string    `json:"username" gorm:"type:varchar(100);uniqueIndex;not null"`
	PasswordHash string    `json:"-" gorm:"type:varchar(255);not null"`
	Role         string    `json:"role" gorm:"type:varchar(50);default:'cashier'"`
	FullName     string    `json:"full_name" gorm:"type:varchar(255)"`
	Email        string    `json:"email" gorm:"type:varchar(255)"`
	Phone        string    `json:"phone" gorm:"type:varchar(20)"`
	AvatarURL    string    `json:"avatar_url" gorm:"type:text"`
	CreatedAt    time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt    time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}
