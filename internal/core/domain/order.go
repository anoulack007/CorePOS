package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrderStatus string

const (
	OrderStatusPending   OrderStatus = "pending"
	OrderStatusCompleted OrderStatus = "completed"
	OrderStatusVoid      OrderStatus = "void"
	OrderStatusRefund    OrderStatus = "refund"
)

type PaymentStatus string

const (
	PaymentStatusUnpaid  PaymentStatus = "unpaid"
	PaymentStatusPartial PaymentStatus = "partial"
	PaymentStatusPaid    PaymentStatus = "paid"
)

// Order represents a sales bill/receipt.
type Order struct {
	ID            uuid.UUID     `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	StoreID       uuid.UUID     `json:"store_id" gorm:"type:uuid;not null"`
	Store         Store         `json:"-" gorm:"foreignKey:StoreID;constraint:OnDelete:CASCADE"`
	UserID        uuid.UUID     `json:"user_id" gorm:"type:uuid;not null"`
	User          User          `json:"-" gorm:"foreignKey:UserID"`
	TotalAmount   float64       `json:"total_amount" gorm:"type:decimal(10,2);not null"`
	Status        OrderStatus   `json:"status" gorm:"type:varchar(20);default:'completed'"`
	PaymentStatus PaymentStatus `json:"payment_status" gorm:"type:varchar(20);default:'unpaid'"`
	Items         []OrderItem   `json:"items,omitempty" gorm:"foreignKey:OrderID"`
	Payments      []Payment     `json:"payments,omitempty" gorm:"foreignKey:OrderID"`
	CreatedAt     time.Time     `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt     time.Time     `json:"updated_at" gorm:"autoUpdateTime"`
}

func (o *Order) BeforeCreate(tx *gorm.DB) error {
	if o.ID == uuid.Nil {
		o.ID = uuid.New()
	}
	return nil
}

// OrderItem represents a product line within an order.
type OrderItem struct {
	ID                uuid.UUID  `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	OrderID           uuid.UUID  `json:"order_id" gorm:"type:uuid;not null"`
	ProductID         *uuid.UUID `json:"product_id" gorm:"type:uuid"`
	Product           *Product   `json:"product,omitempty" gorm:"foreignKey:ProductID;constraint:OnDelete:SET NULL"`
	Quantity          int        `json:"quantity" gorm:"not null"`
	UnitPrice         float64    `json:"unit_price" gorm:"type:decimal(10,2);not null"`
	CostPriceSnapshot float64   `json:"cost_price_snapshot" gorm:"type:decimal(10,2);not null"`
	Subtotal          float64    `json:"subtotal" gorm:"type:decimal(10,2);not null"`
}

func (oi *OrderItem) BeforeCreate(tx *gorm.DB) error {
	if oi.ID == uuid.Nil {
		oi.ID = uuid.New()
	}
	return nil
}

// Payment records a payment against an order.
type Payment struct {
	ID            uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	OrderID       uuid.UUID `json:"order_id" gorm:"type:uuid;not null"`
	Amount        float64   `json:"amount" gorm:"type:decimal(10,2);not null"`
	PaymentMethod string    `json:"payment_method" gorm:"type:varchar(50);not null"`
	PaidAt        time.Time `json:"paid_at" gorm:"autoCreateTime"`
}

func (p *Payment) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return nil
}
