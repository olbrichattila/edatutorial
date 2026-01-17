package order

import (
	"database/sql"

	"github.com/olbrichattila/edatutorial/shared/actions"
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
	reversedQuantity := int(-quantity)

	sql := `INSERT INTO stocks (product_id, quantity)
		VALUES (?, ?) AS new
		ON DUPLICATE KEY UPDATE
			quantity = stocks.quantity + new.quantity`

	_, err := dbexecutor.ExecuteUpdateSQL(tx, sql, productID, reversedQuantity)
	if err != nil {
		return err
	}

	return nil
}
