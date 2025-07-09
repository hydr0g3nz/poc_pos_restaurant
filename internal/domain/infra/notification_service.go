package infra

import "context"

// NotificationService interface for notifications (extended)
type NotificationService interface {
	// SendPushNotification sends push notification to user
	SendPushNotification(ctx context.Context, userID int, title, body string, data map[string]interface{}) error

	// SendBulkNotification sends notification to multiple users
	SendBulkNotification(ctx context.Context, userIDs []int, title, body string, data map[string]interface{}) error

	// RegisterDevice registers device for push notifications
	RegisterDevice(ctx context.Context, userID int, deviceToken, platform string) error

	// UnregisterDevice removes device from notifications
	UnregisterDevice(ctx context.Context, deviceToken string) error

	// Kitchen notifications
	SendKitchenNotification(ctx context.Context, orderID int, message string) error

	// Table service notifications
	SendTableNotification(ctx context.Context, tableID int, message string) error

	// Order status notifications
	SendOrderStatusNotification(ctx context.Context, orderID int, status string) error
}
