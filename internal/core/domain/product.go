package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Category groups products within a store.
type Category struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	StoreID   uuid.UUID `json:"store_id" gorm:"type:uuid;not null"`
	Store     Store     `json:"-" gorm:"foreignKey:StoreID;constraint:OnDelete:CASCADE"`
	Name      string    `json:"name" gorm:"type:varchar(100);not null"`
	IconURL   string    `json:"icon_url" gorm:"type:text"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
}

func (c *Category) BeforeCreate(tx *gorm.DB) error {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	return nil
}

// Product represents a sellable item in a store.

type Product struct {
	ID            uuid.UUID      `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	StoreID       uuid.UUID      `json:"store_id" gorm:"type:uuid;not null"`
	Store         Store          `json:"-" gorm:"foreignKey:StoreID;constraint:OnDelete:CASCADE"`
	CategoryID    *uuid.UUID     `json:"category_id" gorm:"type:uuid;default:null"`
	Category      *Category      `json:"category,omitempty" gorm:"foreignKey:CategoryID;constraint:OnDelete:SET NULL"`
	Name          string         `json:"name" gorm:"type:varchar(255);not null"`
	Barcode       string         `json:"barcode" gorm:"type:varchar(100)"`
	CostPrice     float64        `json:"cost_price" gorm:"type:decimal(10,2);default:0.00"`
	Price         float64        `json:"price" gorm:"type:decimal(10,2);not null"`
	StockQuantity int            `json:"stock_quantity" gorm:"default:0"`
	ImageURL      string         `json:"image_url" gorm:"type:text"`
	CreatedAt     time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt     time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt     gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
}

func (p *Product) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return nil
}

// InventoryMovement tracks stock in/out movements.
type InventoryMovement struct {
	ID              uuid.UUID  `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	StoreID         uuid.UUID  `json:"store_id" gorm:"type:uuid;not null"`
	Store           Store      `json:"-" gorm:"foreignKey:StoreID;constraint:OnDelete:CASCADE"`
	ProductID       uuid.UUID  `json:"product_id" gorm:"type:uuid;not null"`
	Product         Product    `json:"-" gorm:"foreignKey:ProductID;constraint:OnDelete:CASCADE"`
	UserID          *uuid.UUID `json:"user_id" gorm:"type:uuid"`
	User            *User      `json:"-" gorm:"foreignKey:UserID;constraint:OnDelete:SET NULL"`
	MovementType    string     `json:"movement_type" gorm:"type:varchar(50);not null"`
	QuantityChanged int        `json:"quantity_changed" gorm:"not null"`
	ReferenceID     *uuid.UUID `json:"reference_id" gorm:"type:uuid"`
	Notes           string     `json:"notes" gorm:"type:text"`
	CreatedAt       time.Time  `json:"created_at" gorm:"autoCreateTime"`
}

func (i *InventoryMovement) BeforeCreate(tx *gorm.DB) error {
	if i.ID == uuid.Nil {
		i.ID = uuid.New()
	}
	return nil
}
