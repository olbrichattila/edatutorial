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

func (r *repository) GetInvoiceId(orderUUID string) (int64, error) {
	sql := "SELECT id FROM invoice_heads WHERE order_uuid = ?"

	rows, err := dbexecutor.RunSelectSQL(r.db, sql, orderUUID)
	if err != nil {
		return 0, err
	}

	if len(rows) == 0 {
		return 0, fmt.Errorf("invoice does not exists")
	}

	return rows[0]["id"].(int64), nil
}

// Episode 004 orderUUID added
func (r *repository) CreateInvoice(orderUUID string) (invoiceID int64, err error) {
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

	invoiceID, err = r.createInvoice(tx, orderUUID)

	return
}

func (r *repository) createInvoice(tx *sql.Tx, orderUUID string) (int64, error) {
	sql := "SELECT * FROM order_heads WHERE uuid = ? AND cancelled = 0"

	rows, err := dbexecutor.RunSelectSQL(tx, sql, orderUUID)
	if err != nil {
		return 0, err
	}

	if len(rows) == 0 {
		return 0, fmt.Errorf("nothing to store")
	}

	// Episode 004 add order UUID as de-dupe key
	invoiceID, err := r.createInvoiceHead(tx, rows[0], orderUUID)
	if err != nil {
		return 0, err
	}

	orderID, ok := rows[0]["id"].(int64)
	if !ok {
		return 0, fmt.Errorf("cannot get order ID")
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

// Episode 004 add order UUID as de-dupe key
func (r *repository) createInvoiceHead(tx *sql.Tx, row map[string]any, orderUUID string) (int64, error) {
	sql := `INSERT INTO invoice_heads (user_id, email, order_uuid) VALUES (?, ?, ?)`

	// Episode 004 added order UUID
	return dbexecutor.ExecuteInsertSQL(tx, sql, row["user_id"], row["email"], orderUUID)
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
