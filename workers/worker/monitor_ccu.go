package worker

import (
	"database/sql"
	repo "moodle/internal/repository/moodle"
	"moodle/pkg/mongodbiface"
	"strings"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

type MonitorWorker struct {
	mariaRepo *repo.MariaDB
	mongoRepo *repo.MongoDB
}

func NewMonitorWorker(mongo mongodbiface.DB, maria *sql.DB) *MonitorWorker {
	mariaRepo := repo.NewMariaDB(maria)
	mongoRepo := repo.NewMongoDB(mongo)
	
	return &MonitorWorker{
		mariaRepo: mariaRepo,
		mongoRepo: mongoRepo,
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
	w.Write(dbname, data)
	return nil
}


func (w *MonitorWorker) Write(dbname string, data []int) error {
	moodle, err := w.mongoRepo.Get(strings.ReplaceAll(dbname, "_", "-"))
	// check if error is no documents in result then ignore
	if err != nil {
		moodle.Packages.Ccu = 0
		moodle.Packages.CcuExtraMax = 0
	}
	ccu := moodle.Packages.Ccu
	ccumax := moodle.Packages.CcuExtraMax

    // Create a client to connect to InfluxDB
    client := influxdb2.NewClientWithOptions("http://123.30.234.141:8086", "lmspoller:lmspoller", influxdb2.DefaultOptions().SetBatchSize(100))

    // Create a write API for the "mydb" database
    writeAPI := client.WriteAPI("lms", "lms")

    // Define a data point to write
    p := influxdb2.NewPointWithMeasurement("lmsuser").
		AddField("online_user",(data[0])).
		AddField("dbname", dbname).
		AddField("ccu", ccu).
		AddField("ccumax", ccumax).
		AddTag("dbname_tag", dbname).
        SetTime(time.Now())

    // Write the data point to InfluxDB
    writeAPI.WritePoint(p)

    // Close the client when you're done
    client.Close()
	return nil
}
