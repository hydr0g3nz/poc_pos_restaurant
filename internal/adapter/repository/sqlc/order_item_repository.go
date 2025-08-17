package repository

import (
	"context"
	"database/sql"
	"errors"

	sqlc "github.com/hydr0g3nz/poc_pos_restuarant/internal/adapter/sqlc/generated"
	"github.com/hydr0g3nz/poc_pos_restuarant/internal/domain/entity"
	"github.com/hydr0g3nz/poc_pos_restuarant/internal/domain/repository"
	"github.com/hydr0g3nz/poc_pos_restuarant/internal/domain/vo"
	"github.com/hydr0g3nz/poc_pos_restuarant/utils"
	"github.com/jackc/pgx/v5/pgxpool"
)

type orderItemRepository struct {
	db      *pgxpool.Pool
	queries *sqlc.Queries
}

func NewOrderItemRepository(db *pgxpool.Pool) repository.OrderItemRepository {
	return &orderItemRepository{
		db:      db,
		queries: sqlc.New(db),
	}
}

func (r *orderItemRepository) Create(ctx context.Context, item *entity.OrderItem) (*entity.OrderItem, error) {
	dbOrderItem, err := r.queries.CreateOrderItem(ctx, sqlc.CreateOrderItemParams{
		OrderID:   int32(item.OrderID),
		ItemID:    int32(item.ItemID),
		Quantity:  int32(item.Quantity),
		UnitPrice: utils.ConvertToPGNumericFromFloat(item.UnitPrice.AmountBaht()),
		Notes:     utils.ConvertToText(item.Notes),
		Name:      utils.ConvertToText(item.Name),
	})
	if err != nil {
		return nil, err
	}

	return r.dbOrderItemToEntity(dbOrderItem)
}

func (r *orderItemRepository) GetByID(ctx context.Context, id int) (*entity.OrderItem, error) {
	dbOrderItem, err := r.queries.GetOrderItemByID(ctx, int32(id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return r.dbOrderItemToEntity(dbOrderItem)
}

func (r *orderItemRepository) Update(ctx context.Context, item *entity.OrderItem) (*entity.OrderItem, error) {
	dbOrderItem, err := r.queries.UpdateOrderItem(ctx, sqlc.UpdateOrderItemParams{
		ID:        int32(item.ID),
		OrderID:   int32(item.OrderID),
		ItemID:    int32(item.ItemID),
		Quantity:  int32(item.Quantity),
		UnitPrice: utils.ConvertToPGNumericFromFloat(item.UnitPrice.AmountBaht()),
		Notes:     utils.ConvertToText(item.Notes),
		Name:      utils.ConvertToText(item.Name),
	})
	if err != nil {
		return nil, err
	}

	return r.dbOrderItemToEntity(dbOrderItem)
}

func (r *orderItemRepository) Delete(ctx context.Context, id int) error {
	return r.queries.DeleteOrderItem(ctx, int32(id))
}

func (r *orderItemRepository) ListByOrder(ctx context.Context, orderID int) ([]*entity.OrderItem, error) {
	dbOrderItems, err := r.queries.GetOrderItemsByOrderID(ctx, int32(orderID))
	if err != nil {
		return nil, err
	}

	return r.dbOrderItemsToEntities(dbOrderItems)
}

func (r *orderItemRepository) DeleteByOrder(ctx context.Context, orderID int) error {
	return r.queries.DeleteOrderItemsByOrderID(ctx, int32(orderID))
}

func (r *orderItemRepository) GetByOrderAndItem(ctx context.Context, orderID, itemID int) (*entity.OrderItem, error) {
	dbOrderItem, err := r.queries.GetOrderItemByOrderAndItem(ctx, sqlc.GetOrderItemByOrderAndItemParams{
		OrderID: int32(orderID),
		ItemID:  int32(itemID),
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return r.dbOrderItemToEntity(dbOrderItem)
}

// Helper methods for conversion
func (r *orderItemRepository) dbOrderItemToEntity(dbItem *sqlc.OrderItem) (*entity.OrderItem, error) {
	money, err := vo.NewMoneyFromBaht(utils.FromPgNumericToFloat(dbItem.UnitPrice))
	if err != nil {
		return nil, err
	}

	return &entity.OrderItem{
		ID:        int(dbItem.ID),
		OrderID:   int(dbItem.OrderID),
		ItemID:    int(dbItem.ItemID),
		Quantity:  int(dbItem.Quantity),
		UnitPrice: money,
		Notes:     utils.FromPgTextToString(dbItem.Notes),
		CreatedAt: dbItem.CreatedAt.Time,
		UpdatedAt: dbItem.UpdatedAt.Time,
		Name:      utils.FromPgTextToString(dbItem.Name),
	}, nil
}

func (r *orderItemRepository) dbOrderItemsToEntities(dbItems []*sqlc.OrderItem) ([]*entity.OrderItem, error) {
	entities := make([]*entity.OrderItem, len(dbItems))
	for i, dbItem := range dbItems {
		entity, err := r.dbOrderItemToEntity(dbItem)
		if err != nil {
			return nil, err
		}
		entities[i] = entity
	}
	return entities, nil
}
