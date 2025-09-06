package usecase

import (
	"context"
	"fmt"

	"github.com/hydr0g3nz/poc_pos_restuarant/config"
	"github.com/hydr0g3nz/poc_pos_restuarant/internal/domain/entity"
	errs "github.com/hydr0g3nz/poc_pos_restuarant/internal/domain/error"
	"github.com/hydr0g3nz/poc_pos_restuarant/internal/domain/infra"
	"github.com/hydr0g3nz/poc_pos_restuarant/internal/domain/repository"
)

// tableUsecase implements TableUsecase interface
type tableUsecase struct {
	tableRepo repository.TableRepository
	logger    infra.Logger
	config    *config.Config
}

// NewTableUsecase creates a new table usecase
func NewTableUsecase(
	tableRepo repository.TableRepository,
	logger infra.Logger,
	config *config.Config,
) TableUsecase {
	return &tableUsecase{
		tableRepo: tableRepo,
		logger:    logger,
		config:    config,
	}
}

// CreateTable creates a new table
func (u *tableUsecase) CreateTable(ctx context.Context, req *CreateTableRequest) (*TableResponse, error) {
	u.logger.Info("Creating table", "tableNumber", req.TableNumber, "seating", req.Seating)

	// Check if table number already exists
	existingTable, err := u.tableRepo.GetByNumber(ctx, req.TableNumber)
	if err != nil {
		u.logger.Error("Error checking existing table", "error", err, "tableNumber", req.TableNumber)
		return nil, fmt.Errorf("failed to check existing table: %w", err)
	}
	if existingTable != nil {
		u.logger.Warn("Table number already exists", "tableNumber", req.TableNumber)
		return nil, errs.ErrDuplicateTableNumber
	}

	// Create table entity
	table, err := entity.NewTable(req.TableNumber, req.Seating)
	if err != nil {
		u.logger.Error("Error creating table entity", "error", err, "tableNumber", req.TableNumber)
		return nil, err
	}

	// Save to database
	createdTable, err := u.tableRepo.Create(ctx, table)
	if err != nil {
		u.logger.Error("Error creating table", "error", err, "tableNumber", req.TableNumber)
		return nil, fmt.Errorf("failed to create table: %w", err)
	}

	u.logger.Info("Table created successfully", "tableID", createdTable.ID, "tableNumber", createdTable.TableNumber)

	return u.toTableResponse(createdTable), nil
}

// GetTable retrieves table by ID
func (u *tableUsecase) GetTable(ctx context.Context, id int) (*TableResponse, error) {
	u.logger.Debug("Getting table", "tableID", id)

	table, err := u.tableRepo.GetByID(ctx, id)
	if err != nil {
		u.logger.Error("Error getting table", "error", err, "tableID", id)
		return nil, fmt.Errorf("failed to get table: %w", err)
	}
	if table == nil {
		u.logger.Warn("Table not found", "tableID", id)
		return nil, errs.ErrTableNotFound
	}

	return u.toTableResponse(table), nil
}

// GetTableByNumber retrieves table by table number
func (u *tableUsecase) GetTableByNumber(ctx context.Context, number int) (*TableResponse, error) {
	u.logger.Debug("Getting table by number", "tableNumber", number)

	table, err := u.tableRepo.GetByNumber(ctx, number)
	if err != nil {
		u.logger.Error("Error getting table by number", "error", err, "tableNumber", number)
		return nil, fmt.Errorf("failed to get table by number: %w", err)
	}
	if table == nil {
		u.logger.Warn("Table not found", "tableNumber", number)
		return nil, errs.ErrTableNotFound
	}

	return u.toTableResponse(table), nil
}

// GetTableByQRCode retrieves table by QR code
func (u *tableUsecase) GetTableByQRCode(ctx context.Context, qrCode string) (*TableResponse, error) {
	u.logger.Debug("Getting table by QR code", "qrCode", qrCode)

	table, err := u.tableRepo.GetByQRCode(ctx, qrCode)
	if err != nil {
		u.logger.Error("Error getting table by QR code", "error", err, "qrCode", qrCode)
		return nil, fmt.Errorf("failed to get table by QR code: %w", err)
	}
	if table == nil {
		u.logger.Warn("Table not found", "qrCode", qrCode)
		return nil, errs.ErrTableNotFound
	}

	return u.toTableResponse(table), nil
}

