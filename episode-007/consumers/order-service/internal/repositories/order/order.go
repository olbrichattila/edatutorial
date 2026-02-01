package order

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/olbrichattila/edatutorial/shared/actions"
	"github.com/olbrichattila/edatutorial/shared/dbexecutor"
	"producer.example/internal/contracts"
)

var (
	ErrRaceConditionDetected = errors.New("race condition detected")
	ErrNotEnoughStock        = errors.New("not enough stock")
)

func New(db *sql.DB) contracts.OrderRepository {
	return &repository{
		db: db,
	}
}

type repository struct {
	db *sql.DB
}

func (r *repository) Save(ord actions.OrderCreatedAction) (lastInsertID int64, err error) {
	tx, err := r.db.Begin()
	if err != nil {
		return
	}

	defer func() {
		if err != nil {
			rollbackErr := tx.Rollback()
			if rollbackErr != nil {
				err = rollbackErr
			}
			return
		}

		if commitErr := tx.Commit(); commitErr != nil {
			err = commitErr
		}
	}()

	lastInsertID, err = r.saveHead(tx, ord)
	if err != nil {
		return 0, err
	}

	err = r.saveItems(tx, lastInsertID, ord.Items)
	if err != nil {
		return 0, err
	}

	return
}

func (r *repository) saveHead(tx *sql.Tx, ord actions.OrderCreatedAction) (int64, error) {
	//  Episode 003 Add idempotency uuid
	sql := "INSERT INTO order_heads (uuid, user_id, email) VALUES (?, ?, ?)"

	lastInsertID, err := dbexecutor.ExecuteInsertSQL(tx, sql, ord.UUID, ord.UserID, ord.Email)
	if err != nil {
		return 0, err
	}

	return lastInsertID, nil
}

func (r *repository) saveItems(tx *sql.Tx, orderHeadID int64, items []actions.OrderItem) error {
	sql := "INSERT INTO order_items (order_id, product_id, quantity) VALUES (?,?,?)"

	for _, item := range items {
		_, err := dbexecutor.ExecuteInsertSQL(tx, sql, orderHeadID, item.ProductID, item.Quantity)
		if err != nil {
			return err
		}

		err = r.updateStock(tx, item.ProductID, item.Quantity)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *repository) updateStock(tx *sql.Tx, productID string, quantity uint) error {
	// Episode 005, now we assume that the stock row exists, added when the product was added to the system
	// and use optimistic versioning to handle race conditions
	currentQuantity, version, err := r.getCurrentStockAndVersion(tx, productID)
	if err != nil {
		return err
	}

	if quantity > currentQuantity {
		return ErrNotEnoughStock
	}

	sql := `UPDATE stocks 
				SET quantity = quantity - ?,
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
func (r *repository) getCurrentStockAndVersion(tx *sql.Tx, productID string) (uint, int, error) {
	sql := `SELECT
				quantity, version FROM stocks 
			WHERE
				product_id = ?`

	var quantity uint
	var version int
	err := tx.QueryRow(sql, productID).Scan(&quantity, &version)
	if err != nil {
		return 0, 0, fmt.Errorf("cannot fetch from database {%w}", err)
	}

	return quantity, version, nil
}
