package worker

import (
	"fmt"
	"moodle/config"
	"moodle/internal/core/domain"
	repo "moodle/internal/repository/moodle"
	nfs "moodle/pkg/client"
	nfs4 "moodle/pkg/client/nfs4"
	"moodle/pkg/mongodbiface"
	"net/http"
	"strings"
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
	moodleid := dirinfo.Path[17:53]
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
	fieldKey := "current_storage"
	fieldValue := dirinfo.Size
	fieldKey2 := "current_default"
	fieldValue2 := moodle.Packages.DocumentStorage * 1000000
	fieldKey3 := "max_storage"
	fieldValue3 := moodle.Packages.DocumentStorageExtraMax * 1000000
	// timestamp := time.Now().Unix()

	metrics := fmt.Sprintf("%s,%s=%s %s=%d,%s=%d,%s=%d", measurementName, tagKey, tagValue, fieldKey, fieldValue, fieldKey2, fieldValue2, fieldKey3, fieldValue3)
	resp, err := http.Post(url, "application/octet-stream", strings.NewReader(metrics))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	return nil
}

