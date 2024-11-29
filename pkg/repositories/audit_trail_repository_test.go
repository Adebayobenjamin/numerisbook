package repositories

import (
	"context"
	"log"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"
)

func getMockDB() (*sqlx.DB, sqlmock.Sqlmock, *auditTrailRepository) {
	mockDB, mock, err := sqlmock.New()
	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")
	if err != nil {
		log.Fatalf("Failed to create mock: %v", err)
	}

	repo := &auditTrailRepository{
		db:     sqlxDB,
		logger: &zerolog.Logger{},
	}

	return sqlxDB, mock, repo
}
func TestAuditTrailRepository_GetAllCustomerAuditTrails(t *testing.T) {
	db, mock, repo := getMockDB()

	defer db.Close()

	customerID := uint(1)
	limit := 10
	offset := 0

	rows := sqlmock.NewRows([]string{"id", "event_type", "log_level", "message", "invoice_id", "customer_id", "created_at", "updated_at"}).
		AddRow(1, "INVOICE_CREATED", "INFO", "Invoice created", 1, customerID, time.Now(), time.Now())

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM audit_trails WHERE customer_id = ? AND deleted_at IS NULL ORDER BY created_at DESC LIMIT ? OFFSET ?`)).
		WithArgs(customerID, limit, offset).
		WillReturnRows(rows)

	repo.GetAllCustomerAuditTrails(context.Background(), customerID, limit, offset)
}
