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

func (r *repository) Save(msg string) error {
	sql := `INSERT INTO logs (message) VALUES (?)`
	_, err := dbexecutor.ExecuteInsertSQL(r.db, sql, msg)

	return err
}
