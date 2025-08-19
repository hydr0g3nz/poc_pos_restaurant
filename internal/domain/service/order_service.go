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
	orderRepo     repository.OrderRepository
	orderItemRepo repository.OrderItemRepository
	tableRepo     repository.TableRepository
	menuItemRepo  repository.MenuItemRepository
}

func NewOrderService(
	orderRepo repository.OrderRepository,
	orderItemRepo repository.OrderItemRepository,
	tableRepo repository.TableRepository,
	menuItemRepo repository.MenuItemRepository,
) OrderService {
	return &orderService{
		orderRepo:     orderRepo,
		orderItemRepo: orderItemRepo,
		tableRepo:     tableRepo,
		menuItemRepo:  menuItemRepo,
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

	return order.CalculateTotal(), nil
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
	if err := generateReceiptPDF(order, w); err != nil {
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
func generateReceiptPDF(receipt *entity.Order, writer io.Writer) error {
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
	pdf.CellFormat(0, 6, "ใบเสร็จรับเงิน", "", 1, "C", false, 0, "")
	pdf.SetFont("NotoSansThai", "", 9)
	pdf.CellFormat(0, 5, "ร้านอาหารดีเลิศ", "", 1, "C", false, 0, "")
	pdf.CellFormat(0, 5, "123 ถนนสุขุมวิท กรุงเทพฯ 10110", "", 1, "C", false, 0, "")
	pdf.CellFormat(0, 5, "โทร: 02-123-4567", "", 1, "C", false, 0, "")
	pdf.Ln(2)

	// Receipt info
	pdf.SetFont("NotoSansThai", "", 8)
	pdf.CellFormat(0, 5, fmt.Sprintf("เลขที่: %d", receipt.ID), "", 1, "L", false, 0, "")
	pdf.CellFormat(0, 5, fmt.Sprintf("วันที่: %s", receipt.CreatedAt.Format("02/01/2006 15:04")), "", 1, "L", false, 0, "")
	pdf.CellFormat(0, 5, fmt.Sprintf("ลูกค้า: %s", "test"), "", 1, "L", false, 0, "")
	pdf.Ln(2)
	pdf.Line(0, pdf.GetY(), 80, pdf.GetY())
	pdf.Ln(2)

	// Items (simple list)
	for _, item := range receipt.Items {
		line := fmt.Sprintf("%s x%d   %.2f บาท", item.Name, item.Quantity, float64(item.Quantity)*item.UnitPrice.AmountBaht())
		pdf.CellFormat(0, 5, line, "", 1, "L", false, 0, "")
	}
	pdf.Ln(2)
	pdf.Line(0, pdf.GetY(), 80, pdf.GetY())
	pdf.Ln(2)

	// Summary
	subtotal := receipt.CalculateTotal()
	discount := receipt.CalculateDiscount()
	tax := receipt.CalculateTax()
	var total vo.Money
	var err error
	if !discount.IsZero() {
		if total, err = subtotal.Subtract(discount); err != nil {
			return fmt.Errorf("failed to calculate total after discount: %w", err)
		}
	}
	if !tax.IsZero() {
		total = total.Add(tax)
	}

	pdf.CellFormat(0, 5, fmt.Sprintf("ยอดรวม: %.2f บาท", subtotal.AmountBaht()), "", 1, "R", false, 0, "")
	if discount.AmountBaht() > 0 {
		pdf.CellFormat(0, 5, fmt.Sprintf("ส่วนลด %.1f%%: -%.2f บาท", DISCOUNT_RATE, discount.AmountBaht()), "", 1, "R", false, 0, "")
	}
	if tax.AmountBaht() > 0 {
		pdf.CellFormat(0, 5, fmt.Sprintf("VAT %.1f%%: %.2f บาท", TAX_RATE, tax.AmountBaht()), "", 1, "R", false, 0, "")
	}
	pdf.SetFont("NotoSansThai", "B", 10)
	pdf.CellFormat(0, 6, fmt.Sprintf("ยอดสุทธิ: %.2f บาท", total.AmountBaht()), "", 1, "R", false, 0, "")
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
