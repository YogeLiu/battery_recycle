package v1

import (
	"battery-erp-backend/internal/services"

	"github.com/gin-gonic/gin"
)

// SetupRoutes configures all API routes
func SetupRoutes(engine *gin.Engine, services *services.Services) {
	// API v1 group
	v1 := engine.Group("/jxc/v1")

	// Controllers
	authController := NewAuthController(services.Auth)
	userController := NewUserController(services.UserService)
	categoryController := NewCategoryController(services.CategoryService)
	inboundController := NewInboundController(services.InboundService)
	outboundController := NewOutboundController(services.OutboundService)
	inventoryController := NewInventoryController(services.InventoryService)
	reportController := NewReportController(services.ReportService)

	// Auth routes (no middleware)
	authRoutes := v1.Group("/auth")
	{
		authRoutes.POST("/login", authController.Login)
	}

	// Protected routes (with auth middleware)
	authMiddleware := NewAuthMiddleware(services.Auth)

	// User routes
	userRoutes := v1.Group("/users")
	userRoutes.Use(authMiddleware.RequireAuth(), authMiddleware.RequireRole("super_admin"))
	{
		userRoutes.GET("", userController.GetAll)
		userRoutes.POST("", userController.Create)
		userRoutes.GET("/:id", userController.GetByID)
		userRoutes.PUT("/:id", userController.Update)
		userRoutes.DELETE("/:id", userController.Delete)
	}

	// Category routes
	categoryRoutes := v1.Group("/categories")
	categoryRoutes.Use(authMiddleware.RequireAuth())
	{
		categoryRoutes.GET("", categoryController.GetAll)
		categoryRoutes.POST("", authMiddleware.RequireRole("super_admin"), categoryController.Create)
		categoryRoutes.GET("/:id", categoryController.GetByID)
		categoryRoutes.PUT("/:id", authMiddleware.RequireRole("super_admin"), categoryController.Update)
		categoryRoutes.DELETE("/:id", authMiddleware.RequireRole("super_admin"), categoryController.Delete)
	}

	// Inbound routes
	inboundRoutes := v1.Group("/inbound/orders")
	inboundRoutes.Use(authMiddleware.RequireAuth())
	{
		inboundRoutes.POST("/search", inboundController.GetAll)
		inboundRoutes.POST("", inboundController.Create)
		inboundRoutes.GET("/:id", inboundController.GetByID)
		inboundRoutes.PUT("/:id", inboundController.Update)
		inboundRoutes.DELETE("/:id", inboundController.Delete)
	}

	// Outbound routes
	outboundRoutes := v1.Group("/outbound/orders")
	outboundRoutes.Use(authMiddleware.RequireAuth())
	{
		outboundRoutes.GET("", outboundController.GetAll)
		outboundRoutes.POST("", outboundController.Create)
		outboundRoutes.GET("/:id", outboundController.GetByID)
		outboundRoutes.PUT("/:id", outboundController.Update)
		outboundRoutes.DELETE("/:id", outboundController.Delete)
	}

	// Inventory routes
	inventoryRoutes := v1.Group("/inventory")
	inventoryRoutes.Use(authMiddleware.RequireAuth())
	{
		inventoryRoutes.GET("", inventoryController.GetAll)
		inventoryRoutes.GET("/:categoryId", inventoryController.GetByCategoryID)
	}

	// Report routes
	reportRoutes := v1.Group("/reports")
	reportRoutes.Use(authMiddleware.RequireAuth())
	{
		reportRoutes.GET("/summary", reportController.GetSummary)
	}
}
