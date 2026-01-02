package order

import (
	"database/sql"
	"fmt"

	"github.com/olbrichattila/edatutorial/shared/dbexecutor"
	"producer.example/internal/contracts"
)

func New(db *sql.DB) contracts.OrderRepository {
	return &repository{
		db: db,
	}
}

type repository struct {
	db *sql.DB
}

func (r *repository) Cancel(orderID int64) (err error) {
	tx, err := r.db.Begin()
	if err != nil {
		return
	}

	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				fmt.Printf("rollback error: %v\n", rollbackErr)
			}

			return
		}

		if commitErr := tx.Commit(); commitErr != nil {
			err = commitErr
		}
	}()

	err = r.cancelOrder(tx, orderID)

	return
}

func (r *repository) cancelOrder(tx *sql.Tx, orderID int64) error {
	sql := "UPDATE order_heads set cancelled = 1 WHERE id = ?"

	_, err := dbexecutor.ExecuteUpdateSQL(tx, sql, orderID)
	if err != nil {
		return err
	}

	return r.updateStockPerItem(tx, orderID)
}

func (r *repository) updateStockPerItem(tx *sql.Tx, orderID int64) error {
	sql := `SELECT product_id, quantity FROM order_items WHERE order_id = ?`
	rows, err := dbexecutor.RunSelectSQL(tx, sql, orderID)
	if err != nil {
		return err
	}

	for _, row := range rows {
		productID, ok := row["product_id"].([]uint8)
		if !ok {
			return fmt.Errorf("productID is not []unit8")
		}

		quantity, ok := row["quantity"].(int64)
		if !ok {
			return fmt.Errorf("quantity is not int64")
		}

		err := r.updateStock(tx, string(productID), quantity)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *repository) updateStock(tx *sql.Tx, productID string, quantity int64) error {
	sql := `INSERT INTO stocks (product_id, quantity)
		VALUES (?, ?) AS new
		ON DUPLICATE KEY UPDATE
			quantity = stocks.quantity + new.quantity`

	_, err := dbexecutor.ExecuteUpdateSQL(tx, sql, productID, quantity)
	if err != nil {
		return err
	}

	return nil
}
