package service

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"

	"codeberg.org/go-pdf/fpdf"
	"github.com/hydr0g3nz/poc_pos_restuarant/internal/domain/entity"
	errs "github.com/hydr0g3nz/poc_pos_restuarant/internal/domain/error"
	"github.com/hydr0g3nz/poc_pos_restuarant/internal/domain/repository"
	"github.com/hydr0g3nz/poc_pos_restuarant/internal/domain/vo"
	"github.com/skip2/go-qrcode"
)

const (
	TAX_RATE      = 0.07 // 7% VAT
	DISCOUNT_RATE = 0.10 // 10% discount
)

// OrderService provides domain logic for orders
type OrderService interface {
	// ValidateOrderCreation validates if order can be created for table
	ValidateOrderCreation(ctx context.Context, tableID int) error

	// CalculateOrderTotal calculates total amount for order
	CalculateOrderTotal(ctx context.Context, order *entity.Order) (vo.Money, error)

	// ValidateOrderItem validates order item before adding
	ValidateOrderItem(ctx context.Context, orderID, itemID int, quantity int) error

	// ProcessOrderClosure processes order closure with validations
	ProcessOrderClosure(ctx context.Context, orderID int) error

	ReceiptPdf(ctx context.Context, order *entity.Order) ([]byte, error)

	QRCodePdf(ctx context.Context, receipt *entity.Order) ([]byte, error)
}

type orderService struct {
	orderRepo           repository.OrderRepository
	orderItemRepo       repository.OrderItemRepository
	orderItemOptionRepo repository.OrderItemOptionRepository
	menuOptionRepo      repository.MenuOptionRepository
	optionValueRepo     repository.OptionValueRepository
	tableRepo           repository.TableRepository
	menuItemRepo        repository.MenuItemRepository
}

func NewOrderService(
	orderRepo repository.OrderRepository,
	orderItemRepo repository.OrderItemRepository,
	orderItemOptionRepo repository.OrderItemOptionRepository,
	menuOptionRepo repository.MenuOptionRepository,
	optionValueRepo repository.OptionValueRepository,
	tableRepo repository.TableRepository,
	menuItemRepo repository.MenuItemRepository,
) OrderService {
	return &orderService{
		orderRepo:           orderRepo,
		orderItemRepo:       orderItemRepo,
		orderItemOptionRepo: orderItemOptionRepo,
		menuOptionRepo:      menuOptionRepo,
		optionValueRepo:     optionValueRepo,
		tableRepo:           tableRepo,
		menuItemRepo:        menuItemRepo,
	}
}

func (s *orderService) ValidateOrderCreation(ctx context.Context, tableID int) error {
	// Check if table exists
	table, err := s.tableRepo.GetByID(ctx, tableID)
	if err != nil {
		return fmt.Errorf("failed to get table: %w", err)
	}
	if table == nil {
		return errs.ErrTableNotFound
	}

	// Check if table already has an open order
	openOrder, err := s.orderRepo.GetOpenOrderByTable(ctx, tableID)
	if err != nil {
		return fmt.Errorf("failed to check open order: %w", err)
	}
	if openOrder != nil {
		return errs.ErrTableAlreadyHasOpenOrder
	}

	return nil
}

func (s *orderService) CalculateOrderTotal(ctx context.Context, order *entity.Order) (vo.Money, error) {
	if order == nil {
		return vo.Money{}, errs.ErrOrderNotFound
	}

	// Get order items if not loaded
	if order.Items == nil {
		items, err := s.orderItemRepo.ListByOrder(ctx, order.ID)
		if err != nil {
			return vo.Money{}, fmt.Errorf("failed to get order items: %w", err)
		}
		order.Items = items
	}

	// Calculate total including options
	total, _ := vo.NewMoneyFromSatang(0)
	for _, item := range order.Items {
		itemSubtotal := item.CalculateSubtotal()

		// Add option prices
		options, err := s.orderItemOptionRepo.GetByOrderItemID(ctx, item.ID)
		if err == nil { // Don't fail if options can't be loaded
			for _, option := range options {
				itemSubtotal = itemSubtotal.Add(option.AdditionalPrice.Multiply(float64(item.Quantity)))
			}
		}

		total = total.Add(itemSubtotal)
	}

	return total, nil
}

