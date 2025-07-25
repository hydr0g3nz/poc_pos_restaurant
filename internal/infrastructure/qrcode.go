package infrastructure

import (
	"bytes"
	"context"
	"log"

	"github.com/yeqown/go-qrcode/v2"
	"github.com/yeqown/go-qrcode/writer/standard"
)

type QRCodeService struct {
}

func NewQRCodeService() *QRCodeService {
	return &QRCodeService{}
}

type bufferWriteCloser struct {
	*bytes.Buffer
}

func (bwc bufferWriteCloser) Close() error {
	// ไม่มีอะไรต้องปิดสำหรับ buffer แต่จำเป็นต้องมีเพื่อ implement io.WriteCloser
	return nil
}

func (s *QRCodeService) GenerateQRCodeImage(ctx context.Context, data string) ([]byte, error) {
	qr, err := qrcode.New(data)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	bwc := bufferWriteCloser{&buf}
	writer := standard.NewWithWriter(&bwc)

	// เขียน QR Code ลง buffer
	if err := qr.Save(writer); err != nil {
		log.Printf("Error saving QR code: %v", err)
		return nil, err
	}

	return buf.Bytes(), nil
}
