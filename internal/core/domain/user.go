package domain

import "time"

// User represents a staff member or cashier belonging to a store.
type User struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	StoreID      uint      `json:"store_id" gorm:"not null"`
	Store        Store     `json:"-" gorm:"foreignKey:StoreID;constraint:OnDelete:CASCADE"`
	Username     string    `json:"username" gorm:"type:varchar(100);uniqueIndex;not null"`
	PasswordHash string    `json:"-" gorm:"type:varchar(255);not null"`
	Role         string    `json:"role" gorm:"type:varchar(50);default:'cashier'"`
	CreatedAt    time.Time `json:"created_at" gorm:"autoCreateTime"`
}
