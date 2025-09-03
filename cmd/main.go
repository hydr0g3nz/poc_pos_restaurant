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
	mockAdapter "github.com/hydr0g3nz/poc_pos_restuarant/internal/adapter/mock"
	"github.com/hydr0g3nz/poc_pos_restuarant/internal/adapter/presenter"
	gormRepo "github.com/hydr0g3nz/poc_pos_restuarant/internal/adapter/repository/gorm"
	migrater "github.com/hydr0g3nz/poc_pos_restuarant/internal/adapter/repository/migration"
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
	db, err := infrastructure.ConnectGorm(&cfg.Database)
	if err != nil {
		logger.Fatal("Failed to connect to database", "error", err)
	}
	defer infrastructure.CloseGorm(db)
	if err := migrater.AutoMigrate(db); err != nil {
		logger.Fatal("Failed to migrate database", "error", err)
	}
	// Setup cache
	cache := infrastructure.NewRedisClient(cfg.Cache)
	defer cache.Close()
	// printerService, err := infrastructure.NewPrinterService(cfg.Printer.URL)
	// if err != nil {
	// 	logger.Fatal("Failed to connect to printer service", "error", err)
	// }
	// defer printerService.Close()
	printerMock := mockAdapter.NewPrinterService()
	// infra
	qrcodeGenerator := infrastructure.NewQRCodeService()

	errorPresenter := presenter.NewErrorPresenter(logger)
	// Setup repositories
	repoContainer := gormRepo.NewRepositoryContainer(db)
	userRepo := repoContainer.UserRepository()
	categoryRepo := repoContainer.CategoryRepository()
	menuItemRepo := repoContainer.MenuItemRepository()
	tableRepo := repoContainer.TableRepository()
	orderRepo := repoContainer.OrderRepository()
	orderItemRepo := repoContainer.OrderItemRepository()
	paymentRepo := repoContainer.PaymentRepository()
	revenueRepo := repoContainer.RevenueRepository() // New revenue repository
	kitchenStationRepo := repoContainer.KitchenStationRepository()
	orderItemOptionRepo := repoContainer.OrderItemOptionRepository()
	menuOptionRepo := repoContainer.MenuOptionRepository()
	optionValueRepo := repoContainer.OptionValueRepository()
	txManager := repoContainer.TxManager()
	// menuItemOptionRepo := repoContainer.MenuItemOptionRepository()

	// Setup domain services
	orderService := service.NewOrderService(
		orderRepo,
		orderItemRepo,
		orderItemOptionRepo,
		menuOptionRepo,
		optionValueRepo,
		tableRepo,
		menuItemRepo,
	)
	qrCodeService := service.NewQRCodeService(cfg.App.QRcodeURL, qrcodeGenerator, orderRepo) // New QR code service (pass in "tableRepo)
	// revenueService := service.NewRevenueService(revenueRepo, paymentRepo, orderRepo) // New revenue service

	// Setup use cases
	userUsecase := usecase.NewUserUsecase(userRepo, logger, cfg)
	categoryUsecase := usecase.NewCategoryUsecase(categoryRepo, logger, cfg)
	menuItemUsecase := usecase.NewMenuItemUsecase(menuItemRepo, categoryRepo, kitchenStationRepo, logger, cfg)
	tableUsecase := usecase.NewTableUsecase(tableRepo, logger, cfg)
	orderItemOptionUsecase := usecase.NewOrderItemOptionUsecase(orderItemOptionRepo, orderItemRepo, menuOptionRepo, optionValueRepo, orderRepo, logger, cfg)
	orderUsecase := usecase.NewOrderUsecase(
		orderItemOptionUsecase,
		orderRepo,
		orderItemRepo,
		tableRepo,
		menuItemRepo,
		orderService,
		qrCodeService,
		printerMock,
		// printerService,
		txManager,
		logger, cfg)
	paymentUsecase := usecase.NewPaymentUsecase(paymentRepo, orderRepo, orderService, logger, cfg)
	// qrCodeUsecase := usecase.NewQRCodeUsecase(tableRepo, orderRepo, qrCodeService, orderUsecase, logger, cfg)
	revenueUsecase := usecase.NewRevenueUsecase(revenueRepo, paymentRepo, orderRepo, logger, cfg) // New revenue usecase
	kitchenUsecase := usecase.NewKitchenUsecase(orderItemRepo, orderRepo, menuItemRepo, tableRepo, orderItemOptionRepo, menuOptionRepo, optionValueRepo, logger, cfg)
	kitchenStationUsecase := usecase.NewKitchenStationUsecase(kitchenStationRepo, logger, cfg)
	// menuOptionUsecase := usecase.NewMenuOptionUsecase(menuOptionRepo, logger, cfg)
	menuWithOptionsUsecase := usecase.NewMenuWithOptionsUsecase(repoContainer)
	menuOptionMgmtUsecase := usecase.NewMenuOptionManagementUsecase(repoContainer)
	// Setup controllers
	userController := controller.NewUserController(userUsecase, errorPresenter)
	categoryController := controller.NewCategoryController(categoryUsecase, errorPresenter)
	menuItemController := controller.NewMenuItemController(menuItemUsecase, errorPresenter)
	tableController := controller.NewTableController(tableUsecase, errorPresenter)
	orderController := controller.NewOrderController(orderUsecase, errorPresenter)
	paymentController := controller.NewPaymentController(paymentUsecase, errorPresenter)
	revenueController := controller.NewRevenueController(revenueUsecase, errorPresenter) // New revenue controller
	kitchenController := controller.NewKitchenController(kitchenUsecase, kitchenStationUsecase, errorPresenter)
	customController := controller.NewCustomerController(categoryUsecase, menuItemUsecase, orderUsecase, errorPresenter)
	// menuOptionController := controller.NewMenuOptionController(menuOptionUsecase, errorPresenter)
	menuOptionController := controller.NewMenuWithOptionsController(menuWithOptionsUsecase, menuOptionMgmtUsecase, errorPresenter)

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
	kitchenController.RegisterRoutes(api)
	customController.RegisterRoutes(api)
	menuOptionController.RegisterRoutes(api)

	// Graceful shutdown
	go func() {
		logger.Info("Server starting", "port", cfg.Server.Port)
		if err := app.Listen(cfg.Server.Host + ":" + cfg.Server.Port); err != nil {
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
