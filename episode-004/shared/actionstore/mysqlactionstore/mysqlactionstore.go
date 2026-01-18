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
	"github.com/olbrichattila/edatutorial/shared/dbexecutor"
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

func (a *actionStore) GetMetadata(consumer string, actionMetaData actions.MetaData) (string, error) {
	panic("unimplemented")
}

func (a *actionStore) SetDuplicateAction(consumer string, actionMetaData actions.MetaData, metadata string) error {
	_, err := a.db.Exec(
		"INSERT INTO processed_actions (idempotency_key, metadata) VALUES (?,?)",
		a.idempotencyKey(consumer, actionMetaData),
		metadata,
	)

	if a.isDuplicateKeyError(err) {
		return nil
	}

	return nil
}

func (a *actionStore) IsDuplicateAction(consumer string, actionMetaData actions.MetaData) (bool, string, error) {
	rows, err := dbexecutor.RunSelectSQL(
		a.db,
		"SELECT metadata FROM processed_actions WHERE idempotency_key = ?",
		a.idempotencyKey(consumer, actionMetaData),
	)
	if err != nil {
		return false, "", err
	}

	if len(rows) > 0 {
		return true, string(rows[0]["metadata"].([]byte)), nil
	}

	return false, "", nil
}

func (a *actionStore) idempotencyKey(consumer string, actionMetaData actions.MetaData) string {
	return consumer + "_" + actionMetaData.ActionType + "_" + actionMetaData.CorrelationID + "_" + strconv.FormatInt(actionMetaData.Index, 10)
}

func (a *actionStore) isDuplicateKeyError(err error) bool {
	var mysqlErr *mysql.MySQLError
	if errors.As(err, &mysqlErr) {
		return mysqlErr.Number == 1062
	}
	return false
}
