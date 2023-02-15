package mariadbiface

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

// MariaDB interface
type MariaDB interface {
	Connect(dbname string) *sql.DB
}
