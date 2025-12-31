package order

import (
	"database/sql"

	"github.com/olbrichattila/edatutorial/shared/dbexecutor"
	"producer.example/internal/contracts"
	"producer.example/internal/dto"
)

func New(db *sql.DB) contracts.OrderRepository {
	return &repository{
		db: db,
	}
}

type repository struct {
	db *sql.DB
}

func (r *repository) Save(ord dto.Order) (lastInsertId int64, err error) {
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

	lastInsertId, err = r.saveHead(tx, ord)
	if err != nil {
		return 0, err
	}

	err = r.saveItems(tx, lastInsertId, ord.Items)
	if err != nil {
		return 0, err
	}

	return
}

func (r *repository) saveHead(tx *sql.Tx, ord dto.Order) (int64, error) {
	sql := "INSERT INTO order_heads (user_id, email) VALUES (?, ?)"

	lastInsertID, err := dbexecutor.ExecuteInsertSQL(tx, sql, ord.UserID, ord.Email)
	if err != nil {
		return 0, err
	}

	return lastInsertID, nil
}

func (r *repository) saveItems(tx *sql.Tx, orderHeadId int64, items []dto.Item) error {
	sql := "INSERT INTO order_items (order_id, product_id, quantity) VALUES (?,?,?)"

	for _, item := range items {
		_, err := dbexecutor.ExecuteInsertSQL(tx, sql, orderHeadId, item.ProductID, item.Quantity)
		if err != nil {
			return err
		}
	}

	return nil
}
