package service

import (
	"context"
	"fmt"

	"github.com/hydr0g3nz/poc_pos_restuarant/internal/domain/entity"
	errs "github.com/hydr0g3nz/poc_pos_restuarant/internal/domain/error"
	"github.com/hydr0g3nz/poc_pos_restuarant/internal/domain/repository"
	"github.com/hydr0g3nz/poc_pos_restuarant/internal/domain/vo"
)

// OrderService provides domain logic for orders
type OrderService interface {
	// ValidateOrderCreation validates if order can be created for table
	ValidateOrderCreation(ctx context.Context, tableID int) error

	// CalculateOrderTotal calculates total amount for order
	CalculateOrderTotal(ctx context.Context, order *entity.Order) (vo.Money, error)

	// ValidateOrderItem validates order item before adding
	ValidateOrderItem(ctx context.Context, orderID, itemID int, quantity int) error

	// ProcessOrderClosure processes order closure with validations
	ProcessOrderClosure(ctx context.Context, orderID int) error
}

type orderService struct {
	orderRepo     repository.OrderRepository
	orderItemRepo repository.OrderItemRepository
	tableRepo     repository.TableRepository
	menuItemRepo  repository.MenuItemRepository
}

func NewOrderService(
	orderRepo repository.OrderRepository,
	orderItemRepo repository.OrderItemRepository,
	tableRepo repository.TableRepository,
	menuItemRepo repository.MenuItemRepository,
) OrderService {
	return &orderService{
		orderRepo:     orderRepo,
		orderItemRepo: orderItemRepo,
		tableRepo:     tableRepo,
		menuItemRepo:  menuItemRepo,
	}
}

func (s *orderService) ValidateOrderCreation(ctx context.Context, tableID int) error {
	// Check if table exists
	table, err := s.tableRepo.GetByID(ctx, tableID)
	if err != nil {
		return fmt.Errorf("failed to get table: %w", err)
	}
	if table == nil {
		return errs.ErrTableNotFound
	}

	// Check if table already has an open order
	openOrder, err := s.orderRepo.GetOpenOrderByTable(ctx, tableID)
	if err != nil {
		return fmt.Errorf("failed to check open order: %w", err)
	}
	if openOrder != nil {
		return errs.ErrTableAlreadyHasOpenOrder
	}

	return nil
}

func (s *orderService) CalculateOrderTotal(ctx context.Context, order *entity.Order) (vo.Money, error) {
	if order == nil {
		return vo.Money{}, errs.ErrOrderNotFound
	}

	// Get order items if not loaded
	if order.Items == nil {
		items, err := s.orderItemRepo.ListByOrder(ctx, order.ID)
		if err != nil {
			return vo.Money{}, fmt.Errorf("failed to get order items: %w", err)
		}
		order.Items = items
	}

	return order.CalculateTotal(), nil
}

func (s *orderService) ValidateOrderItem(ctx context.Context, orderID, itemID int, quantity int) error {
	// Check if order exists and is open
	order, err := s.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return fmt.Errorf("failed to get order: %w", err)
	}
	if order == nil {
		return errs.ErrOrderNotFound
	}
	if !order.IsOpen() {
		return errs.ErrOrderNotOpen
	}

	// Check if menu item exists
	menuItem, err := s.menuItemRepo.GetByID(ctx, itemID)
	if err != nil {
		return fmt.Errorf("failed to get menu item: %w", err)
	}
	if menuItem == nil {
		return errs.ErrMenuItemNotFound
	}

	// Validate quantity
	if quantity <= 0 {
		return errs.ErrInvalidQuantity
	}

	return nil
}

func (s *orderService) ProcessOrderClosure(ctx context.Context, orderID int) error {
	// Check if order exists
	order, err := s.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return fmt.Errorf("failed to get order: %w", err)
	}
	if order == nil {
		return errs.ErrOrderNotFound
	}

	// Check if order is already closed
	if order.IsClosed() {
		return errs.ErrOrderAlreadyClosed
	}

	// Check if order has items
	items, err := s.orderItemRepo.ListByOrder(ctx, orderID)
	if err != nil {
		return fmt.Errorf("failed to get order items: %w", err)
	}
	if len(items) == 0 {
		return errs.ErrEmptyOrder
	}

	return nil
}
