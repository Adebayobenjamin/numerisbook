package services

import (
	"context"

	"github.com/Adebayobenjamin/numerisbook/pkg/models"
	repositories_interfaces "github.com/Adebayobenjamin/numerisbook/pkg/repositories/interfaces"
	services_interfaces "github.com/Adebayobenjamin/numerisbook/pkg/services/interfaces"
	"github.com/rs/zerolog"
)

type customerService struct {
	logger             *zerolog.Logger
	customerRepository repositories_interfaces.CustomerRepository
}

// GetCustomerByID implements services_interfaces.CustomerService.
func (c *customerService) GetCustomerByID(ctx context.Context, customerID uint) (*models.Customer, error) {
	return c.customerRepository.GetCustomerByID(ctx, customerID)
}

func NewCustomerService(logger *zerolog.Logger, customerRepository repositories_interfaces.CustomerRepository) services_interfaces.CustomerService {
	return &customerService{
		logger:             logger,
		customerRepository: customerRepository,
	}
}
