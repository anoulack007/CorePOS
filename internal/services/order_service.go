package services

import (
	"errors"
	"fmt"
	"math"
	"sort"

	"github.com/anoulack007/core-pos/internal/core/domain"
	"github.com/anoulack007/core-pos/internal/core/ports"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type orderService struct {
	db   *gorm.DB
	repo ports.OrderRepository
}

type lockedOrderLine struct {
	product   domain.Product
	quantity  int
	unitCents int64
	costCents int64
}

func NewOrderService(db *gorm.DB, repo ports.OrderRepository) ports.OrderService {
	return &orderService{db: db, repo: repo}
}

func (s *orderService) GetAllOrders(storeID uuid.UUID) ([]domain.Order, error) {
	return s.repo.FindAll(storeID)
}

func (s *orderService) GetOrder(storeID, id uuid.UUID) (*domain.Order, error) {
	return s.repo.FindByID(storeID, id)
}

func (s *orderService) CreateOrder(storeID, userID uuid.UUID, lines []ports.OrderLineInput) (*domain.Order, error) {
	normalized, err := normalizeOrderLines(lines)
	if err != nil {
		return nil, err
	}

	order := &domain.Order{
		StoreID:       storeID,
		UserID:        userID,
		Status:        domain.OrderStatusCompleted,
		PaymentStatus: domain.PaymentStatusUnpaid,
	}
	err = s.db.Transaction(func(tx *gorm.DB) error {
		lockedLines := make([]lockedOrderLine, 0, len(normalized))
		var totalCents int64

		for _, line := range normalized {
			var product domain.Product
			err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
				Where("store_id = ? AND id = ?", storeID, line.ProductID).
				First(&product).Error
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("%w: %s", domain.ErrProductNotFound, line.ProductID)
			}
			if err != nil {
				return err
			}
			if product.StockQuantity < line.Quantity {
				return fmt.Errorf("%w: %s", domain.ErrInsufficientStock, product.Name)
			}

			unitCents := moneyToCents(product.Price)
			costCents := moneyToCents(product.CostPrice)
			if unitCents < 0 || costCents < 0 || int64(line.Quantity) > math.MaxInt64/maxInt64(unitCents, 1) {
				return fmt.Errorf("%w: invalid product price or quantity", domain.ErrInvalidOrder)
			}
			lineTotal := unitCents * int64(line.Quantity)
			if totalCents > math.MaxInt64-lineTotal {
				return fmt.Errorf("%w: order total is too large", domain.ErrInvalidOrder)
			}
			totalCents += lineTotal
			lockedLines = append(lockedLines, lockedOrderLine{
				product: product, quantity: line.Quantity, unitCents: unitCents, costCents: costCents,
			})
		}

		order.TotalAmount = centsToMoney(totalCents)
		if err := tx.Omit("Store", "User", "Items", "Payments").Create(order).Error; err != nil {
			return err
		}

		for _, line := range lockedLines {
			productID := line.product.ID
			item := domain.OrderItem{
				OrderID:           order.ID,
				ProductID:         &productID,
				Quantity:          line.quantity,
				UnitPrice:         centsToMoney(line.unitCents),
				CostPriceSnapshot: centsToMoney(line.costCents),
				Subtotal:          centsToMoney(line.unitCents * int64(line.quantity)),
			}
			if err := tx.Omit("Product").Create(&item).Error; err != nil {
				return err
			}

			newStock := line.product.StockQuantity - line.quantity
			if err := tx.Model(&domain.Product{}).
				Where("store_id = ? AND id = ?", storeID, productID).
				Update("stock_quantity", newStock).Error; err != nil {
				return err
			}

			movement := domain.InventoryMovement{
				StoreID:         storeID,
				ProductID:       productID,
				UserID:          &userID,
				MovementType:    "SALE",
				QuantityChanged: -line.quantity,
				ReferenceID:     &order.ID,
				Notes:           "Order checkout",
			}
			if err := tx.Create(&movement).Error; err != nil {
				return err
			}
			order.Items = append(order.Items, item)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return order, nil
}

func (s *orderService) VoidOrder(storeID, id, actorUserID uuid.UUID) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		var order domain.Order
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("store_id = ? AND id = ?", storeID, id).
			First(&order).Error; err != nil {
			return err
		}
		if order.Status == domain.OrderStatusVoid || order.Status == domain.OrderStatusRefund {
			return domain.ErrOrderNotVoidable
		}
		if order.PaymentStatus != domain.PaymentStatusUnpaid {
			return domain.ErrPaidOrderCannotVoid
		}

		var items []domain.OrderItem
		if err := tx.Where("order_id = ?", id).Order("product_id asc").Find(&items).Error; err != nil {
			return err
		}
		for _, item := range items {
			if item.ProductID == nil {
				return domain.ErrOrderNotVoidable
			}
			var product domain.Product
			if err := tx.Unscoped().Clauses(clause.Locking{Strength: "UPDATE"}).
				Where("store_id = ? AND id = ?", storeID, *item.ProductID).
				First(&product).Error; err != nil {
				return err
			}
			if err := tx.Unscoped().Model(&domain.Product{}).
				Where("store_id = ? AND id = ?", storeID, product.ID).
				Update("stock_quantity", product.StockQuantity+item.Quantity).Error; err != nil {
				return err
			}

			movement := domain.InventoryMovement{
				StoreID:         storeID,
				ProductID:       product.ID,
				UserID:          &actorUserID,
				MovementType:    "RETURN",
				QuantityChanged: item.Quantity,
				ReferenceID:     &order.ID,
				Notes:           "Order voided",
			}
			if err := tx.Create(&movement).Error; err != nil {
				return err
			}
		}

		return tx.Model(&domain.Order{}).
			Where("store_id = ? AND id = ?", storeID, id).
			Update("status", domain.OrderStatusVoid).Error
	})
}

func normalizeOrderLines(lines []ports.OrderLineInput) ([]ports.OrderLineInput, error) {
	if len(lines) == 0 {
		return nil, fmt.Errorf("%w: at least one item is required", domain.ErrInvalidOrder)
	}

	quantities := make(map[uuid.UUID]int, len(lines))
	for _, line := range lines {
		if line.ProductID == uuid.Nil || line.Quantity <= 0 {
			return nil, fmt.Errorf("%w: product and positive quantity are required", domain.ErrInvalidOrder)
		}
		if quantities[line.ProductID] > math.MaxInt-line.Quantity {
			return nil, fmt.Errorf("%w: quantity is too large", domain.ErrInvalidOrder)
		}
		quantities[line.ProductID] += line.Quantity
	}

	normalized := make([]ports.OrderLineInput, 0, len(quantities))
	for productID, quantity := range quantities {
		normalized = append(normalized, ports.OrderLineInput{ProductID: productID, Quantity: quantity})
	}
	sort.Slice(normalized, func(i, j int) bool {
		return normalized[i].ProductID.String() < normalized[j].ProductID.String()
	})
	return normalized, nil
}

func moneyToCents(amount float64) int64 {
	return int64(math.Round(amount * 100))
}

func centsToMoney(cents int64) float64 {
	return float64(cents) / 100
}

func maxInt64(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}
