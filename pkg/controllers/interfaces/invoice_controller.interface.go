package controller_interfaces

import "github.com/gin-gonic/gin"

type InvoiceController interface {
	Create(ctx *gin.Context)
	GetStatistics(ctx *gin.Context)
	GetCustomerInvoices(ctx *gin.Context)
	Duplicate(ctx *gin.Context)
	GetShareableLink(ctx *gin.Context)
	GetCustomerAuditTrails(ctx *gin.Context)
	GetSingleInvoiceAuditTrails(ctx *gin.Context)
	SetReminder(ctx *gin.Context)
	GetDetails(ctx *gin.Context)
	ConfirmPayment(ctx *gin.Context)
}
