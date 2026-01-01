package http

import (
	"github.com/gin-gonic/gin"
	"github.com/shoelfikar/voucher-management-system/internal/delivery/http/handler"
)

// SetupRouter configures and returns the Gin router with all routes
func SetupRouter(
	authHandler *handler.AuthHandler,
	voucherHandler *handler.VoucherHandler,
	authMiddleware gin.HandlerFunc,
	corsMiddleware gin.HandlerFunc,
) *gin.Engine {
	r := gin.Default()

	r.Use(corsMiddleware)

	// Health check endpoint (public)
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "Voucher Management System API is running",
		})
	})

	api := r.Group("/api/v1")
	{
		// Auth routes (public)
		api.POST("/auth/login", authHandler.Login)

		protected := api.Group("")
		protected.Use(authMiddleware)
		{
			// Voucher routes
			vouchers := protected.Group("/vouchers")
			{
				vouchers.GET("", voucherHandler.GetAll)
				vouchers.GET("/:id", voucherHandler.GetByID)
				vouchers.POST("", voucherHandler.Create)
				vouchers.PUT("/:id", voucherHandler.Update)
				vouchers.DELETE("/:id", voucherHandler.Delete)

				vouchers.POST("/upload-csv", voucherHandler.ImportCSV)
				vouchers.POST("/upload-batch", voucherHandler.UploadBatch)
				vouchers.GET("/export", voucherHandler.ExportCSV)
			}
		}
	}

	return r
}
