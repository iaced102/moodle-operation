package worker

import (
	"fmt"
	"log"
	"moodle/config"
	"moodle/internal/core/domain"
	repo "moodle/internal/repository/moodle"
	nfs "moodle/pkg/client"
	nfs4 "moodle/pkg/client/nfs4"
	"moodle/pkg/mongodbiface"
	"net/http"
	"strings"
	"database/sql"
)

type MonitorNFS struct {
	mongoRepo *repo.MongoDB
	NFSClient	 *nfs4.NfsClient
}

func NewMonitorNFSWorker(mongo mongodbiface.DB) *MonitorNFS {
	mongoRepo := repo.NewMongoDB(mongo)
	nfsclient := nfs.NewNFSClient()

	return &MonitorNFS{
		mongoRepo: mongoRepo,
		NFSClient: nfsclient,
	}
}

type DirInfo struct {
	MoodleID string `json:"moodleid"`
	Name string `json:"name"`
	Size uint64 `json:"size"`
}


func (n *MonitorNFS) Run() {
	// dirinfo, err := n.GetNFSSize()
	// get all from nfs_tracking
	nfsTracking, err := n.mongoRepo.GetNFSTracking()
	if err != nil {
		// log.Println(err)
		return
	}
	fmt.Println(nfsTracking, nfsTracking.Data)
	for _, data := range nfsTracking.Data {
		go n.Write(data)
	}
}


func (n *MonitorNFS) GetNFSSize() ([]DirInfo, error) {
	files, err := n.NFSClient.GetFileList("/moodle")
	if err != nil {
		return nil, err
	}

	var dirInfo []DirInfo
	// check if type is dir
	for _, file := range files {
		if file.IsDir {
			dirInfo = append(dirInfo, DirInfo{
				MoodleID: file.Name[:36],
				Name: file.Name,
				Size: file.Size,
			})
		}
	}
	return dirInfo, nil
}

func (n *MonitorNFS) Write(dirinfo domain.NFSTracking) error {
	if len(dirinfo.Path) < 53 {
		log.Printf("Path too short, skipping: %s", dirinfo.Path)
		return nil
	}
	moodleid := dirinfo.Path[13:49]
	moodle, err := n.mongoRepo.Get(moodleid)
	// check if error is no documents in result then ignore
	if err != nil {
		moodle.Packages.DocumentStorage = 0
		moodle.Packages.DocumentStorageExtraMax = 0
	}
	url := config.VICTORIAMETRIC
	measurementName := "lmsnfs_storage"
	tagKey := "moodleid"
	tagValue := moodleid

	// Calculate DB size (bytes) using CCU monitor's Maria config logic
	dbname := strings.ReplaceAll(moodleid, "-", "_")
	var dbSizeBytes int64
	if mapping, err := n.mongoRepo.GetMoodleMariaMappingByMoodleId(moodleid); err == nil && mapping.MariaId != "" {
		if cfg, err := n.mongoRepo.GetMariaConfigById(mapping.MariaId); err == nil {
			// Decrypt password
			dbPass := cfg.DBPass
			if key, err := getSecretKeyFromConfig(); err == nil {
				if dec, err := decryptPasswordCFB(dbPass, key); err == nil {
					dbPass = dec
				} else {
					log.Printf("failed to decrypt DB password for %s: %v", moodleid, err)
				}
			} else {
				log.Printf("failed to get secret key: %v", err)
			}

			dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", cfg.DBUser, dbPass, cfg.DBHost, cfg.DBPort, dbname)
			if db, err := sql.Open("mysql", dsn); err == nil {
				defer db.Close()
				row := db.QueryRow("SELECT IFNULL(SUM(data_length+index_length),0) FROM information_schema.tables WHERE table_schema = ?", dbname)
				if err := row.Scan(&dbSizeBytes); err != nil {
					log.Printf("failed to scan db size for %s: %v", dbname, err)
				}
			}
		}
	}

	// Fallback to default Maria connection if mapping/config not found or size is zero
	if dbSizeBytes == 0 {
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", config.MARIAUSER, config.MARIAPASSWORD, config.MARIAHOSTW, config.MARIAPORT, dbname)
		if db, err := sql.Open("mysql", dsn); err == nil {
			defer db.Close()
			row := db.QueryRow("SELECT IFNULL(SUM(data_length+index_length),0) FROM information_schema.tables WHERE table_schema = ?", dbname)
			if err := row.Scan(&dbSizeBytes); err != nil {
				log.Printf("fallback: failed to scan db size for %s: %v", dbname, err)
			}
		}
	}

	// Combine NFS size and DB size into a single metric for lmsnfs_storage
	nfsBytes := int64(dirinfo.Size)
	totalBytes := nfsBytes + dbSizeBytes/1000
	fieldKey := "current_storage"
	fieldValue := totalBytes
	fieldKey2 := "current_default"
	fieldValue2 := moodle.Packages.DocumentStorage * 1000000
	fieldKey3 := "max_storage"
	fieldValue3 := moodle.Packages.DocumentStorageExtraMax * 1000000

	metrics := fmt.Sprintf("%s,%s=%s %s=%d,%s=%d,%s=%d", measurementName, tagKey, tagValue, fieldKey, fieldValue, fieldKey2, fieldValue2, fieldKey3, fieldValue3)
	resp, err := http.Post(url, "application/octet-stream", strings.NewReader(metrics))
	if err != nil {
		return nil
	}
	defer resp.Body.Close()

	return nil
}
