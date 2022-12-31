package adapter

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
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

func (m *MariaAdapter) Connect() *sql.DB {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", m.Username, m.Password, m.Host, m.Port, m.Database))
	if err != nil {
		log.Fatal(err)
	}
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(100)
	db.SetMaxIdleConns(100)
	return db
}



// select data from mysql
func (m *MariaAdapter) Select(selectSQL []string, fromSQL string) []interface{} {
	db := m.Connect()
	defer db.Close()
	queryString := fmt.Sprintf("SELECT %s FROM %s", strings.Join(selectSQL, ","), fromSQL)
	rows, err := db.Query(queryString)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		log.Fatal(err)
	}

	count := len(columns)
	tableData := make([]interface{}, count)
	values := make([]interface{}, count)
	scanArgs := make([]interface{}, count)
	for i := range values {
		scanArgs[i] = &values[i]
	}

	for rows.Next() {
		err = rows.Scan(scanArgs...)
		if err != nil {
			log.Fatal(err)
		}
		for i, col := range values {
			if col != nil {
				tableData[i] = col
			}
		}
	}
	fmt.Println(tableData)
	return tableData
}


func (m *MariaAdapter) Update(updateSQL []string, fromSQL string) {
	db := m.Connect()
	defer db.Close()
	queryString := fmt.Sprintf("UPDATE %s SET %s", fromSQL, strings.Join(updateSQL, ","))
	_, err := db.Query(queryString)
	if err != nil {
		log.Fatal(err)
	}
}


func (m *MariaAdapter) Insert(insertSQL []string, fromSQL string) {
	db := m.Connect()
	defer db.Close()
	queryString := fmt.Sprintf("INSERT INTO %s VALUES %s", fromSQL, strings.Join(insertSQL, ","))
	_, err := db.Query(queryString)
	if err != nil {
		log.Fatal(err)
	}
}


func (m *MariaAdapter) Delete(deleteSQL []string, fromSQL string) {
	db := m.Connect()
	defer db.Close()
	queryString := fmt.Sprintf("DELETE FROM %s WHERE %s", fromSQL, strings.Join(deleteSQL, ","))
	_, err := db.Query(queryString)
	if err != nil {
		log.Fatal(err)
	}
}

