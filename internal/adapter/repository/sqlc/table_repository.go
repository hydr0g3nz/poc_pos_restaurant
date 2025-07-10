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

type tableRepository struct {
	db      *pgxpool.Pool
	queries *sqlc.Queries
}

func NewTableRepository(db *pgxpool.Pool) repository.TableRepository {
	return &tableRepository{
		db:      db,
		queries: sqlc.New(db),
	}
}

func (r *tableRepository) Create(ctx context.Context, table *entity.Table) (*entity.Table, error) {
	dbTable, err := r.queries.CreateTable(ctx, sqlc.CreateTableParams{
		TableNumber: int32(table.TableNumber.Number()),
		QrCode:      table.QRCode,
		Seating:     int32(table.Seating),
		IsActive:    utils.ConvertToPGBool(table.IsActive),
	})
	if err != nil {
		return nil, err
	}

	return r.dbTableToEntity(dbTable)
}

func (r *tableRepository) GetByID(ctx context.Context, id int) (*entity.Table, error) {
	dbTable, err := r.queries.GetTableByID(ctx, int32(id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return r.dbTableToEntity(dbTable)
}

func (r *tableRepository) GetByNumber(ctx context.Context, number int) (*entity.Table, error) {
	dbTable, err := r.queries.GetTableByNumber(ctx, int32(number))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return r.dbTableToEntity(dbTable)
}

func (r *tableRepository) GetByQRCode(ctx context.Context, qrCode string) (*entity.Table, error) {
	dbTable, err := r.queries.GetTableByQRCode(ctx, qrCode)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return r.dbTableToEntity(dbTable)
}

func (r *tableRepository) Update(ctx context.Context, table *entity.Table) (*entity.Table, error) {
	dbTable, err := r.queries.UpdateTable(ctx, sqlc.UpdateTableParams{
		ID:          int32(table.ID),
		TableNumber: int32(table.TableNumber.Number()),
		QrCode:      table.QRCode,
		Seating:     int32(table.Seating),
		IsActive:    utils.ConvertToPGBool(table.IsActive),
	})
	if err != nil {
		return nil, err
	}

	return r.dbTableToEntity(dbTable)
}

func (r *tableRepository) Delete(ctx context.Context, id int) error {
	return r.queries.DeleteTable(ctx, int32(id))
}

func (r *tableRepository) List(ctx context.Context) ([]*entity.Table, error) {
	dbTables, err := r.queries.ListTables(ctx)
	if err != nil {
		return nil, err
	}

	return r.dbTablesToEntities(dbTables)
}

func (r *tableRepository) HasOrders(ctx context.Context, tableID int) (bool, error) {
	result, err := r.queries.CheckTableHasOrders(ctx, int32(tableID))
	if err != nil {
		return false, err
	}
	return result, nil
}

// Helper methods for conversion
func (r *tableRepository) dbTableToEntity(dbTable *sqlc.Table) (*entity.Table, error) {
	tableNumber, err := vo.NewTableNumber(int(dbTable.TableNumber))
	if err != nil {
		return nil, err
	}

	return &entity.Table{
		ID:          int(dbTable.ID),
		TableNumber: tableNumber,
		QRCode:      dbTable.QrCode,
		Seating:     int(dbTable.Seating),
		IsActive:    utils.ConvertToBool(dbTable.IsActive),
	}, nil
}

func (r *tableRepository) dbTablesToEntities(dbTables []*sqlc.Table) ([]*entity.Table, error) {
	entities := make([]*entity.Table, len(dbTables))
	for i, dbTable := range dbTables {
		entity, err := r.dbTableToEntity(dbTable)
		if err != nil {
			return nil, err
		}
		entities[i] = entity
	}
	return entities, nil
}
