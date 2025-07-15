// cmd/main.go (Updated with Revenue functionality)
package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/hydr0g3nz/poc_pos_restuarant/config"
	"github.com/hydr0g3nz/poc_pos_restuarant/internal/adapter/controller"
	sqlcRepo "github.com/hydr0g3nz/poc_pos_restuarant/internal/adapter/repository/sqlc"
	usecase "github.com/hydr0g3nz/poc_pos_restuarant/internal/application"
	"github.com/hydr0g3nz/poc_pos_restuarant/internal/domain/service"
	"github.com/hydr0g3nz/poc_pos_restuarant/internal/infrastructure"
)

func main() {
	// Load configuration
	cfg := config.LoadFromEnv()

	// Setup logger
	logger, err := infrastructure.NewLogger(cfg.IsProduction())
	if err != nil {
		log.Fatalf("Failed to create logger: %v", err)
	}
	defer logger.Sync()

	logger.Info("Starting application")

	// Setup database
	db, err := infrastructure.ConnectDB(&cfg.Database)
	if err != nil {
		logger.Fatal("Failed to connect to database", "error", err)
	}
	defer infrastructure.CloseDB(db)

	// Setup cache
	cache := infrastructure.NewRedisClient(cfg.Cache)
	defer cache.Close()

	// Setup repositories
	userRepo := sqlcRepo.NewUserRepository(db)
	categoryRepo := sqlcRepo.NewCategoryRepository(db)
	menuItemRepo := sqlcRepo.NewMenuItemRepository(db)
	tableRepo := sqlcRepo.NewTableRepository(db)
	orderRepo := sqlcRepo.NewOrderRepository(db)
	orderItemRepo := sqlcRepo.NewOrderItemRepository(db)
	paymentRepo := sqlcRepo.NewPaymentRepository(db)
	revenueRepo := sqlcRepo.NewRevenueRepository(db) // New revenue repository

	// Setup domain services
	orderService := service.NewOrderService(orderRepo, orderItemRepo, tableRepo, menuItemRepo)
	qrCodeService := service.NewQRCodeService(tableRepo)
	// revenueService := service.NewRevenueService(revenueRepo, paymentRepo, orderRepo) // New revenue service

	// Setup use cases
	userUsecase := usecase.NewUserUsecase(userRepo, logger, cfg)
	categoryUsecase := usecase.NewCategoryUsecase(categoryRepo, logger, cfg)
	menuItemUsecase := usecase.NewMenuItemUsecase(menuItemRepo, categoryRepo, logger, cfg)
	tableUsecase := usecase.NewTableUsecase(tableRepo, logger, cfg)
	orderUsecase := usecase.NewOrderUsecase(orderRepo, orderItemRepo, tableRepo, menuItemRepo, orderService, logger, cfg)
	paymentUsecase := usecase.NewPaymentUsecase(paymentRepo, orderRepo, orderService, logger, cfg)
	qrCodeUsecase := usecase.NewQRCodeUsecase(tableRepo, orderRepo, qrCodeService, orderUsecase, logger, cfg)
	revenueUsecase := usecase.NewRevenueUsecase(revenueRepo, paymentRepo, orderRepo, logger, cfg) // New revenue usecase

	// Setup controllers
	userController := controller.NewUserController(userUsecase)
	categoryController := controller.NewCategoryController(categoryUsecase)
	menuItemController := controller.NewMenuItemController(menuItemUsecase)
	tableController := controller.NewTableController(tableUsecase, qrCodeUsecase)
	orderController := controller.NewOrderController(orderUsecase)
	paymentController := controller.NewPaymentController(paymentUsecase)
	revenueController := controller.NewRevenueController(revenueUsecase) // New revenue controller

	// Setup fiber server
	app := infrastructure.NewFiber(infrastructure.ServerConfig{
		Address:      cfg.Server.Port,
		ReadTimeout:  time.Duration(cfg.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.Server.WriteTimeout) * time.Second,
	})

	// Register routes
	api := app.Group("/api/v1")
	userController.RegisterRoutes(api)
	categoryController.RegisterRoutes(api)
	menuItemController.RegisterRoutes(api)
	tableController.RegisterRoutes(api)
	orderController.RegisterRoutes(api)
	paymentController.RegisterRoutes(api)
	revenueController.RegisterRoutes(api) // Register revenue routes

	// Graceful shutdown
	go func() {
		logger.Info("Server starting", "port", cfg.Server.Port)
		if err := app.Listen(":" + cfg.Server.Port); err != nil {
			logger.Fatal("Failed to start server", "error", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	// Shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := app.ShutdownWithContext(ctx); err != nil {
		logger.Error("Server forced to shutdown", "error", err)
	}

	logger.Info("Server exited")
}

// NewFiberServer creates a new Fiber server instance
func NewFiberServer(config *config.Config) *infrastructure.FiberApp {
	return infrastructure.NewFiber(infrastructure.ServerConfig{
		Address:      config.Server.Port,
		ReadTimeout:  time.Duration(config.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(config.Server.WriteTimeout) * time.Second,
		IdleTimeout:  time.Duration(config.Server.ReadTimeout) * time.Second,
	})
}
