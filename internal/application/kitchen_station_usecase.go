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

// KitchenStationUsecase interface defines the contract for kitchen station business logic

// kitchenStationUsecase implements KitchenStationUsecase interface
type kitchenStationUsecase struct {
	kitchenStationRepo repository.KitchenStationRepository
	logger             infra.Logger
	config             *config.Config
}

// NewKitchenStationUsecase creates a new kitchen station usecase
func NewKitchenStationUsecase(
	kitchenStationRepo repository.KitchenStationRepository,
	logger infra.Logger,
	config *config.Config,
) KitchenStationUsecase {
	return &kitchenStationUsecase{
		kitchenStationRepo: kitchenStationRepo,
		logger:             logger,
		config:             config,
	}
}

// CreateKitchenStation creates a new kitchen station
func (u *kitchenStationUsecase) CreateKitchenStation(ctx context.Context, req *CreateKitchenStationRequest) (*KitchenStationOnlyResponse, error) {
	u.logger.Info("Creating kitchen station", "name", req.Name)

	// Create kitchen station entity
	kitchenStation, err := entity.NewKitchenStation(req.Name, req.IsAvailable)
	if err != nil {
		u.logger.Error("Error creating kitchen station entity", "error", err, "name", req.Name)
		return nil, err
	}

	// Save to database
	createdStation, err := u.kitchenStationRepo.Create(ctx, kitchenStation)
	if err != nil {
		u.logger.Error("Error creating kitchen station", "error", err, "name", req.Name)
		return nil, fmt.Errorf("failed to create kitchen station: %w", err)
	}

	u.logger.Info("Kitchen station created successfully", "stationID", createdStation.ID, "name", createdStation.Name)

	return u.toKitchenStationResponse(createdStation), nil
}

// GetKitchenStation retrieves kitchen station by ID
func (u *kitchenStationUsecase) GetKitchenStation(ctx context.Context, id int) (*KitchenStationOnlyResponse, error) {
	u.logger.Debug("Getting kitchen station", "stationID", id)

	kitchenStation, err := u.kitchenStationRepo.GetByID(ctx, id)
	if err != nil {
		u.logger.Error("Error getting kitchen station", "error", err, "stationID", id)
		return nil, fmt.Errorf("failed to get kitchen station: %w", err)
	}
	if kitchenStation == nil {
		u.logger.Warn("Kitchen station not found", "stationID", id)
		return nil, errs.NewNotFoundError("kitchen station", id)
	}

	return u.toKitchenStationResponse(kitchenStation), nil
}

// UpdateKitchenStation updates kitchen station information
func (u *kitchenStationUsecase) UpdateKitchenStation(ctx context.Context, id int, req *UpdateKitchenStationRequest) (*KitchenStationOnlyResponse, error) {
	u.logger.Info("Updating kitchen station", "stationID", id, "name", req.Name)

	// Get current kitchen station
	currentStation, err := u.kitchenStationRepo.GetByID(ctx, id)
	if err != nil {
		u.logger.Error("Error getting current kitchen station", "error", err, "stationID", id)
		return nil, fmt.Errorf("failed to get kitchen station: %w", err)
	}
	if currentStation == nil {
		return nil, errs.NewNotFoundError("kitchen station", id)
	}

	// Update fields
	currentStation.Name = req.Name
	currentStation.IsAvailable = req.IsAvailable
	// Update kitchen station
	updatedStation, err := u.kitchenStationRepo.Update(ctx, currentStation)
	if err != nil {
		u.logger.Error("Error updating kitchen station", "error", err, "stationID", id)
		return nil, fmt.Errorf("failed to update kitchen station: %w", err)
	}

	u.logger.Info("Kitchen station updated successfully", "stationID", id)

	return u.toKitchenStationResponse(updatedStation), nil
}

// DeleteKitchenStation deletes a kitchen station
func (u *kitchenStationUsecase) DeleteKitchenStation(ctx context.Context, id int) error {
	u.logger.Info("Deleting kitchen station", "stationID", id)

	// Check if kitchen station exists
	kitchenStation, err := u.kitchenStationRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get kitchen station: %w", err)
	}
	if kitchenStation == nil {
		return errs.NewNotFoundError("kitchen station", id)
	}

	// Delete kitchen station
	if err := u.kitchenStationRepo.Delete(ctx, id); err != nil {
		u.logger.Error("Error deleting kitchen station", "error", err, "stationID", id)
		return fmt.Errorf("failed to delete kitchen station: %w", err)
	}

	u.logger.Info("Kitchen station deleted successfully", "stationID", id)
	return nil
}

// ListKitchenStations retrieves all kitchen stations
func (u *kitchenStationUsecase) ListKitchenStations(ctx context.Context, onlyAvailable bool) ([]*KitchenStationOnlyResponse, error) {
	u.logger.Debug("Listing kitchen stations")

	kitchenStations, err := u.kitchenStationRepo.List(ctx, onlyAvailable, 1000, 0) // Get all stations
	if err != nil {
		u.logger.Error("Error listing kitchen stations", "error", err)
		return nil, fmt.Errorf("failed to list kitchen stations: %w", err)
	}

	return u.toKitchenStationResponses(kitchenStations), nil
}

// Helper methods

// toKitchenStationResponse converts entity to response
func (u *kitchenStationUsecase) toKitchenStationResponse(station *entity.KitchenStation) *KitchenStationOnlyResponse {
	return &KitchenStationOnlyResponse{
		ID:          station.ID,
		Name:        station.Name,
		IsAvailable: station.IsAvailable,
	}
}

// toKitchenStationResponses converts slice of entities to responses
func (u *kitchenStationUsecase) toKitchenStationResponses(stations []*entity.KitchenStation) []*KitchenStationOnlyResponse {
	responses := make([]*KitchenStationOnlyResponse, len(stations))
	for i, station := range stations {
		responses[i] = u.toKitchenStationResponse(station)
	}
	return responses
}
