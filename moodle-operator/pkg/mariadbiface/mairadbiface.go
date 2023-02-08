package mariadbiface

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

// MariaDB interface
type MariaDB interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
}
