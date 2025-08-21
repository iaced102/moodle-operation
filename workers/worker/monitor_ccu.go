package worker

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"moodle/config"
	"moodle/internal/core/domain"
	repo "moodle/internal/repository/moodle"
	"moodle/pkg/mongodbiface"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-redis/redis"
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
			// Decrypt password before connecting
			dbPass := cfg.DBPass
			if key, err := getSecretKeyFromConfig(); err == nil {
				if dec, err := decryptPasswordCFB(dbPass, key); err == nil {
					dbPass = dec
				} else {
					log.Printf("failed to decrypt DB password for %s: %v", moodleID, err)
				}
			} else {
				log.Printf("failed to get secret key: %v", err)
			}
			dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", cfg.DBUser, dbPass, cfg.DBHost, cfg.DBPort, dbname)
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

	w.CheckCCUAndAlert(moodle, data[0])
	if err != nil {
		return nil
		// panic(err)
	}
	defer resp.Body.Close()

	return nil
}

func (w *MonitorWorker) CheckCCUAndAlert(moodle domain.Moodle, currentCCU int) error {
	defaultCCU := moodle.Packages.Ccu
	maxCCU := moodle.Packages.CcuExtraMax

	// Check if we should send alert (only once per hour per site)
	shouldSendAlert := false
	alertType := ""

	if currentCCU > maxCCU {
		// ALERT: CCU > max CCU
		shouldSendAlert = true
		alertType = "alert"
	} else if currentCCU > defaultCCU {
		// WARN: CCU > default CCU
		shouldSendAlert = true
		alertType = "warn"
	}

	if shouldSendAlert {
		// Check if alert was already sent in the last hour
		alreadySent, err := w.mongoRepo.CheckCCUAlertSent(moodle.Id, alertType)
		if err != nil {
			log.Printf("Failed to check CCU alert status for %s: %v", moodle.Id, err)
			return err
		}

		if !alreadySent {
			// Send alert
			if err := w.sendTelegram(moodle, currentCCU, defaultCCU, maxCCU, alertType); err != nil {
				log.Printf("Failed to send CCU alert for %s: %v", moodle.Id, err)
				return err
			}

			// Log alert to MongoDB
			ccuAlert := domain.CCUAlertTracking{
				MoodleId:   moodle.Id,
				AlertType:  alertType,
				CurrentCCU: currentCCU,
				DefaultCCU: defaultCCU,
				MaxCCU:     maxCCU,
				CreatedAt:  time.Now(),
			}

			if err := w.mongoRepo.CreateCCUAlertTracking(ccuAlert); err != nil {
				log.Printf("Failed to save CCU alert tracking for %s: %v", moodle.Id, err)
			} else {
				log.Printf("CCU %s sent and saved for moodle %s: current=%d, default=%d, max=%d",
					alertType, moodle.Id, currentCCU, defaultCCU, maxCCU)
			}
		} else {
			log.Printf("CCU %s already sent for moodle %s in the last hour, skipping", alertType, moodle.Id)
		}
	}

	return nil
}

// sendTelegram sends message to Telegram using the same API as alarm_nfs.go
func (w *MonitorWorker) sendTelegram(m domain.Moodle, currentCCU, defaultCCU, maxCCU int, alertType string) error {
	apiUrl := "https://api.telegram.org/bot248926246:AAETwv7hzpk8zv6j9aRHDITYcnIRUGydS80/sendMessage"

	var text string
	if alertType == "alert" {
		text = fmt.Sprintf("[CCU ALERT] KH: %s, Email: %s, Site: %s, Current CCU: %d, Default CCU: %d, Max CCU: %d",
			m.Name, m.Email, m.WebSiteName, currentCCU, defaultCCU, maxCCU)
	} else {
		text = fmt.Sprintf("[CCU WARNING] KH: %s, Email: %s, Site: %s, Current CCU: %d, Default CCU: %d, Max CCU: %d",
			m.Name, m.Email, m.WebSiteName, currentCCU, defaultCCU, maxCCU)
	}

	message := map[string]interface{}{
		"chat_id":           -1002960560112,
		"message_thread_id": 2,
		"text":              text,
	}

	jsonValue, _ := json.Marshal(message)
	req, err := http.NewRequest("POST", apiUrl, bytes.NewBuffer(jsonValue))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return fmt.Errorf("telegram error: %s", resp.Status)
	}
	return nil
}

func getSecretKeyFromConfig() ([]byte, error) {
	keyBytes, err := hex.DecodeString(config.DB_ENC_KEY_HEX)
	if err != nil {
		return nil, err
	}
	if len(keyBytes) != 32 {
		return nil, fmt.Errorf("SecretKey must be 32 bytes (got %d)", len(keyBytes))
	}
	return keyBytes, nil
}

func decryptPasswordCFB(encryptedHex string, key []byte) (string, error) {
	ciphertext, err := hex.DecodeString(encryptedHex)
	if err != nil {
		return "", err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	if len(ciphertext) < aes.BlockSize {
		return "", fmt.Errorf("ciphertext too short")
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]
	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ciphertext, ciphertext)
	return string(ciphertext), nil
}
