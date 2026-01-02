package invoice

import (
	"database/sql"
	"fmt"
	"math/rand"

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

func (r *repository) CreateInvoice(orderID int64) (invoiceID int64, err error) {
	tx, err := r.db.Begin()
	if err != nil {
		return
	}

	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}

		tx.Commit()
	}()

	invoiceID, err = r.createInvoice(tx, orderID)

	return
}

func (r *repository) createInvoice(tx *sql.Tx, orderID int64) (int64, error) {
	sql := "SELECT * FROM order_heads WHERE id = ? AND cancelled = 0"

	rows, err := dbexecutor.RunSelectSQL(tx, sql, orderID)
	if err != nil {
		return 0, err
	}

	if len(rows) == 0 {
		return 0, fmt.Errorf("nothing to store")
	}

	invoiceID, err := r.createInvoiceHead(tx, rows[0])
	if err != nil {
		return 0, err
	}

	sql = "SELECT * FROM order_items WHERE order_id = ?"
	rows, err = dbexecutor.RunSelectSQL(tx, sql, orderID)
	if err != nil {
		return 0, err
	}

	for _, row := range rows {
		err := r.createInvoiceItem(tx, invoiceID, row)
		if err != nil {
			return 0, err
		}
	}

	return invoiceID, nil
}

func (r *repository) createInvoiceHead(tx *sql.Tx, row map[string]any) (int64, error) {
	sql := `INSERT INTO invoice_heads (user_id, email) VALUES (?, ?)`

	return dbexecutor.ExecuteInsertSQL(tx, sql, row["user_id"], row["email"])
}

func (r *repository) createInvoiceItem(tx *sql.Tx, invoiceID int64, orderRow map[string]any) error {
	sql := `INSERT INTO invoice_items (
		order_id,
		invoice_id,
		product_id,
		quantity,
		price
	) VALUES (?, ?, ?, ?, ?)`

	randomPrice := float64(rand.Int63n(100000)) / 100

	_, err := dbexecutor.ExecuteInsertSQL(
		tx,
		sql,
		orderRow["order_id"],
		invoiceID,
		orderRow["product_id"],
		orderRow["quantity"],
		randomPrice,
	)
	return err
}
