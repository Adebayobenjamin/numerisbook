package repositories_interfaces

import (
	"context"

	"github.com/Adebayobenjamin/numerisbook/pkg/models"
)

type CustomerRepository interface {
	GetCustomerByID(ctx context.Context, customerID uint) (*models.Customer, error)
}
