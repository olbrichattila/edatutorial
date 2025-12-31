package order

import (
	"database/sql"

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

func (r *repository) Cancel(orderId int64) (err error) {
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

	err = r.cancelOrder(tx, orderId)

	return
}

func (r *repository) cancelOrder(tx *sql.Tx, orderId int64) error {
	sql := "UPDATE order_heads set cancelled = 1 WHERE id = ?"

	_, err := dbexecutor.ExecuteUpdateSQL(tx, sql, orderId)
	if err != nil {
		return err
	}

	return r.updateStockPerItem(tx, orderId)
}

func (r *repository) updateStockPerItem(tx *sql.Tx, orderId int64) error {
	sql := `SELECT product_id, quantity FROM order_items WHERE order_id = ?`
	rows, err := dbexecutor.RunSelectSQL(tx, sql, orderId)
	if err != nil {
		return err
	}

	for _, row := range rows {
		err := r.updateStock(tx, string(row["product_id"].([]uint8)), int(row["quantity"].(int64)))
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *repository) updateStock(tx *sql.Tx, productId string, quantity int) error {
	sql := `INSERT INTO stocks (product_id, quantity)
		VALUES (?, ?) AS new
		ON DUPLICATE KEY UPDATE
			quantity = stocks.quantity + new.quantity`

	_, err := dbexecutor.ExecuteUpdateSQL(tx, sql, productId, quantity)
	if err != nil {
		return err
	}

	return nil
}
