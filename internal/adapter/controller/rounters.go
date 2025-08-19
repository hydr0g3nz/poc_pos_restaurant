package controller

import "github.com/gofiber/fiber/v2"

// RegisterRoutes registers the routes for the category controller
func (c *CategoryController) RegisterRoutes(router fiber.Router) {
	categoryGroup := router.Group("/categories")

	// Public routes
	categoryGroup.Get("/", c.ListCategories)
	categoryGroup.Get("/search", c.GetCategoryByName) // GET /categories/search?name=ของคาว
	categoryGroup.Get("/:id", c.GetCategory)

	// Admin routes (require admin role in real implementation)
	categoryGroup.Post("/", c.CreateCategory)
	categoryGroup.Put("/:id", c.UpdateCategory)
	categoryGroup.Delete("/:id", c.DeleteCategory)
}

// RegisterRoutes registers the routes for the menu item controller
func (c *MenuItemController) RegisterRoutes(router fiber.Router) {
	menuItemGroup := router.Group("/menu-items")

	// Public routes
	menuItemGroup.Get("/", c.ListMenuItems)
	menuItemGroup.Get("/search", c.SearchMenuItems)                       // GET /menu-items/search?q=ข้าวผัด
	menuItemGroup.Get("/category/:categoryId", c.ListMenuItemsByCategory) // GET /menu-items/category/1
	menuItemGroup.Get("/:id", c.GetMenuItem)

	// Admin routes (require admin role in real implementation)
	menuItemGroup.Post("/", c.CreateMenuItem)
	menuItemGroup.Put("/:id", c.UpdateMenuItem)
	menuItemGroup.Delete("/:id", c.DeleteMenuItem)
}
func (c *OrderController) RegisterRoutes(router fiber.Router) {
	orderGroup := router.Group("/orders")

	// Order routes
	orderGroup.Post("/", c.CreateOrder)
	orderGroup.Get("/", c.ListOrders)
	orderGroup.Get("/search", c.GetOrdersByStatus)        // GET /orders/search?status=open
	orderGroup.Get("/date-range", c.GetOrdersByDateRange) // GET /orders/date-range?start_date=2024-01-01&end_date=2024-01-31
	orderGroup.Get("/:id", c.GetOrder)
	orderGroup.Get("/:id/items", c.GetOrderWithItems)
	orderGroup.Put("/:id", c.UpdateOrder)
	orderGroup.Put("/:id/close", c.CloseOrder)
	orderGroup.Get("/:id/print/receipt", c.PrintOrderReceipt) // Print order by ID
	orderGroup.Get("/:id/print/qrcode", c.PrintOrderQRCode)   // Print order by ID
	// Order by table routes
	orderGroup.Get("/table/:tableId", c.ListOrdersByTable)
	orderGroup.Get("/table/:tableId/open", c.GetOpenOrderByTable)

	// Order items routes
	orderGroup.Post("/items", c.AddOrderItem)
	orderGroup.Put("/items/:id", c.UpdateOrderItem)
	orderGroup.Delete("/items/:id", c.RemoveOrderItem)
	orderGroup.Get("/:orderId/items", c.ListOrderItems)
	orderGroup.Get("/:orderId/total", c.CalculateOrderTotal)
}

// RegisterRoutes registers the routes for the payment controller
func (c *PaymentController) RegisterRoutes(router fiber.Router) {
	paymentGroup := router.Group("/payments")

	// Payment routes
	paymentGroup.Post("/", c.ProcessPayment)
	paymentGroup.Get("/", c.ListPayments)
	paymentGroup.Get("/search", c.ListPaymentsByMethod)        // GET /payments/search?method=cash
	paymentGroup.Get("/date-range", c.ListPaymentsByDateRange) // GET /payments/date-range?start_date=2024-01-01&end_date=2024-01-31
	paymentGroup.Get("/:id", c.GetPayment)
	paymentGroup.Get("/order/:orderId", c.GetPaymentByOrder)
}

// RegisterRoutes registers the routes for the revenue controller
func (c *RevenueController) RegisterRoutes(router fiber.Router) {
	revenueGroup := router.Group("/revenue")

	// Daily revenue routes
	revenueGroup.Get("/daily", c.GetDailyRevenue)            // GET /revenue/daily?date=2024-01-01
	revenueGroup.Get("/daily/range", c.GetDailyRevenueRange) // GET /revenue/daily/range?start_date=2024-01-01&end_date=2024-01-31

	// Monthly revenue routes
	revenueGroup.Get("/monthly", c.GetMonthlyRevenue)            // GET /revenue/monthly?year=2024&month=1
	revenueGroup.Get("/monthly/range", c.GetMonthlyRevenueRange) // GET /revenue/monthly/range?start_date=2024-01-01&end_date=2024-12-31

	// Total revenue route
	revenueGroup.Get("/total", c.GetTotalRevenue) // GET /revenue/total?start_date=2024-01-01&end_date=2024-12-31
}

// RegisterRoutes registers the routes for the table controller
func (c *TableController) RegisterRoutes(router fiber.Router) {
	tableGroup := router.Group("/tables")

	// Table CRUD operations
	tableGroup.Post("/", c.CreateTable)
	tableGroup.Get("/", c.ListTables)
	tableGroup.Get("/:id", c.GetTable)
	tableGroup.Put("/:id", c.UpdateTable)
	tableGroup.Delete("/:id", c.DeleteTable)

	// Table by number
	tableGroup.Get("/number/:number", c.GetTableByNumber)

	// QR code operations
	tableGroup.Get("/qr", c.GetTableByQRCode)               // GET /tables/qr?qr_code=/order?table=1
	tableGroup.Post("/:id/qr-code", c.GenerateQRCode)       // POST /tables/1/qr-code
	tableGroup.Get("/scan", c.ScanQRCode)                   // GET /tables/scan?qr_code=/order?table=1
	tableGroup.Post("/scan/order", c.CreateOrderFromQRCode) // POST /tables/scan/order?qr_code=/order?table=1
}
