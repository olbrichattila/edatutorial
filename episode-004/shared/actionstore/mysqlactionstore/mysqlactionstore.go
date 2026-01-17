// Episode 004 Action store for idempotency
package mysqlactionstore

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"

	"github.com/go-sql-driver/mysql"
	"github.com/olbrichattila/edatutorial/shared/actions"
	"github.com/olbrichattila/edatutorial/shared/actionstore/contracts"
)

func New(db *sql.DB) (contracts.ActionStore, error) {
	if db == nil {
		return nil, fmt.Errorf("mysql action store, db required")
	}
	return &actionStore{
		db: db,
	}, nil
}

type actionStore struct {
	db *sql.DB
}

func (a *actionStore) IsDuplicateAction(actionMetaData actions.MetaData) (bool, error) {
	_, err := a.db.Exec(
		"INSERT INTO processed_actions (idempotency_key) VALUES (?)",
		a.idempotencyKey(actionMetaData),
	)
	if a.isDuplicateKeyError(err) {
		return true, nil
	}

	return false, err
}

func (a *actionStore) idempotencyKey(actionMetaData actions.MetaData) string {
	return actionMetaData.CorrelationID + "_" + actionMetaData.CausationID + "_" + strconv.FormatInt(actionMetaData.Index, 10)
}

func (a *actionStore) isDuplicateKeyError(err error) bool {
	var mysqlErr *mysql.MySQLError
	if errors.As(err, &mysqlErr) {
		return mysqlErr.Number == 1062
	}
	return false
}
