package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Adebayobenjamin/numerisbook/pkg/common"
	"github.com/Adebayobenjamin/numerisbook/pkg/common/exceptions"
	controller_interfaces "github.com/Adebayobenjamin/numerisbook/pkg/controllers/interfaces"
	request_dto "github.com/Adebayobenjamin/numerisbook/pkg/dtos/request"
	"github.com/Adebayobenjamin/numerisbook/pkg/helper"
	"github.com/Adebayobenjamin/numerisbook/pkg/models"
	services_interfaces "github.com/Adebayobenjamin/numerisbook/pkg/services/interfaces"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type invoiceController struct {
	logger          *zerolog.Logger
	invoiceService  services_interfaces.InvoiceService
	auditService    services_interfaces.AuditService
	reminderService services_interfaces.RemiderService
	customerService services_interfaces.CustomerService
}

// ConfirmPayment implements controller_interfaces.InvoiceController.
func (i *invoiceController) ConfirmPayment(ctx *gin.Context) {
	var request request_dto.PaymentConfirmationRequest
	var err error
	var invoice *models.Invoice

	if err = ctx.ShouldBindJSON(&request); err != nil {
		exceptions.ThrowUnProcessableEntityException(ctx, err.Error())
		return
	}

	// get customer from context (set by middleware)
	customer, err := i.getCustomerFromContext(ctx)
	if err != nil {
		exceptions.ThrowBadRequestException(ctx, err.Error())
		return
	}

	// get invoice from params and fetch from db
	invoice, err = i.getInvoiceDetailsFromParams(ctx, customer.ID)
	if err != nil {
		exceptions.ThrowBadRequestException(ctx, err.Error())
		return
	}

	defer func() {
		// audit payment confirmation
		if err == nil && invoice != nil {
			i.auditService.CreateAuditTrail(
				ctx,
				models.EventTypePaymentConfirmed,
				models.LogLevelInfo,
				fmt.Sprintf("Confirmed Payment for Invoice %s/%s", invoice.InvoiceNumber, customer.Name),
				invoice.ID,
				customer.ID,
			)
		}
	}()

	// TODO: call service to validate payment amount
	err = i.invoiceService.ValidatePaymentAmount(ctx, request.Amount, invoice, request.IsPartial)
	if err != nil {
		exceptions.ThrowBadRequestException(ctx, err.Error())
		return
	}

	// TODO: call service to confirm payment
	err = i.invoiceService.ConfirmPayment(ctx, invoice.ID, request.Amount, request.PaymentDate, request.IsPartial)
	if err != nil {
		exceptions.ThrowBadRequestException(ctx, err.Error())
		return
	}

	err = i.invoiceService.SetInvoiceStatusIfFullyPaid(ctx, invoice)
	if err != nil {
		exceptions.ThrowBadRequestException(ctx, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, common.BuildSuccessResponse("payment confirmed successfully", nil))
}

// Create implements controller_interfaces.InvoiceController.
func (i *invoiceController) Create(ctx *gin.Context) {
	var request request_dto.CreateInvoiceRequest
	var err error
	var invoice *models.Invoice

	if err = ctx.ShouldBindJSON(&request); err != nil {
		exceptions.ThrowUnProcessableEntityException(ctx, err.Error())
		return
	}

	customer, err := i.getCustomerFromContext(ctx)
	if err != nil {
		exceptions.ThrowBadRequestException(ctx, err.Error())
		return
	}

	customer, err = i.customerService.GetCustomerByID(ctx, customer.ID)
	if err != nil {
		exceptions.ThrowBadRequestException(ctx, err.Error())
		return
	}

	defer func() {
		// audit invoice creation
		if err == nil && invoice != nil {
			i.auditService.CreateAuditTrail(
				ctx,
				models.EventTypeInvoiceCreated,
				models.LogLevelInfo,
				fmt.Sprintf("Created Invoice %s/%s", invoice.InvoiceNumber, customer.Name),
				invoice.ID,
				customer.ID,
			)
		}
	}()

	// call service to create invoice
	invoice, err = i.invoiceService.CreateInvoice(ctx, customer.ID, &request)
	if err != nil {
		exceptions.ThrowBadRequestException(ctx, err.Error())
		return
	}

	// if reminder schedules are provided, set reminders
	if len(request.ReminderSchedules) > 0 {
		err = i.reminderService.SetInvoiceReminders(ctx, invoice, customer.ID, request.ReminderSchedules)
		if err != nil {
			exceptions.ThrowBadRequestException(ctx, err.Error())
			return
		}
	}

	ctx.JSON(http.StatusCreated, common.BuildSuccessResponse("invoice created successfully", invoice))
}

// Duplicate implements controller_interfaces.InvoiceController.
func (i *invoiceController) Duplicate(ctx *gin.Context) {
	invoiceID := ctx.Param("invoice_id")
	var newInvoice *models.Invoice

	if invoiceID == "" {
		exceptions.ThrowUnProcessableEntityException(ctx, "invoice id is required")
		return
	}

	invoiceIDUint, err := strconv.ParseUint(invoiceID, 10, 64)
	if err != nil {
		exceptions.ThrowUnProcessableEntityException(ctx, "invalid invoice id")
		return
	}

	customer, err := i.getCustomerFromContext(ctx)
	if err != nil {
		exceptions.ThrowBadRequestException(ctx, err.Error())
		return
	}

	defer func() {
		i.auditService.CreateAuditTrail(
			ctx,
			models.EventTypeInvoiceCreated,
			models.LogLevelInfo,
			fmt.Sprintf("Created Invoice %s/%s", newInvoice.InvoiceNumber, customer.Name),
			newInvoice.ID,
			customer.ID,
		)

	}()

	// TODO: call service to copy invoice details
	existingInvoice, err := i.invoiceService.GetInvoiceByIDandCustomer(ctx, uint(invoiceIDUint), customer.ID)
	if err != nil {
		exceptions.ThrowBadRequestException(ctx, err.Error())
		return
	}

	// TODO: call service to create the new invoice
	newInvoice, err = i.invoiceService.DuplicateInvoice(ctx, existingInvoice)
	if err != nil {
		exceptions.ThrowBadRequestException(ctx, err.Error())
		return
	}

	ctx.JSON(http.StatusCreated, common.BuildSuccessResponse("invoice duplicated successfully", newInvoice))
}

// GetCustomerInvoices implements controller_interfaces.InvoiceController.
func (i *invoiceController) GetCustomerInvoices(ctx *gin.Context) {
	var request request_dto.GetAllRequest

	if err := ctx.ShouldBindQuery(&request); err != nil {
		exceptions.ThrowUnProcessableEntityException(ctx, err.Error())
		return
	}

	customer, err := i.getCustomerFromContext(ctx)
	if err != nil {
		exceptions.ThrowBadRequestException(ctx, err.Error())
		return
	}

	invoices, err := i.invoiceService.GetCustomerInvoices(ctx, request.Limit, request.Page, customer.ID)
	if err != nil {
		exceptions.ThrowBadRequestException(ctx, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, common.BuildSuccessResponse("invoices fetched successfully", invoices))
}

// GetCustomerAuditTrails implements controller_interfaces.InvoiceController.
func (i *invoiceController) GetCustomerAuditTrails(ctx *gin.Context) {
	var request request_dto.GetAllRequest

	if err := ctx.ShouldBindQuery(&request); err != nil {
		exceptions.ThrowUnProcessableEntityException(ctx, err.Error())
		return
	}

	customer, err := i.getCustomerFromContext(ctx)
	if err != nil {
		exceptions.ThrowBadRequestException(ctx, err.Error())
		return
	}

	//TODO: call service to get all audit trails
	auditTrails, err := i.auditService.GetCustomerAuditTrails(ctx, uint(customer.ID), request.Limit, request.Page)
	if err != nil {
		exceptions.ThrowBadRequestException(ctx, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, common.BuildSuccessResponse("audit trails fetched successfully", auditTrails))
}

// GetDetails implements controller_interfaces.InvoiceController.
func (i *invoiceController) GetDetails(ctx *gin.Context) {
	invoiceID := ctx.Param("invoice_id")

	if invoiceID == "" {
		exceptions.ThrowUnProcessableEntityException(ctx, "invoice id is required")
		return
	}

	invoiceIDUint, err := strconv.ParseUint(invoiceID, 10, 64)
	if err != nil {
		exceptions.ThrowUnProcessableEntityException(ctx, "invalid invoice id")
		return
	}

	//TODO: call service to get invoice details
	invoiceDetails, err := i.invoiceService.GetInvoiceDetails(ctx, uint(invoiceIDUint))
	if err != nil {
		exceptions.ThrowBadRequestException(ctx, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, common.BuildSuccessResponse("invoice details fetched successfully", invoiceDetails))
}

// GetShareableLink implements controller_interfaces.InvoiceController.
func (i *invoiceController) GetShareableLink(ctx *gin.Context) {
	invoiceID := ctx.Param("invoice_id")

	if invoiceID == "" {
		exceptions.ThrowUnProcessableEntityException(ctx, "invoice id is required")
	}

	customer, err := i.getCustomerFromContext(ctx)
	if err != nil {
		exceptions.ThrowBadRequestException(ctx, err.Error())
	}

	invoiceIDUint, err := strconv.ParseUint(invoiceID, 10, 64)
	if err != nil {
		exceptions.ThrowUnProcessableEntityException(ctx, "invalid invoice id")
	}

	invoice, err := i.invoiceService.GetInvoiceByIDandCustomer(ctx, uint(invoiceIDUint), customer.ID)
	if err != nil {
		exceptions.ThrowBadRequestException(ctx, err.Error())
	}

	//TODO: call service to get shareable link
	shareableLink, err := i.invoiceService.GetShareableLink(ctx, invoice)
	if err != nil {
		exceptions.ThrowBadRequestException(ctx, err.Error())
	}

	ctx.JSON(http.StatusOK, common.BuildSuccessResponse("shareable link fetched successfully", shareableLink))
}

// GetSingleInvoiceAuditTrails implements controller_interfaces.InvoiceController.
func (i *invoiceController) GetSingleInvoiceAuditTrails(ctx *gin.Context) {
	var request request_dto.GetAllRequest

	if err := ctx.ShouldBindQuery(&request); err != nil {
		exceptions.ThrowUnProcessableEntityException(ctx, err.Error())
	}

	invoiceID := ctx.Param("invoice_id")

	if invoiceID == "" {
		exceptions.ThrowUnProcessableEntityException(ctx, "invoice id is required")
	}

	invoiceIDUint, err := strconv.ParseUint(invoiceID, 10, 64)
	if err != nil {
		exceptions.ThrowUnProcessableEntityException(ctx, "invalid invoice id")
	}

	customer, err := i.getCustomerFromContext(ctx)
	if err != nil {
		exceptions.ThrowBadRequestException(ctx, err.Error())
	}

	//TODO: call service to get single invoice audit trails
	auditTrails, err := i.auditService.GetAuditTrailsByInvoiceID(ctx, uint(invoiceIDUint), customer.ID, request.Limit, request.Page)
	if err != nil {
		exceptions.ThrowBadRequestException(ctx, err.Error())
	}

	ctx.JSON(http.StatusOK, common.BuildSuccessResponse("audit trails fetched successfully", auditTrails))
}

// GetStatistics implements controller_interfaces.InvoiceController.
func (i *invoiceController) GetStatistics(ctx *gin.Context) {
	// TODO: call service to get invoice statistics
	customer, err := i.getCustomerFromContext(ctx)
	if err != nil {
		exceptions.ThrowBadRequestException(ctx, err.Error())
		return
	}

	statistics, err := i.invoiceService.GetInvoiceStatistics(ctx, customer.ID)
	if err != nil {
		exceptions.ThrowBadRequestException(ctx, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, common.BuildSuccessResponse("invoice statistics fetched successfully", statistics))
}

// SetReminder implements controller_interfaces.InvoiceController.
func (i *invoiceController) SetReminder(ctx *gin.Context) {
	var request request_dto.SetReminderRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		exceptions.ThrowUnProcessableEntityException(ctx, err.Error())
		return
	}

	// TODO: call service to set reminder
	customer, err := i.getCustomerFromContext(ctx)
	if err != nil {
		exceptions.ThrowBadRequestException(ctx, err.Error())
		return
	}

	invoiceID := ctx.Param("invoice_id")
	if invoiceID == "" {
		exceptions.ThrowUnProcessableEntityException(ctx, "invoice id is required")
		return
	}

	invoiceIDUint, err := strconv.ParseUint(invoiceID, 10, 64)
	if err != nil {
		exceptions.ThrowUnProcessableEntityException(ctx, "invalid invoice id")
		return
	}

	invoice, err := i.invoiceService.GetInvoiceByIDandCustomer(ctx, uint(invoiceIDUint), customer.ID)
	if err != nil {
		exceptions.ThrowBadRequestException(ctx, err.Error())
		return
	}

	err = i.reminderService.SetInvoiceReminders(ctx, invoice, customer.ID, request.Schedules)
	if err != nil {
		exceptions.ThrowBadRequestException(ctx, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, common.BuildSuccessResponse("reminders set successfully", nil))
}

func (i *invoiceController) getCustomerFromContext(ctx *gin.Context) (*models.Customer, error) {
	customerID, err := helper.GetCustomerIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	return i.customerService.GetCustomerByID(ctx, customerID)
}

func (i *invoiceController) getInvoiceDetailsFromParams(ctx *gin.Context, customerID uint) (*models.Invoice, error) {
	invoiceID := ctx.Param("invoice_id")
	if invoiceID == "" {
		return nil, errors.New("invoice id is required")
	}

	invoiceIDUint, err := strconv.ParseUint(invoiceID, 10, 64)
	if err != nil {
		return nil, errors.New("invalid invoice id")
	}

	return i.invoiceService.GetInvoiceByIDandCustomer(ctx, uint(invoiceIDUint), customerID)
}

func NewInvoiceController(
	logger *zerolog.Logger,
	invoiceService services_interfaces.InvoiceService,
	auditService services_interfaces.AuditService,
	reminderService services_interfaces.RemiderService,
	customerService services_interfaces.CustomerService,
) controller_interfaces.InvoiceController {
	return &invoiceController{
		logger:          logger,
		invoiceService:  invoiceService,
		auditService:    auditService,
		reminderService: reminderService,
		customerService: customerService,
	}
}
