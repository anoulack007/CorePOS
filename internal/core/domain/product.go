package domain

import (
	"time"

	"gorm.io/gorm"
)

// Category groups products within a store.
type Category struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	StoreID   uint      `json:"store_id" gorm:"not null"`
	Store     Store     `json:"-" gorm:"foreignKey:StoreID;constraint:OnDelete:CASCADE"`
	Name      string    `json:"name" gorm:"type:varchar(100);not null"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
}

// Product represents a sellable item in a store.
type Product struct {
	ID            uint           `json:"id" gorm:"primaryKey"`
	StoreID       uint           `json:"store_id" gorm:"not null"`
	Store         Store          `json:"-" gorm:"foreignKey:StoreID;constraint:OnDelete:CASCADE"`
	CategoryID    *uint          `json:"category_id" gorm:"default:null"`
	Category      *Category      `json:"category,omitempty" gorm:"foreignKey:CategoryID;constraint:OnDelete:SET NULL"`
	Name          string         `json:"name" gorm:"type:varchar(255);not null"`
	Barcode       string         `json:"barcode" gorm:"type:varchar(100)"`
	CostPrice     float64        `json:"cost_price" gorm:"type:decimal(10,2);default:0.00"`
	Price         float64        `json:"price" gorm:"type:decimal(10,2);not null"`
	StockQuantity int            `json:"stock_quantity" gorm:"default:0"`
	CreatedAt     time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt     time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt     gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
}

// InventoryMovement tracks stock in/out movements.
type InventoryMovement struct {
	ID              uint      `json:"id" gorm:"primaryKey"`
	StoreID         uint      `json:"store_id" gorm:"not null"`
	Store           Store     `json:"-" gorm:"foreignKey:StoreID;constraint:OnDelete:CASCADE"`
	ProductID       uint      `json:"product_id" gorm:"not null"`
	Product         Product   `json:"-" gorm:"foreignKey:ProductID;constraint:OnDelete:CASCADE"`
	UserID          *uint     `json:"user_id"`
	User            *User     `json:"-" gorm:"foreignKey:UserID;constraint:OnDelete:SET NULL"`
	MovementType    string    `json:"movement_type" gorm:"type:varchar(50);not null"`
	QuantityChanged int       `json:"quantity_changed" gorm:"not null"`
	ReferenceID     *uint     `json:"reference_id"`
	Notes           string    `json:"notes" gorm:"type:text"`
	CreatedAt       time.Time `json:"created_at" gorm:"autoCreateTime"`
}