// UpdateTable updates table information
func (u *tableUsecase) UpdateTable(ctx context.Context, id int, req *UpdateTableRequest) (*TableResponse, error) {
	u.logger.Info("Updating table", "tableID", id, "tableNumber", req.TableNumber, "seating", req.Seating)

	// Get current table
	currentTable, err := u.tableRepo.GetByID(ctx, id)
	if err != nil {
		u.logger.Error("Error getting current table", "error", err, "tableID", id)
		return nil, fmt.Errorf("failed to get table: %w", err)
	}
	if currentTable == nil {
		return nil, errs.ErrTableNotFound
	}

	// Check if new table number is different and unique
	if req.TableNumber != currentTable.TableNumber {
		existingTable, err := u.tableRepo.GetByNumber(ctx, req.TableNumber)
		if err != nil {
			u.logger.Error("Error checking table number uniqueness", "error", err, "tableNumber", req.TableNumber)
			return nil, fmt.Errorf("failed to check table number uniqueness: %w", err)
		}
		if existingTable != nil {
			return nil, errs.ErrDuplicateTableNumber
		}

		// Update table number and regenerate QR code
		newTableNumber, err := entity.NewTable(req.TableNumber, req.Seating)
		if err != nil {
			u.logger.Error("Error creating new table number", "error", err, "tableNumber", req.TableNumber)
			return nil, err
		}
		currentTable.TableNumber = newTableNumber.TableNumber
	}

	// Update seating
	currentTable.Seating = req.Seating
	currentTable.IsActive = req.IsAvailable
	// Update table
	updatedTable, err := u.tableRepo.Update(ctx, currentTable)
	if err != nil {
		u.logger.Error("Error updating table", "error", err, "tableID", id)
		return nil, fmt.Errorf("failed to update table: %w", err)
	}

	u.logger.Info("Table updated successfully", "tableID", id)

	return u.toTableResponse(updatedTable), nil
}

// DeleteTable deletes a table
func (u *tableUsecase) DeleteTable(ctx context.Context, id int) error {
	u.logger.Info("Deleting table", "tableID", id)

	// Check if table exists
	table, err := u.tableRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get table: %w", err)
	}
	if table == nil {
		return errs.ErrTableNotFound
	}

	// Check if table has orders
	hasOrders, err := u.tableRepo.HasOrders(ctx, id)
	if err != nil {
		u.logger.Error("Error checking table orders", "error", err, "tableID", id)
		return fmt.Errorf("failed to check table orders: %w", err)
	}
	if hasOrders {
		u.logger.Warn("Cannot delete table with orders", "tableID", id)
		return errs.ErrCannotDeleteTableWithOrders
	}

	// Delete table
	if err := u.tableRepo.Delete(ctx, id); err != nil {
		u.logger.Error("Error deleting table", "error", err, "tableID", id)
		return fmt.Errorf("failed to delete table: %w", err)
	}

	u.logger.Info("Table deleted successfully", "tableID", id)
	return nil
}

// ListTables retrieves all tables
func (u *tableUsecase) ListTables(ctx context.Context) ([]*TableResponse, error) {
	u.logger.Debug("Listing tables")

	tables, err := u.tableRepo.List(ctx)
	if err != nil {
		u.logger.Error("Error listing tables", "error", err)
		return nil, fmt.Errorf("failed to list tables: %w", err)
	}

	return u.toTableResponses(tables), nil
}

// Helper methods

// toTableResponse converts entity to response
func (u *tableUsecase) toTableResponse(table *entity.Table) *TableResponse {
	t := TableResponse{
		ID:          table.ID,
		TableNumber: table.TableNumber,
		Seating:     table.Seating,
		IsAvailable: table.IsActive,
	}
	if table.CurrentOrder != nil {
		t.CurrentOrder = &OrderTableDetails{
			OrderID:     table.CurrentOrder.OrderID,
			OrderNumber: table.CurrentOrder.OrderNumber,
			Status:      table.CurrentOrder.Status,
			QRCode:      table.CurrentOrder.QRCode,
			CreatedAt:   table.CurrentOrder.CreatedAt,
		}
	}
	fmt.Println("table resp:", t)
	return &t
}

// toTableResponses converts slice of entities to responses
func (u *tableUsecase) toTableResponses(tables []*entity.Table) []*TableResponse {
	responses := make([]*TableResponse, len(tables))
	for i, table := range tables {
		responses[i] = u.toTableResponse(table)
	}
	return responses
}
