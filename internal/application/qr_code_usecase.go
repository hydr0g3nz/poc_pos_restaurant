package usecase

// import (
// 	"github.com/hydr0g3nz/poc_pos_restuarant/config"
// 	"github.com/hydr0g3nz/poc_pos_restuarant/internal/domain/infra"
// 	"github.com/hydr0g3nz/poc_pos_restuarant/internal/domain/repository"
// 	"github.com/hydr0g3nz/poc_pos_restuarant/internal/domain/service"
// )

// // qrCodeUsecase implements QRCodeUsecase interface
// type qrCodeUsecase struct {
// 	tableRepo     repository.TableRepository
// 	orderRepo     repository.OrderRepository
// 	qrCodeService service.QRCodeService
// 	orderUsecase  OrderUsecase
// 	logger        infra.Logger
// 	config        *config.Config
// }

// // NewQRCodeUsecase creates a new QR code usecase
// func NewQRCodeUsecase(
// 	tableRepo repository.TableRepository,
// 	orderRepo repository.OrderRepository,
// 	qrCodeService service.QRCodeService,
// 	orderUsecase OrderUsecase,
// 	logger infra.Logger,
// 	config *config.Config,
// ) QRCodeUsecase {
// 	return &qrCodeUsecase{
// 		tableRepo:     tableRepo,
// 		orderRepo:     orderRepo,
// 		qrCodeService: qrCodeService,
// 		orderUsecase:  orderUsecase,
// 		logger:        logger,
// 		config:        config,
// 	}
// }

// // ScanQRCode scans QR code and returns table information with order status
// // func (u *qrCodeUsecase) ScanQRCode(ctx context.Context, qrCode string) (*QRCodeScanResponse, error) {
// // 	u.logger.Info("Scanning QR code", "qrCode", qrCode)

// // 	// Validate and get table from QR code
// // 	table, err := u.qrCodeService.ValidateQRCode(ctx, qrCode)
// // 	if err != nil {
// // 		u.logger.Error("Error validating QR code", "error", err, "qrCode", qrCode)
// // 		return nil, err
// // 	}

// // 	// Check if table is active
// // 	if !table.IsActive {
// // 		u.logger.Warn("Table is not active", "tableID", table.ID)
// // 		return nil, errs.ErrTableNotAvailable
// // 	}

// // 	// Check if table has open order
// // 	openOrder, err := u.orderRepo.GetOpenOrderByTable(ctx, table.ID)
// // 	if err != nil {
// // 		u.logger.Error("Error checking open order", "error", err, "tableID", table.ID)
// // 		return nil, fmt.Errorf("failed to check open order: %w", err)
// // 	}

// // 	response := &QRCodeScanResponse{
// // 		TableID: table.ID,
// // 		Table: &TableResponse{
// // 			ID:          table.ID,
// // 			TableNumber: table.TableNumber.Number(),
// // 			QRCode:      table.QRCode,
// // 			Seating:     table.Seating,
// // 		},
// // 		HasOpenOrder: openOrder != nil,
// // 	}

// // 	// If there's an open order, include it in response
// // 	if openOrder != nil {
// // 		response.OpenOrder = &OrderResponse{
// // 			ID:        openOrder.ID,
// // 			TableID:   openOrder.TableID,
// // 			Status:    openOrder.Status.String(),
// // 			CreatedAt: openOrder.CreatedAt,
// // 		}
// // 		if openOrder.ClosedAt != nil {
// // 			response.OpenOrder.ClosedAt = openOrder.ClosedAt
// // 		}

// // 		u.logger.Info("QR code scanned - table has open order", "tableID", table.ID, "orderID", openOrder.ID)
// // 	} else {
// // 		u.logger.Info("QR code scanned - table is available", "tableID", table.ID)
// // 	}

// // 	return response, nil
// // }

// // // CreateOrderFromQRCode creates a new order from QR code scan
// // func (u *qrCodeUsecase) CreateOrderFromQRCode(ctx context.Context, qrCode string) (*OrderResponse, error) {
// // 	u.logger.Info("Creating order from QR code", "qrCode", qrCode)

// // 	// First scan QR code to get table info and check availability
// // 	scanResult, err := u.ScanQRCode(ctx, qrCode)
// // 	if err != nil {
// // 		u.logger.Error("Error scanning QR code", "error", err, "qrCode", qrCode)
// // 		return nil, err
// // 	}

// // 	// Check if table already has open order
// // 	if scanResult.HasOpenOrder {
// // 		u.logger.Warn("Table already has open order", "tableID", scanResult.TableID, "orderID", scanResult.OpenOrder.ID)
// // 		return nil, errs.ErrTableAlreadyHasOpenOrder
// // 	}

// // 	// Create new order for the table
// // 	createOrderReq := &CreateOrderRequest{
// // 		TableID: scanResult.TableID,
// // 	}

// // 	order, err := u.orderUsecase.CreateOrder(ctx, createOrderReq)
// // 	if err != nil {
// // 		u.logger.Error("Error creating order from QR code", "error", err, "tableID", scanResult.TableID)
// // 		return nil, err
// // 	}

// // 	u.logger.Info("Order created successfully from QR code", "orderID", order.ID, "tableID", scanResult.TableID)

// // 	return order, nil
// // }
