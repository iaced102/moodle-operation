package worker

import (
	"database/sql"
	"fmt"
	"moodle/config"
	repo "moodle/internal/repository/moodle"
	"moodle/pkg/mongodbiface"
	"net/http"
	"strings"
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
		return nil
	}
	// if len(dbnames) == 0  then return nil
	if len(dbnames) == 0 {
		return nil
	}
	// for  database in dbnames get userdata
	for _, dbname := range dbnames {
		// ignore if dbname is moodle
		if dbname == "moodle" {
			continue
		}
		go w.WriteUser(dbname)
	}
	return nil
}

// getuser data then write to influxdb to monitoring
func (w *MonitorWorker) WriteUser(dbname string) error {
	// get user
	data , err := w.mariaRepo.CountUserOnline(dbname)
	if err != nil {
		return nil
	}
	// if len(data) == 0 then return nil
	if len(data) == 0 {
		data = append(data, 0)
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

	url := config.VICTORIAMETRIC
	measurementName := "lmsccu"
	tagKey := "dbname"
	tagValue := dbname
	fieldKey := "online_user"
	fieldValue := data[0]
	fieldKey2 := "ccu"
	fieldValue2 := ccu
	fieldKey3 := "ccumax"
	fieldValue3 := ccumax
	// timestamp := time.Now().Unix()

	metrics := fmt.Sprintf("%s,%s=%s %s=%d,%s=%d,%s=%d", measurementName, tagKey, tagValue, fieldKey, fieldValue, fieldKey2, fieldValue2, fieldKey3, fieldValue3)
	resp, err := http.Post(url, "application/octet-stream", strings.NewReader(metrics))
	if err != nil {
		return nil
		// panic(err)
	}
	defer resp.Body.Close()

	return nil
}
