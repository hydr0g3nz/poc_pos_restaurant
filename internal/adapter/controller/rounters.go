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
	orderGroup.Get("/items", c.ListOrdersWithItems)
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
	tableGroup.Get("/qr", c.GetTableByQRCode)

}

// เพิ่มใน internal/adapter/controller/routers.go

// RegisterRoutes registers the routes for the menu option controller
func (c *MenuOptionController) RegisterRoutes(router fiber.Router) {
	menuOptionGroup := router.Group("/menu-options")

	// Public routes
	menuOptionGroup.Get("/", c.ListMenuOptions)
	menuOptionGroup.Get("/search", c.GetMenuOptionsByType) // GET /menu-options/search?type=spice_level
	menuOptionGroup.Get("/:id", c.GetMenuOption)

	// Admin routes
	menuOptionGroup.Post("/", c.CreateMenuOption)
	menuOptionGroup.Put("/:id", c.UpdateMenuOption)
	menuOptionGroup.Delete("/:id", c.DeleteMenuOption)
}

// RegisterRoutes registers the routes for the option value controller
func (c *OptionValueController) RegisterRoutes(router fiber.Router) {
	optionValueGroup := router.Group("/option-values")

	// Public routes
	optionValueGroup.Get("/:id", c.GetOptionValue)
	optionValueGroup.Get("/option/:optionId", c.GetOptionValuesByOptionID) // GET /option-values/option/1

	// Admin routes
	optionValueGroup.Post("/", c.CreateOptionValue)
	optionValueGroup.Put("/:id", c.UpdateOptionValue)
	optionValueGroup.Delete("/:id", c.DeleteOptionValue)
}

// RegisterRoutes registers the routes for the menu item option controller
func (c *MenuItemOptionController) RegisterRoutes(router fiber.Router) {
	menuItemOptionGroup := router.Group("/menu-item-options")

	// Routes
	menuItemOptionGroup.Get("/item/:itemId", c.GetMenuItemOptions)                           // GET /menu-item-options/item/1
	menuItemOptionGroup.Post("/", c.AddOptionToMenuItem)                                     // POST /menu-item-options
	menuItemOptionGroup.Delete("/item/:itemId/option/:optionId", c.RemoveOptionFromMenuItem) // DELETE /menu-item-options/item/1/option/2
}

// RegisterRoutes registers the routes for the kitchen controller
func (c *KitchenController) RegisterRoutes(router fiber.Router) {
	kitchenGroup := router.Group("/kitchen")

	// Kitchen routes
	kitchenGroup.Get("/", c.GetKitchenStatation)
	kitchenGroup.Post("/", c.CreateKitchenStatation)
	kitchenGroup.Put("/", c.UpdateKitchenStatation)
	kitchenGroup.Delete("/", c.DeleteKitchenStatation)

	kitchenGroup.Get("/queue", c.GetKitchenQueue)                           // GET /kitchen/queue
	kitchenGroup.Get("/items", c.GetOrderItemsByStatus)                     // GET /kitchen/items?status=preparing
	kitchenGroup.Get("/station/orders", c.GetKitchenOrdersByStation)        // GET /kitchen/station?station=grill
	kitchenGroup.Put("/items/:orderItemId/status", c.UpdateOrderItemStatus) // PUT /kitchen/items/1/status
	kitchenGroup.Put("/items/:orderItemId/ready", c.MarkOrderItemAsReady)   // PUT /kitchen/items/1/ready
	kitchenGroup.Put("/items/:orderItemId/served", c.MarkOrderItemAsServed) // PUT /kitchen/items/1/served
}
func (c *CustomerController) RegisterRoutes(router fiber.Router) {
	customerGroup := router.Group("/customers")
	// Customer routes
	customerGroup.Get("/menu", c.ListMenuItems)
	customerGroup.Get("/menu/items/:id", c.GetMenuItem)
	customerGroup.Get("/category", c.ListCategory)
	// customerGroup.Post("/", c.CreateCustomer)
	// customerGroup.Put("/:id", c.UpdateCustomer)
	// customerGroup.Delete("/:id", c.DeleteCustomer)
}
