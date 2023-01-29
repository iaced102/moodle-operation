package adapter

import (
	"database/sql"
	"fmt"
	"log"
	"os/exec"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type MariaAdapter struct {
	Host     string
	Port     int
	Username string
	Password string
	Database string
}

// new maria adapter
func NewMariaAdapter(host string, port int, username string, password string) *MariaAdapter {
	return &MariaAdapter{
		Host:     host,
		Port:     port,
		Username: username,
		Password: password,
	}
}

func (m *MariaAdapter) Connect(dbname string) *sql.DB {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", m.Username, m.Password, m.Host, m.Port, dbname))
	if err != nil {
		log.Fatal(err)
	}
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(100)
	db.SetMaxIdleConns(100)
	return db
}



// select data from mysql
func (m *MariaAdapter) Select(dbname, query string) map[string]string {
	db := m.Connect(dbname)
	defer db.Close()
	queryString := fmt.Sprintf("%s", query)
	rows, err := db.Query(queryString)
	if err != nil {
		log.Println(err)
		return nil
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		log.Println(err)
		return nil
	}
	// map column and column data as key and value
	columnsMap := make(map[string]string)
	// make a slice for the values
	values := make([]sql.RawBytes, len(columns))
	// rows.Scan wants '[]interface{}' as an argument, so we must copy the
	// references into such a slice
	// see http://code.google.com/p/go-wiki/wiki/InterfaceSlice for details
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}
	// Fetch rows
	for rows.Next() {
		// get RawBytes from data
		err = rows.Scan(scanArgs...)
		if err != nil {
			log.Println(err)
			return nil
		}
		// Now do something with the data.
		// Here we just print each column as a string.
		var value string
		for i, col := range values {
			// Here we can check if the value is nil (NULL value)
			if col == nil {
				value = "NULL"
			} else {
				value = string(col)
			}
			columnsMap[columns[i]] = value
		}
	}
	if err = rows.Err(); err != nil {
		log.Println(err)
		return nil
	}
	return columnsMap
}


// create new database from existing Database
func (m *MariaAdapter) CreateDatabase(dbname string) error {
	db := m.Connect("")
	defer db.Close()
	queryString := fmt.Sprintf("CREATE DATABASE %s", dbname)
	_, err := db.Exec(queryString)
	if err != nil {
		log.Println(err)
		return err
	}
	return err
}

// restore database from sql filepath
func (m *MariaAdapter) RestoreDatabase(dbname, filepath string) error {
	db := m.Connect(dbname)
	defer db.Close()

	comand :=  "mysql -u root -h 45.124.94.112 -p0YU8381WUlk1u9ysVbF4Qb5FigNW8z8uCvPI " + dbname + " < " + filepath
	// print current time
	log.Println("start at: ", time.Now())
	output, err := exec.Command("bash", "-c", comand).Output()
	log.Println("end at: ", time.Now())

	if err != nil {
		log.Println(err)
		return err
	}
	log.Println(string(output))
	return err
}
