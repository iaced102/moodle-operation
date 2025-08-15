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
	// get all moodle IDs from moodles collection
	moodleIDs, err := w.mongoRepo.GetAllMoodleIDs()
	if err != nil {
		return err
	}
	// if len(moodleIDs) == 0 then return nil
	if len(moodleIDs) == 0 {
		return nil
	}
	// for each moodle ID, convert to dbname by replacing "-" with "_"
	for _, moodleID := range moodleIDs {
		dbname := strings.ReplaceAll(moodleID, "-", "_")
		go w.WriteUser(dbname)
	}
	return nil
}



func (w *MonitorWorker) WriteUser(dbname string) error {
	moodleID := strings.ReplaceAll(dbname, "_", "-")

	var data []int

	// Try mapping: moodle_id -> maria_config
	if mapping, err := w.mongoRepo.GetMoodleMariaMappingByMoodleId(moodleID); err == nil && mapping.MariaId != "" {
		if cfg, err := w.mongoRepo.GetMariaConfigById(mapping.MariaId); err == nil {
			dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", cfg.DBUser, cfg.DBPass, cfg.DBHost, cfg.DBPort, dbname)
			if db, err := sql.Open("mysql", dsn); err == nil {
				defer db.Close()
				tmpRepo := repo.NewMariaDB(db)
				data, _ = tmpRepo.CountUserOnline(dbname)
			}
		}
	}

	// Fallback to default Maria connection
	if len(data) == 0 {
		if d, err := w.mariaRepo.CountUserOnline(dbname); err == nil {
			data = d
		}
	}

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
