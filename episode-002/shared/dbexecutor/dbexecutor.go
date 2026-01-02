package dbexecutor

import (
	"database/sql"
	"fmt"

	"github.com/olbrichattila/edatutorial/shared/config"

	_ "github.com/go-sql-driver/mysql"
)

// This encapsulates *sql.DB and *SQL.tx common properties I use to be able
// to work with Unit of Work pattern
type DbExecutor interface {
	Query(query string, args ...any) (*sql.Rows, error)
	Exec(query string, args ...any) (sql.Result, error)
}

func RunSelectSQL(db DbExecutor, sql string, params ...any) ([]map[string]any, error) {
	rows, err := db.Query(sql, params...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	results := make([]map[string]any, 0)
	for rows.Next() {
		result := make(map[string]any, len(columns))

		values := make([]any, len(columns))
		pointers := make([]any, len(columns))

		for i := range columns {
			pointers[i] = &values[i]
		}

		err := rows.Scan(pointers...)
		if err != nil {
			return nil, err
		}

		for i, fieldName := range columns {
			result[fieldName] = values[i]
		}

		results = append(results, result)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return results, nil
}

func ExecuteInsertSQL(db DbExecutor, sql string, params ...any) (int64, error) {
	result, err := db.Exec(sql, params...)
	if err != nil {
		return 0, err
	}

	lastInsertID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return lastInsertID, nil
}

func ExecuteUpdateSQL(db DbExecutor, sql string, params ...any) (int64, error) {
	result, err := db.Exec(sql, params...)
	if err != nil {
		return 0, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return rowsAffected, nil
}

func ConnectToDB() (*sql.DB, error) {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=true",
		config.DBUsername(),
		config.DBPassword(),
		config.DBHost(),
		config.DBPort(),
		config.DBDatabase(),
	)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
