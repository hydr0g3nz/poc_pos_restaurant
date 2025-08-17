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

type menuItemRepository struct {
	db      *pgxpool.Pool
	queries *sqlc.Queries
}

func NewMenuItemRepository(db *pgxpool.Pool) repository.MenuItemRepository {
	return &menuItemRepository{
		db:      db,
		queries: sqlc.New(db),
	}
}

func (r *menuItemRepository) Create(ctx context.Context, item *entity.MenuItem) (*entity.MenuItem, error) {
	dbMenuItem, err := r.queries.CreateMenuItem(ctx, sqlc.CreateMenuItemParams{
		CategoryID:  int32(item.CategoryID),
		Name:        item.Name,
		Description: utils.ConvertToText(item.Description),
		Price:       utils.ConvertToPGNumericFromFloat(item.Price.AmountBaht()),
		IsActive:    utils.ConvertToPGBool(item.IsActive),
	})
	if err != nil {
		return nil, err
	}

	return r.dbMenuItemToEntity(dbMenuItem)
}

func (r *menuItemRepository) GetByID(ctx context.Context, id int) (*entity.MenuItem, error) {
	dbMenuItem, err := r.queries.GetMenuItemByID(ctx, int32(id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return r.dbMenuItemToEntity(dbMenuItem)
}

func (r *menuItemRepository) Update(ctx context.Context, item *entity.MenuItem) (*entity.MenuItem, error) {
	dbMenuItem, err := r.queries.UpdateMenuItem(ctx, sqlc.UpdateMenuItemParams{
		ID:          int32(item.ID),
		CategoryID:  int32(item.CategoryID),
		Name:        item.Name,
		Description: utils.ConvertToText(item.Description),
		Price:       utils.ConvertToPGNumericFromFloat(item.Price.AmountBaht()),
		IsActive:    utils.ConvertToPGBool(item.IsActive),
	})
	if err != nil {
		return nil, err
	}

	return r.dbMenuItemToEntity(dbMenuItem)
}

func (r *menuItemRepository) Delete(ctx context.Context, id int) error {
	return r.queries.DeleteMenuItem(ctx, int32(id))
}

func (r *menuItemRepository) List(ctx context.Context, limit, offset int) ([]*entity.MenuItem, error) {
	dbMenuItems, err := r.queries.ListMenuItems(ctx, sqlc.ListMenuItemsParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, err
	}

	return r.dbMenuItemsToEntities(dbMenuItems)
}

func (r *menuItemRepository) ListByCategory(ctx context.Context, categoryID int, limit, offset int) ([]*entity.MenuItem, error) {
	dbMenuItems, err := r.queries.ListMenuItemsByCategory(ctx, sqlc.ListMenuItemsByCategoryParams{
		CategoryID: int32(categoryID),
		Limit:      int32(limit),
		Offset:     int32(offset),
	})
	if err != nil {
		return nil, err
	}

	return r.dbMenuItemsToEntities(dbMenuItems)
}

func (r *menuItemRepository) Search(ctx context.Context, query string, limit, offset int) ([]*entity.MenuItem, error) {
	dbMenuItems, err := r.queries.SearchMenuItems(ctx, sqlc.SearchMenuItemsParams{
		Column1: utils.ConvertToText(query),
		Limit:   int32(limit),
		Offset:  int32(offset),
	})
	if err != nil {
		return nil, err
	}

	return r.dbMenuItemsToEntities(dbMenuItems)
}

// Additional helper methods

// UpdateStatus updates the active status of a menu item
func (r *menuItemRepository) UpdateStatus(ctx context.Context, id int, isActive bool) (*entity.MenuItem, error) {
	dbMenuItem, err := r.queries.UpdateMenuItemStatus(ctx, sqlc.UpdateMenuItemStatusParams{
		ID:       int32(id),
		IsActive: utils.ConvertToPGBool(isActive),
	})
	if err != nil {
		return nil, err
	}

	return r.dbMenuItemToEntity(dbMenuItem)
}

// GetByIDs gets multiple menu items by their IDs
func (r *menuItemRepository) GetByIDs(ctx context.Context, ids []int) ([]*entity.MenuItem, error) {
	// Convert []int to []int32 for the query
	int32IDs := make([]int32, len(ids))
	for i, id := range ids {
		int32IDs[i] = int32(id)
	}

	dbMenuItems, err := r.queries.GetMenuItemsByIDs(ctx, int32IDs)
	if err != nil {
		return nil, err
	}

	return r.dbMenuItemsToEntities(dbMenuItems)
}

// CountByCategory counts menu items in a category
func (r *menuItemRepository) CountByCategory(ctx context.Context, categoryID int) (int, error) {
	count, err := r.queries.CountMenuItemsByCategory(ctx, int32(categoryID))
	if err != nil {
		return 0, err
	}
	return int(count), nil
}

// Helper methods for conversion
func (r *menuItemRepository) dbMenuItemToEntity(dbMenuItem *sqlc.MenuItem) (*entity.MenuItem, error) {

	money, err := vo.NewMoneyFromBaht(utils.FromPgNumericToFloat(dbMenuItem.Price))
	if err != nil {
		return nil, err
	}

	return &entity.MenuItem{
		ID:          int(dbMenuItem.ID),
		CategoryID:  int(dbMenuItem.CategoryID),
		Name:        dbMenuItem.Name,
		Description: utils.FromPgTextToString(dbMenuItem.Description),
		Price:       money,
		IsActive:    utils.ConvertToBool(dbMenuItem.IsActive),
		CreatedAt:   dbMenuItem.CreatedAt.Time,
		UpdatedAt:   dbMenuItem.UpdatedAt.Time,
	}, nil
}

func (r *menuItemRepository) dbMenuItemsToEntities(dbMenuItems []*sqlc.MenuItem) ([]*entity.MenuItem, error) {
	entities := make([]*entity.MenuItem, len(dbMenuItems))
	for i, dbMenuItem := range dbMenuItems {
		entity, err := r.dbMenuItemToEntity(dbMenuItem)
		if err != nil {
			return nil, err
		}
		entities[i] = entity
	}
	return entities, nil
}