func (s *orderService) ValidateOrderItem(ctx context.Context, orderID, itemID int, quantity int) error {
	// Check if order exists and is open
	order, err := s.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return fmt.Errorf("failed to get order: %w", err)
	}
	if order == nil {
		return errs.ErrOrderNotFound
	}
	if !order.IsOpen() {
		return errs.ErrOrderNotOpen
	}

	// Check if menu item exists
	menuItem, err := s.menuItemRepo.GetByID(ctx, itemID)
	if err != nil {
		return fmt.Errorf("failed to get menu item: %w", err)
	}
	if menuItem == nil {
		return errs.ErrMenuItemNotFound
	}

	// Validate quantity
	if quantity <= 0 {
		return errs.ErrInvalidQuantity
	}

	return nil
}

func (s *orderService) ProcessOrderClosure(ctx context.Context, orderID int) error {
	// Check if order exists
	order, err := s.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return fmt.Errorf("failed to get order: %w", err)
	}
	if order == nil {
		return errs.ErrOrderNotFound
	}

	// Check if order is already closed
	if order.IsClosed() {
		return errs.ErrOrderAlreadyClosed
	}

	// Check if order has items
	items, err := s.orderItemRepo.ListByOrder(ctx, orderID)
	if err != nil {
		return fmt.Errorf("failed to get order items: %w", err)
	}
	if len(items) == 0 {
		return errs.ErrEmptyOrder
	}

	return nil
}
func (s *orderService) ReceiptPdf(ctx context.Context, order *entity.Order) ([]byte, error) {
	if order == nil {
		return nil, errs.ErrOrderNotFound
	}
	w := &bytes.Buffer{}
	if err := s.generateReceiptPDF(order, w); err != nil {
		return nil, fmt.Errorf("failed to generate receipt PDF: %w", err)
	}
	return w.Bytes(), nil
}
func (s *orderService) QRCodePdf(ctx context.Context, receipt *entity.Order) ([]byte, error) {
	if receipt == nil {
		return nil, errs.ErrOrderNotFound
	}
	w := &bytes.Buffer{}
	if err := generateOrderQRCodePDF(receipt, w); err != nil {
		return nil, fmt.Errorf("failed to generate QR code PDF: %w", err)
	}
	return w.Bytes(), nil
}

// Helper struct to hold item with its options
type ItemWithOptions struct {
	Item    *entity.OrderItem
	Options []OptionWithDetails
}

type OptionWithDetails struct {
	Option          *entity.MenuOption
	Value           *entity.OptionValue
	AdditionalPrice vo.Money
}

// loadOrderItemsWithOptions loads order items and their options
func (s *orderService) loadOrderItemsWithOptions(ctx context.Context, order *entity.Order) error {
	if order.Items == nil {
		items, err := s.orderItemRepo.ListByOrder(ctx, order.ID)
		if err != nil {
			return fmt.Errorf("failed to get order items: %w", err)
		}
		order.Items = items
	}
	return nil
}

// getItemOptions retrieves options for a specific order item
func (s *orderService) getItemOptions(ctx context.Context, itemID int) ([]OptionWithDetails, error) {
	var optionDetails []OptionWithDetails

	// Get order item options
	itemOptions, err := s.orderItemOptionRepo.GetByOrderItemID(ctx, itemID)
	if err != nil {
		return optionDetails, err
	}

	for _, itemOption := range itemOptions {
		// Get option details
		option, err := s.menuOptionRepo.GetByID(ctx, itemOption.OptionID)
		if err != nil {
			continue // Skip if option not found
		}

		// Get value details
		value, err := s.optionValueRepo.GetByID(ctx, itemOption.ValueID)
		if err != nil {
			continue // Skip if value not found
		}

		optionDetails = append(optionDetails, OptionWithDetails{
			Option:          option,
			Value:           value,
			AdditionalPrice: itemOption.AdditionalPrice,
		})
	}

	return optionDetails, nil
}

