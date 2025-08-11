package worker

import (
	"database/sql"
	"fmt"
	"github.com/go-redis/redis"
	"log"
	"math"
	"moodle/config"
	repo "moodle/internal/repository/moodle"
	"moodle/pkg/mongodbiface"
	"net/http"
	"strings"
	"time"
	"strconv"
)

type MonitorWorker struct {
	mariaRepo *repo.MariaDB
	mongoRepo *repo.MongoDB
	redis     *redis.Client
}

func NewMonitorWorker(mongo mongodbiface.DB, maria *sql.DB, redis *redis.Client) *MonitorWorker {
	mariaRepo := repo.NewMariaDB(maria)
	mongoRepo := repo.NewMongoDB(mongo)

	return &MonitorWorker{
		mariaRepo: mariaRepo,
		mongoRepo: mongoRepo,
		redis:     redis,
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
		if dbname == "moodle" || dbname == "binlog" {
			continue
		}
		go w.WriteUser(dbname)
	}
	return nil
}

// getuser data then write to influxdb to monitoring
func (w *MonitorWorker) WriteUser(dbname string) error {
	// get user
	data, err := w.mariaRepo.CountUserOnline(dbname)
	if err != nil {
		return nil
	}
	// if len(data) == 0 then return nil
	if len(data) == 0 {
		return nil
	}

	//log.Printf("database %s - CCU: %d", dbname, data)
	// write to influxdb
	w.Write(dbname, data)
	//w.SetScale(dbname, data)
	return nil
}

func (w *MonitorWorker) SetScale(dbname string, current_ccu []int) error {
	// Set replicas pod number of statefulset to redis
	// replicas logic:
	// default with 2 pods
	// 1 pod can serve for 100 CCU, but using at least 2 pods for HA
	// if current ccu > 100 CCU, replicas = ceil (ccu / 100 )
	// scale logic:
	// if replicas > 2 => scale statefulset to number of replicas, keep this number in 60min
	// if replicas <= 2 but scale process is not due => keep old replicas number
	// if replicase <= 2 but scale process is due => scale
	// Get old CCU
	//if dbname != "16d53fd2_9f0c_48a4_921a_251009f9768a" {
	//	return nil
	//}

	ccu := current_ccu[0]

	// process dbname to statefulset name
	dbname = "replicas_" + strings.ReplaceAll(dbname, "_", "-")

	var replicas string
	var new_replicas int
	const BASE_CCU int = 200
	const UNIT_CCU float64 = 100.0
	const DEFAULT_REPLICAS int = 2

	// Default replicas = 2 pods
	if ccu > BASE_CCU {
		new_replicas = int(math.Ceil(float64(ccu) / UNIT_CCU))
	} else {
		new_replicas = DEFAULT_REPLICAS
	}

	replicas, err := w.redis.Get(dbname).Result()

	if err == redis.Nil {
		//log.Printf("Key store %d replicas does not exist", dbname)
		return nil
	} else if err != nil {
		//fmt.Println("Failed to get value:", err)
		return nil
	} else {
		log.Printf("Replicas number for %s is %d", dbname, replicas)
		i_replicas, _ := strconv.Atoi(replicas)
		if new_replicas < i_replicas {
			new_replicas = i_replicas
		}
	}

	if new_replicas > 2 {
		w.redis.Set(dbname, new_replicas, 60*time.Minute).Err()
	}
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
