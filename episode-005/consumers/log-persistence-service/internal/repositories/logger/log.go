package logger

import (
	"database/sql"

	"github.com/olbrichattila/edatutorial/shared/dbexecutor"
	"producer.example/internal/contracts"
)

func New(db *sql.DB) contracts.LoggerRepository {
	return &repository{
		db: db,
	}
}

type repository struct {
	db *sql.DB
}

func (r *repository) Save(logType, actionID, correlationID, causationID, msg string, ind int64) error {
	sql := `INSERT INTO logs (level, action_id, correlation_id, causation_id, ind, message) VALUES (?,?,?,?,?,?)`

	_, err := dbexecutor.ExecuteInsertSQL(r.db, sql, logType, actionID, correlationID, causationID, ind, msg)

	return err
}
