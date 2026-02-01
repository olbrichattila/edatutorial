package order

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/olbrichattila/edatutorial/shared/dbexecutor"
	"producer.example/internal/contracts"
)

var (
	ErrRaceConditionDetected = errors.New("race condition detected")
)

func New(db *sql.DB) contracts.OrderRepository {
	return &repository{
		db: db,
	}
}

type repository struct {
	db *sql.DB
}

func (r *repository) Cancel(orderUUID string) (err error) {
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

	err = r.cancelOrder(tx, orderUUID)

	return
}

func (r *repository) cancelOrder(tx *sql.Tx, orderUUID string) error {
	sql := "UPDATE order_heads SET cancelled = 1 WHERE uuid = ?"

	_, err := dbexecutor.ExecuteUpdateSQL(tx, sql, orderUUID)
	if err != nil {
		return err
	}

	orderID, err := r.getOrderIdByUUID(tx, orderUUID)
	if err != nil {
		return err
	}

	return r.updateStockPerItem(tx, orderID)
}

func (r *repository) getOrderIdByUUID(tx *sql.Tx, orderUUID string) (int64, error) {
	sql := "SELECT id FROM order_heads WHERE uuid = ?"

	result, err := dbexecutor.RunSelectSQL(tx, sql, orderUUID)
	if err != nil {
		return 0, err
	}

	if len(result) == 0 {
		return 0, fmt.Errorf("cannot find order id for order %s", orderUUID)
	}

	if id, ok := result[0]["id"]; ok {
		if int64Id, ok := id.(int64); ok {
			return int64Id, nil
		}
	}

	return 0, fmt.Errorf("cannot cast order id for order %s", orderUUID)
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

// Episode 005 Refactored to handle versions, race condition optimistically
func (r *repository) updateStock(tx *sql.Tx, productID string, quantity int64) error {
	version, err := r.getCurrentVersion(tx, productID)
	if err != nil {
		return err
	}

	// Assume stock record already added when the product entered into inventory
	sql := `UPDATE stocks 
				SET quantity = quantity + ?,
				version = version + 1
			WHERE
				product_id = ?
				AND version = ?`

	rowsAffected, err := dbexecutor.ExecuteUpdateSQL(
		tx,
		sql,
		quantity,
		productID,
		version,
	)
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrRaceConditionDetected
	}

	return nil
}

// Episode 005, get current quantity and version for optimistic race condition handling,
// - business invariant
// - database level race condition
func (r *repository) getCurrentVersion(tx *sql.Tx, productID string) (int, error) {
	sql := `SELECT
				version FROM stocks 
			WHERE
				product_id = ?`

	var version int
	err := tx.QueryRow(sql, productID).Scan(&version)
	if err != nil {
		return 0, fmt.Errorf("cannot fetch from database {%w}", err)
	}

	return version, nil
}
