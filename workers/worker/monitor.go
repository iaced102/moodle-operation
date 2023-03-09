package worker

import (
	"database/sql"
	repo "moodle/internal/repository/moodle"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

type MonitorWorker struct {
	mariaRepo *repo.MariaDB
}

// new mariadb worker
func NewMonitorWorker(maria *sql.DB) *MonitorWorker {
	mariaRepo := repo.NewMariaDB(maria)
	
	return &MonitorWorker{
		mariaRepo: mariaRepo,
	}
}

// tracking user from mdl_user
func (w *MonitorWorker) GetUser() error {
	// get all dbname
	dbnames, err := w.mariaRepo.GetDBName()
	if err != nil {
		return err
	}
	// for  database in dbnames get userdata
	for _, dbname := range dbnames {
		go w.WriteUser(dbname)
	}
	return nil
}

// getuser data then write to influxdb to monitoring
func (w *MonitorWorker) WriteUser(dbname string) error {
	// get user
	data , err := w.mariaRepo.CountUserOnline(dbname)
	if err != nil {
		return err
	}
	// write to influxdb
	Write(dbname, data)
	return nil
}


func Write(dbname string, data []int) error {
    // Create a client to connect to InfluxDB
    client := influxdb2.NewClientWithOptions("http://123.30.234.141:8086", "lmspoller:lmspoller", influxdb2.DefaultOptions().SetBatchSize(100))

    // Create a write API for the "mydb" database
    writeAPI := client.WriteAPI("lms", "lms")

    // Define a data point to write
    p := influxdb2.NewPointWithMeasurement("lmsuser").
		AddField("online_user",(data[0])).
		AddField("dbname", dbname).
		AddTag("dbname_tag", dbname).
        SetTime(time.Now())

    // Write the data point to InfluxDB
    writeAPI.WritePoint(p)

    // Close the client when you're done
    client.Close()
	return nil
}
