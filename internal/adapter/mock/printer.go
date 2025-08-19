package mocks

import (
	"context"
	"fmt"
)

type PrinterService struct{}

// Print(ctx context.Context, content []byte, contentType string) error
// Close() error
func NewPrinterService() *PrinterService {
	return &PrinterService{}
}
func (p *PrinterService) Print(ctx context.Context, content []byte, contentType string) error {
	fmt.Println("Mock Print called with content:", string(content), "and contentType:", contentType)
	return nil
}
func (p *PrinterService) Close() error {
	return nil
}
