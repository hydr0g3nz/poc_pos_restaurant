// internal/adapter/repository/table_repository.go
package repository

import (
	"context"

	"github.com/hydr0g3nz/poc_pos_restuarant/internal/adapter/repository/gorm/model"
	"github.com/hydr0g3nz/poc_pos_restuarant/internal/domain/entity"
	"github.com/hydr0g3nz/poc_pos_restuarant/internal/domain/repository"
	"github.com/hydr0g3nz/poc_pos_restuarant/internal/domain/vo"
	"gorm.io/gorm"
)

type tableRepository struct {
	baseRepository
}

func NewTableRepository(db *gorm.DB) repository.TableRepository {
	return &tableRepository{
		baseRepository: baseRepository{db: db},
	}
}

func (r *tableRepository) Create(ctx context.Context, table *entity.Table) (*entity.Table, error) {
	dbTable := r.entityToModel(table)

	if err := r.db.WithContext(ctx).Create(dbTable).Error; err != nil {
		return nil, err
	}

	return r.modelToEntity(dbTable), nil
}

func (r *tableRepository) GetByID(ctx context.Context, id int) (*entity.Table, error) {
	var dbTable model.Table

	if err := r.db.WithContext(ctx).First(&dbTable, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return r.modelToEntity(&dbTable), nil
}

func (r *tableRepository) GetByNumber(ctx context.Context, number int) (*entity.Table, error) {
	var dbTable model.Table

	if err := r.db.WithContext(ctx).Where("table_number = ?", number).First(&dbTable).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return r.modelToEntity(&dbTable), nil
}

func (r *tableRepository) GetByQRCode(ctx context.Context, qrCode string) (*entity.Table, error) {
	var dbOrder model.Order

	if err := r.db.WithContext(ctx).Preload("Table").Where("qr_code = ?", qrCode).First(&dbOrder).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return nil, nil
}

func (r *tableRepository) Update(ctx context.Context, table *entity.Table) (*entity.Table, error) {
	dbTable := r.entityToModel(table)

	if err := r.db.WithContext(ctx).Save(dbTable).Error; err != nil {
		return nil, err
	}

	return r.modelToEntity(dbTable), nil
}

func (r *tableRepository) Delete(ctx context.Context, id int) error {
	return r.db.WithContext(ctx).Delete(&model.Table{}, id).Error
}

func (r *tableRepository) List(ctx context.Context) ([]*entity.Table, error) {
	var dbTables []model.Table

	if err := r.db.WithContext(ctx).Preload("CurrentOrder", "order_status IN (?)", []string{vo.OrderStatusOpen.String(), vo.OrderStatusOrdered.String()}).Find(&dbTables).Error; err != nil {
		return nil, err
	}

	return r.modelsToEntities(dbTables), nil
}

func (r *tableRepository) HasOrders(ctx context.Context, tableID int) (bool, error) {
	var count int64

	if err := r.db.WithContext(ctx).Model(&model.Order{}).Where("table_id = ?", tableID).Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}

// Helper methods
func (r *tableRepository) entityToModel(table *entity.Table) *model.Table {
	return &model.Table{
		ID:          table.ID,
		TableNumber: table.TableNumber,
		Seating:     table.Seating,
		IsActive:    table.IsActive,
	}
}

func (r *tableRepository) modelToEntity(dbTable *model.Table) *entity.Table {
	return &entity.Table{
		ID:          dbTable.ID,
		TableNumber: dbTable.TableNumber,
		Seating:     dbTable.Seating,
		IsActive:    dbTable.IsActive,
		CurrentOrder: func() *entity.OrderTableDetails {
			if dbTable.CurrentOrder.ID == 0 {
				return nil
			}
			return &entity.OrderTableDetails{
				OrderID:     dbTable.CurrentOrder.ID,
				OrderNumber: dbTable.CurrentOrder.OrderNumber,
				Status:      dbTable.CurrentOrder.OrderStatus,
				QRCode:      dbTable.CurrentOrder.QRCode,
				CreatedAt:   dbTable.CurrentOrder.CreatedAt.Format("2006-01-02 15:04:05"),
			}
		}(),
	}
}

func (r *tableRepository) modelsToEntities(dbTables []model.Table) []*entity.Table {
	entities := make([]*entity.Table, len(dbTables))
	for i, dbTable := range dbTables {
		entities[i] = r.modelToEntity(&dbTable)
	}
	return entities
}
