package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	sqlc "github.com/hydr0g3nz/poc_pos_restuarant/internal/adapter/sqlc/generated"
	"github.com/hydr0g3nz/poc_pos_restuarant/internal/domain/entity"
	"github.com/hydr0g3nz/poc_pos_restuarant/internal/domain/repository"
	"github.com/hydr0g3nz/poc_pos_restuarant/internal/domain/vo"
	"github.com/hydr0g3nz/poc_pos_restuarant/utils"
	"github.com/jackc/pgx/v5/pgxpool"
)

type orderRepository struct {
	db      *pgxpool.Pool
	queries *sqlc.Queries
}

func NewOrderRepository(db *pgxpool.Pool) repository.OrderRepository {
	return &orderRepository{
		db:      db,
		queries: sqlc.New(db),
	}
}

func (r *orderRepository) Create(ctx context.Context, order *entity.Order) (*entity.Order, error) {
	dbOrder, err := r.queries.CreateOrder(ctx, sqlc.CreateOrderParams{
		TableID: int32(order.TableID),
		Status:  sqlc.OrderStatus(order.Status.String()),
		Notes:   utils.ConvertToText(order.Notes),
	})
	if err != nil {
		return nil, err
	}

	return r.dbOrderToEntity(dbOrder)
}

func (r *orderRepository) GetByID(ctx context.Context, id int) (*entity.Order, error) {
	dbOrder, err := r.queries.GetOrderByID(ctx, int32(id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return r.dbOrderToEntity(dbOrder)
}

func (r *orderRepository) GetByIDWithItems(ctx context.Context, id int) (*entity.Order, error) {
	dbOrder, err := r.queries.GetOrderByID(ctx, int32(id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	// Get order items
	dbOrderItems, err := r.queries.GetOrderItemsByOrderID(ctx, int32(id))
	if err != nil {
		return nil, err
	}

	order, err := r.dbOrderToEntity(dbOrder)
	if err != nil {
		return nil, err
	}

	// Convert order items
	orderItems := make([]*entity.OrderItem, len(dbOrderItems))
	for i, dbItem := range dbOrderItems {
		orderItem, err := r.dbOrderItemToEntity(dbItem)
		if err != nil {
			return nil, err
		}
		orderItems[i] = orderItem
	}

	order.Items = orderItems
	return order, nil
}

func (r *orderRepository) Update(ctx context.Context, order *entity.Order) (*entity.Order, error) {
	var closedAt *time.Time
	if order.ClosedAt != nil {
		closedAt = order.ClosedAt
	}

	dbOrder, err := r.queries.UpdateOrder(ctx, sqlc.UpdateOrderParams{
		ID:       int32(order.ID),
		TableID:  int32(order.TableID),
		Status:   sqlc.OrderStatus(order.Status.String()),
		Notes:    utils.ConvertToText(order.Notes),
		ClosedAt: utils.ConvertToPGTimestamp(closedAt),
	})
	if err != nil {
		return nil, err
	}

	return r.dbOrderToEntity(dbOrder)
}

func (r *orderRepository) Delete(ctx context.Context, id int) error {
	return r.queries.DeleteOrder(ctx, int32(id))
}

func (r *orderRepository) List(ctx context.Context, limit, offset int) ([]*entity.Order, error) {
	dbOrders, err := r.queries.ListOrders(ctx, sqlc.ListOrdersParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, err
	}

	return r.dbOrdersToEntities(dbOrders)
}

func (r *orderRepository) ListByTable(ctx context.Context, tableID int, limit, offset int) ([]*entity.Order, error) {
	dbOrders, err := r.queries.ListOrdersByTable(ctx, sqlc.ListOrdersByTableParams{
		TableID: int32(tableID),
		Limit:   int32(limit),
		Offset:  int32(offset),
	})
	if err != nil {
		return nil, err
	}

	return r.dbOrdersToEntities(dbOrders)
}

func (r *orderRepository) GetOpenOrderByTable(ctx context.Context, tableID int) (*entity.Order, error) {
	dbOrder, err := r.queries.GetOpenOrderByTable(ctx, int32(tableID))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return r.dbOrderToEntity(dbOrder)
}

func (r *orderRepository) ListByStatus(ctx context.Context, status string, limit, offset int) ([]*entity.Order, error) {
	dbOrders, err := r.queries.ListOrdersByStatus(ctx, sqlc.ListOrdersByStatusParams{
		Status: sqlc.OrderStatus(status),
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, err
	}

	return r.dbOrdersToEntities(dbOrders)
}

func (r *orderRepository) ListByDateRange(ctx context.Context, startDate, endDate time.Time, limit, offset int) ([]*entity.Order, error) {
	dbOrders, err := r.queries.ListOrdersByDateRange(ctx, sqlc.ListOrdersByDateRangeParams{
		CreatedAt:   utils.ConvertToPGTimestamp(&startDate),
		CreatedAt_2: utils.ConvertToPGTimestamp(&endDate),
		Limit:       int32(limit),
		Offset:      int32(offset),
	})
	if err != nil {
		return nil, err
	}

	return r.dbOrdersToEntities(dbOrders)
}
func (r *orderRepository) GetOrderByQRCode(ctx context.Context, qrCode string) (*entity.Order, error) {
	dbOrder, err := r.queries.GetOrderByQRCode(ctx, utils.ConvertToText(qrCode))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	// Convert the database order to domain entity
	order, err := r.dbOrderToEntity(dbOrder)
	if err != nil {
		return nil, err
	}
	return order, nil
}

// Helper methods for conversion
func (r *orderRepository) dbOrderToEntity(dbOrder *sqlc.Order) (*entity.Order, error) {
	status, err := vo.NewOrderStatus(string(dbOrder.Status))
	if err != nil {
		return nil, err
	}

	order := &entity.Order{
		ID:        int(dbOrder.ID),
		TableID:   int(dbOrder.TableID),
		Status:    status,
		Notes:     utils.FromPgTextToString(dbOrder.Notes),
		CreatedAt: dbOrder.CreatedAt.Time,
		UpdatedAt: dbOrder.UpdatedAt.Time,
	}

	if dbOrder.ClosedAt.Valid {
		order.ClosedAt = &dbOrder.ClosedAt.Time
	}

	return order, nil
}

func (r *orderRepository) dbOrdersToEntities(dbOrders []*sqlc.Order) ([]*entity.Order, error) {
	entities := make([]*entity.Order, len(dbOrders))
	for i, dbOrder := range dbOrders {
		entity, err := r.dbOrderToEntity(dbOrder)
		if err != nil {
			return nil, err
		}
		entities[i] = entity
	}
	return entities, nil
}

func (r *orderRepository) dbOrderItemToEntity(dbItem *sqlc.OrderItem) (*entity.OrderItem, error) {
	money, err := vo.NewMoney(utils.FromPgNumericToFloat(dbItem.UnitPrice))
	if err != nil {
		return nil, err
	}

	orderItem := &entity.OrderItem{
		ID:        int(dbItem.ID),
		OrderID:   int(dbItem.OrderID),
		ItemID:    int(dbItem.ItemID),
		Quantity:  int(dbItem.Quantity),
		UnitPrice: money,
		Notes:     utils.FromPgTextToString(dbItem.Notes),
		CreatedAt: dbItem.CreatedAt.Time,
		UpdatedAt: dbItem.UpdatedAt.Time,
	}

	return orderItem, nil
}
