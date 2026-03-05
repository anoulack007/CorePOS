package domain

import "time"

// OrderStatus represents the state of an order.
type OrderStatus string

const (
	OrderStatusPending   OrderStatus = "pending"
	OrderStatusCompleted OrderStatus = "completed"
	OrderStatusVoid      OrderStatus = "void"
	OrderStatusRefund    OrderStatus = "refund"
)

// PaymentStatus represents the payment state of an order.
type PaymentStatus string

const (
	PaymentStatusUnpaid  PaymentStatus = "unpaid"
	PaymentStatusPartial PaymentStatus = "partial"
	PaymentStatusPaid    PaymentStatus = "paid"
)

// Order represents a sales bill/receipt.
type Order struct {
	ID            uint          `json:"id" gorm:"primaryKey"`
	StoreID       uint          `json:"store_id" gorm:"not null"`
	Store         Store         `json:"-" gorm:"foreignKey:StoreID;constraint:OnDelete:CASCADE"`
	UserID        uint          `json:"user_id" gorm:"not null"`
	User          User          `json:"-" gorm:"foreignKey:UserID"`
	TotalAmount   float64       `json:"total_amount" gorm:"type:decimal(10,2);not null"`
	Status        OrderStatus   `json:"status" gorm:"type:varchar(20);default:'completed'"`
	PaymentStatus PaymentStatus `json:"payment_status" gorm:"type:varchar(20);default:'unpaid'"`
	Items         []OrderItem   `json:"items,omitempty" gorm:"foreignKey:OrderID"`
	Payments      []Payment     `json:"payments,omitempty" gorm:"foreignKey:OrderID"`
	CreatedAt     time.Time     `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt     time.Time     `json:"updated_at" gorm:"autoUpdateTime"`
}

// OrderItem represents a product line within an order.
type OrderItem struct {
	ID                uint     `json:"id" gorm:"primaryKey"`
	OrderID           uint     `json:"order_id" gorm:"not null"`
	ProductID         *uint    `json:"product_id"`
	Product           *Product `json:"product,omitempty" gorm:"foreignKey:ProductID;constraint:OnDelete:SET NULL"`
	Quantity          int      `json:"quantity" gorm:"not null"`
	UnitPrice         float64  `json:"unit_price" gorm:"type:decimal(10,2);not null"`
	CostPriceSnapshot float64  `json:"cost_price_snapshot" gorm:"type:decimal(10,2);not null"`
	Subtotal          float64  `json:"subtotal" gorm:"type:decimal(10,2);not null"`
}

// Payment records a payment against an order (supports split payments).
type Payment struct {
	ID            uint      `json:"id" gorm:"primaryKey"`
	OrderID       uint      `json:"order_id" gorm:"not null"`
	Amount        float64   `json:"amount" gorm:"type:decimal(10,2);not null"`
	PaymentMethod string    `json:"payment_method" gorm:"type:varchar(50);not null"`
	PaidAt        time.Time `json:"paid_at" gorm:"autoCreateTime"`
}
