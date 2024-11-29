package repositories

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Adebayobenjamin/numerisbook/pkg/models"
	repositories_interfaces "github.com/Adebayobenjamin/numerisbook/pkg/repositories/interfaces"
	"github.com/jmoiron/sqlx"
)

type customerRepository struct {
	db *sqlx.DB
}

// GetCustomerByID implements repositories_interfaces.CustomerRepository.
func (c *customerRepository) GetCustomerByID(ctx context.Context, customerID uint) (*models.Customer, error) {
	query := `
		SELECT * FROM customers 
		WHERE id = ? AND deleted_at IS NULL`

	var customer models.Customer
	err := c.db.GetContext(ctx, &customer, query, customerID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("customer not found")
		}
		return nil, fmt.Errorf("failed to get customer: %w", err)
	}

	return &customer, nil
}

func NewCustomerRepository(db *sqlx.DB) repositories_interfaces.CustomerRepository {
	return &customerRepository{
		db: db,
	}
}
