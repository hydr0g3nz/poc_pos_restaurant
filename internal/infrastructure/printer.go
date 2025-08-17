package infrastructure

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/coder/websocket"
)

type printerServer struct {
	con *websocket.Conn
}

func NewPrinterService(url string) (*printerServer, error) {
	p := &printerServer{}
	if err := p.connect(url); err != nil {
		return nil, err
	}
	log.Println("Printer service connected to", url)
	return p, nil
}
func (s *printerServer) Print(ctx context.Context, content []byte, contentType string) error {
	if s.con == nil {
		return fmt.Errorf("printer connection is not established")
	}
	return s.con.Write(ctx, websocket.MessageText, content)
}
func (s *printerServer) Close() error {
	if s.con != nil {
		return s.con.Close(websocket.StatusNormalClosure, "closing printer connection")
	}
	log.Println("Printer connection is already closed")
	return nil
}
func (s *printerServer) connect(url string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, _, err := websocket.Dial(ctx, url, nil)
	if err != nil {
		log.Fatal("dial error:", err)
	}
	s.con = conn
	return nil
}
