package invoice

import (
	"database/sql"
	"fmt"

	"github.com/olbrichattila/edatutorial/shared/dbexecutor"
	"producer.example/internal/contracts"
)

func New(db *sql.DB) contracts.InvoiceRepository {
	return &repository{
		db: db,
	}
}

type repository struct {
	db *sql.DB
}

func (r *repository) Head(invoiceID int64) (map[string]any, error) {
	sql := `SELECT * FROM invoice_heads WHERE id = ?`
	result, err := dbexecutor.RunSelectSQL(r.db, sql, invoiceID)
	if err != nil {
		return nil, err
	}

	if len(result) == 0 {
		return nil, fmt.Errorf("Invoice is missing")
	}

	return result[0], nil
}

func (r *repository) Items(invoiceID int64) ([]map[string]any, error) {
	sql := `SELECT * FROM invoice_items WHERE invoice_id = ?`

	return dbexecutor.RunSelectSQL(r.db, sql, invoiceID)
}
