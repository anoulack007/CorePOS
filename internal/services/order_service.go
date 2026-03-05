package services

import (
	"fmt"

	"github.com/anoulack007/core-pos/internal/core/domain"
	"github.com/anoulack007/core-pos/internal/core/ports"
)

type orderService struct {
	orderRepo   ports.OrderRepository
	productRepo ports.ProductRepository
	paymentRepo ports.PaymentRepository
}

// NewOrderService creates a new OrderService.
func NewOrderService(
	orderRepo ports.OrderRepository,
	productRepo ports.ProductRepository,
	paymentRepo ports.PaymentRepository,
) ports.OrderService {
	return &orderService{
		orderRepo:   orderRepo,
		productRepo: productRepo,
		paymentRepo: paymentRepo,
	}
}

func (s *orderService) GetAllOrders(storeID uint) ([]domain.Order, error) {
	return s.orderRepo.FindAll(storeID)
}

func (s *orderService) GetOrder(storeID, id uint) (*domain.Order, error) {
	return s.orderRepo.FindByID(storeID, id)
}

func (s *orderService) CreateOrder(order *domain.Order) error {
	// Calculate totals from items
	var total float64
	for i := range order.Items {
		item := &order.Items[i]
		item.Subtotal = float64(item.Quantity) * item.UnitPrice
		total += item.Subtotal
	}
	order.TotalAmount = total
	order.Status = domain.OrderStatusCompleted
	order.PaymentStatus = domain.PaymentStatusUnpaid

	return s.orderRepo.Create(order)
}

func (s *orderService) VoidOrder(storeID, id uint) error {
	order, err := s.orderRepo.FindByID(storeID, id)
	if err != nil {
		return err
	}
	if order.Status == domain.OrderStatusVoid {
		return fmt.Errorf("order %d is already voided", id)
	}
	return s.orderRepo.UpdateStatus(storeID, id, domain.OrderStatusVoid)
}
