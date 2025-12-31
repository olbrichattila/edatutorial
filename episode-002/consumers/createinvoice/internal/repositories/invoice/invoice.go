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

func (r *repository) CreateInvoice(orderId int64) (err error) {
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

	err = r.createInvoice(tx, orderId)

	return
}

func (r *repository) createInvoice(tx *sql.Tx, orderId int64) error {
	sql := "SELECT * FROM order_heads WHERE id = ? AND cancelled = 0"

	rows, err := dbexecutor.RunSelectSQL(tx, sql, orderId)
	if err != nil {
		return err
	}

	if len(rows) == 0 {
		return fmt.Errorf("Nothing to store")
	}

	err = r.createInvoiceHead(tx, rows[0])
	if err != nil {
		return err
	}

	sql = "SELECT * FROM order_items WHERE order_id = ? AND cancelled = 0"

	rows, err = dbexecutor.RunSelectSQL(tx, sql, orderId)
	if err != nil {
		return err
	}

	for _, row := range rows {
		err := r.createInvoiceItem(tx, row)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *repository) createInvoiceHead(tx *sql.Tx, row map[string]any) error {
	sql := `INSERT INTO order_heads (user_id, email) VALUES (?, ?)`

	_, err := dbexecutor.ExecuteInsertSQL(tx, sql, row["user_id"], row["email"])
	return err
}

// TODO store invoice head
// CREATE TABLE order_items (
// id INT AUTO_INCREMENT PRIMARY KEY,
// order_id INT NOT NULL,
// product_id CHAR(36) NOT NULL,
// quantity INT NOT NULL,
// created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
// FOREIGN KEY (order_id) REFERENCES order_heads(id)
// );
func (r *repository) createInvoiceItem(tx *sql.Tx, row map[string]any) error {
	return nil
}
