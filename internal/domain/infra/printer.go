package infra

import "context"

type PrinterService interface {
	Print(ctx context.Context, content []byte, contentType string) error
	Close() error
}
