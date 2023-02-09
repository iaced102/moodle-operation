package mariadbiface

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

// MariaDB interface
type MariaDB interface {
	Connect(dbname string) *sql.DB
	Select(dbname, query string) map[string]string
	CreateDatabase(dbname string) error
	RestoreDatabase(dbname, filepath string) error
	DropDatabase(dbname string) error
	Update(dbname, shortname, fullname string) error
}