func (s *orderService) generateReceiptPDF(order *entity.Order, writer io.Writer) error {
	pdf := fpdf.NewCustom(&fpdf.InitType{
		OrientationStr: "P",
		UnitStr:        "mm",
		SizeStr:        "",
		Size: fpdf.SizeType{
			Wd: 80,  // 80mm width
			Ht: 300, // increased height to accommodate options
		},
	})
	pdf.AddPage()

	// Add Thai font
	pdf.AddUTF8Font("NotoSansThai", "", `E:\h_lab\go\poc_pos_restaurant\font\NotoSansThai-Regular.ttf`)
	pdf.AddUTF8Font("NotoSansThai", "B", `E:\h_lab\go\poc_pos_restaurant\font\NotoSansThai-Bold.ttf`)

	pdf.SetLeftMargin(5)
	pdf.SetRightMargin(5)

	// Header
	pdf.SetFont("NotoSansThai", "B", 12)
	pdf.CellFormat(0, 6, "ใบเสร็จรับเงิน", "", 1, "C", false, 0, "")
	pdf.SetFont("NotoSansThai", "", 9)
	pdf.CellFormat(0, 5, "ร้านอาหารดีเลิศ", "", 1, "C", false, 0, "")
	pdf.CellFormat(0, 5, "123 ถนนสุขุมวิท กรุงเทพฯ 10110", "", 1, "C", false, 0, "")
	pdf.CellFormat(0, 5, "โทร: 02-123-4567", "", 1, "C", false, 0, "")
	pdf.Ln(2)

	// Receipt info
	pdf.SetFont("NotoSansThai", "", 8)
	pdf.CellFormat(0, 5, fmt.Sprintf("เลขที่: %d", order.ID), "", 1, "L", false, 0, "")
	pdf.CellFormat(0, 5, fmt.Sprintf("วันที่: %s", order.CreatedAt.Format("02/01/2006 15:04")), "", 1, "L", false, 0, "")
	pdf.CellFormat(0, 5, fmt.Sprintf("โต๊ะ: %d", order.TableID), "", 1, "L", false, 0, "")
	pdf.Ln(2)
	pdf.Line(0, pdf.GetY(), 80, pdf.GetY())
	pdf.Ln(2)

	// Items with options
	ctx := context.Background()
	var totalWithOptions float64

	for _, item := range order.Items {
		// Main item
		itemPrice := item.UnitPrice.AmountBaht()
		itemSubtotal := float64(item.Quantity) * itemPrice

		pdf.SetFont("NotoSansThai", "", 8)
		pdf.CellFormat(0, 4, item.Name, "", 1, "L", false, 0, "")
		pdf.CellFormat(0, 4, fmt.Sprintf("  %d x %.2f บาท = %.2f บาท",
			item.Quantity, itemPrice, itemSubtotal), "", 1, "L", false, 0, "")

		// Get and display options
		options, err := s.getItemOptions(ctx, item.ID)
		if err == nil && len(options) > 0 {
			pdf.SetFont("NotoSansThai", "", 7)
			for _, opt := range options {
				optionPrice := opt.AdditionalPrice.AmountBaht()
				optionTotal := float64(item.Quantity) * optionPrice

				if optionPrice > 0 {
					pdf.CellFormat(0, 3, fmt.Sprintf("    + %s: %s (+%.2f บาท x%d = +%.2f บาท)",
						opt.Option.Name, opt.Value.Name, optionPrice, item.Quantity, optionTotal), "", 1, "L", false, 0, "")
					itemSubtotal += optionTotal
				} else {
					pdf.CellFormat(0, 3, fmt.Sprintf("    + %s: %s",
						opt.Option.Name, opt.Value.Name), "", 1, "L", false, 0, "")
				}
			}
		}

		totalWithOptions += itemSubtotal
		pdf.Ln(1)
	}

	pdf.Ln(2)
	pdf.Line(0, pdf.GetY(), 80, pdf.GetY())
	pdf.Ln(2)

	// Summary with corrected calculations
	subtotal, _ := vo.NewMoneyFromBaht(totalWithOptions)
	discount := subtotal.Multiply(DISCOUNT_RATE)
	afterDiscount, _ := subtotal.Subtract(discount)
	tax := afterDiscount.Multiply(TAX_RATE)
	finalTotal := afterDiscount.Add(tax)

	pdf.SetFont("NotoSansThai", "", 8)
	pdf.CellFormat(0, 5, fmt.Sprintf("ยอดรวม: %.2f บาท", subtotal.AmountBaht()), "", 1, "R", false, 0, "")

	if discount.AmountBaht() > 0 {
		pdf.CellFormat(0, 5, fmt.Sprintf("ส่วนลด %.0f%%: -%.2f บาท", DISCOUNT_RATE*100, discount.AmountBaht()), "", 1, "R", false, 0, "")
	}

	if tax.AmountBaht() > 0 {
		pdf.CellFormat(0, 5, fmt.Sprintf("VAT %.0f%%: %.2f บาท", TAX_RATE*100, tax.AmountBaht()), "", 1, "R", false, 0, "")
	}

	pdf.SetFont("NotoSansThai", "B", 10)
	pdf.CellFormat(0, 6, fmt.Sprintf("ยอดสุทธิ: %.2f บาท", finalTotal.AmountBaht()), "", 1, "R", false, 0, "")
	pdf.Ln(4)

	// Footer
	pdf.SetFont("NotoSansThai", "", 8)
	pdf.CellFormat(0, 5, "ขอบคุณที่ใช้บริการ", "", 1, "C", false, 0, "")
	pdf.CellFormat(0, 5, "Thank you for your business!", "", 1, "C", false, 0, "")

	return pdf.Output(writer)
}
func generateOrderQRCodePDF(receipt *entity.Order, writer io.Writer) error {
	pdf := fpdf.NewCustom(&fpdf.InitType{
		OrientationStr: "P",
		UnitStr:        "mm",
		SizeStr:        "",
		Size: fpdf.SizeType{
			Wd: 80,  // 80mm width
			Ht: 250, // adjustable height
		},
	})
	pdf.AddPage()

	// Add Thai font
	pdf.AddUTF8Font("NotoSansThai", "", `E:\h_lab\go\poc_pos_restaurant\font\NotoSansThai-Regular.ttf`)
	pdf.AddUTF8Font("NotoSansThai", "B", `E:\h_lab\go\poc_pos_restaurant\font\NotoSansThai-Bold.ttf`)

	pdf.SetLeftMargin(5)
	pdf.SetRightMargin(5)

	// Header
	pdf.SetFont("NotoSansThai", "B", 12)
	pdf.CellFormat(0, 6, "Order QR Code", "", 1, "C", false, 0, "")
	pdf.SetFont("NotoSansThai", "", 9)
	pdf.CellFormat(0, 5, "โต๊ะ: "+fmt.Sprintf("%d", receipt.TableID), "", 1, "C", false, 0, "")
	pdf.CellFormat(0, 5, "วันที่: "+receipt.CreatedAt.Format("02/01/2006 15:04"), "", 1, "C", false, 0, "")
	pdf.Ln(2)

	// QR Code
	if receipt.QRCode != "" {
		qrFile := "temp_qr.png"
		err := qrcode.WriteFile(receipt.QRCode, qrcode.Medium, 256, qrFile)
		if err != nil {
			return fmt.Errorf("failed to create QR code: %v", err)
		}
		defer os.Remove(qrFile)

		qrX := (80.0 - 40.0) / 2
		pdf.Image(qrFile, qrX, pdf.GetY(), 40, 40, false, "", 0, "")
		pdf.Ln(45)
		pdf.SetFont("NotoSansThai", "B", 8)
		pdf.CellFormat(0, 5, "สแกน QR เพื่อสั่งอาหารหรือดูโปรโมชั่น", "", 1, "C", false, 0, "")
		pdf.Ln(2)
	}

	// Footer
	pdf.SetFont("NotoSansThai", "", 8)
	pdf.CellFormat(0, 5, "สแกน QR เพื่อสั่งอาหารหรือดูโปรโมชั่น", "", 1, "C", false, 0, "")
	pdf.CellFormat(0, 5, "ขอบคุณที่ใช้บริการ", "", 1, "C", false, 0, "")

	return pdf.Output(writer)
}
