package services_interfaces

import (
	"context"

	"github.com/Adebayobenjamin/numerisbook/pkg/models"
)

type CustomerService interface {
	GetCustomerByID(ctx context.Context, customerID uint) (*models.Customer, error)
}
