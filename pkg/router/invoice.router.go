package router

import (
	controller_interfaces "github.com/Adebayobenjamin/numerisbook/pkg/controllers/interfaces"
	"github.com/Adebayobenjamin/numerisbook/pkg/middleware"
	"github.com/gin-gonic/gin"
)

func NewInvoiceRouter(invoiceController controller_interfaces.InvoiceController, router *gin.RouterGroup) *gin.RouterGroup {
	invoiceRouter := router.Group("/invoices")
	invoiceRouter.Use(middlewares.RequiresAuthHeader())

	// Create and manage invoices
	invoiceRouter.POST("", invoiceController.Create)
	invoiceRouter.GET("/:invoice_id", invoiceController.GetDetails)
	invoiceRouter.GET("/statistics", invoiceController.GetStatistics)
	invoiceRouter.GET("", invoiceController.GetCustomerInvoices)

	// Payment confirmation
	invoiceRouter.POST("/:invoice_id/confirm-payment", invoiceController.ConfirmPayment)

	// Reminders
	invoiceRouter.POST("/:invoice_id/reminders", invoiceController.SetReminder)

	// Shareable link
	invoiceRouter.GET("/:invoice_id/shareable-link", invoiceController.GetShareableLink)

	// Invoice duplication
	invoiceRouter.POST("/:invoice_id/duplicate", invoiceController.Duplicate)

	// Audit trails
	invoiceRouter.GET("/audit-trails", invoiceController.GetCustomerAuditTrails)
	invoiceRouter.GET("/:invoice_id/audit-trails", invoiceController.GetSingleInvoiceAuditTrails)

	return invoiceRouter
}
