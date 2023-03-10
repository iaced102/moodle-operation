package worker

import (
	"moodle/internal/core/domain"
	repo "moodle/internal/repository/moodle"
	nfs "moodle/pkg/client"
	nfs4 "moodle/pkg/client/nfs4"
	"moodle/pkg/mongodbiface"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
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

    // Create a client to connect to InfluxDB
    client := influxdb2.NewClientWithOptions("http://123.30.234.141:8086", "lmspoller:lmspoller", influxdb2.DefaultOptions().SetBatchSize(100))

    // Create a write API for the "mydb" database
    writeAPI := client.WriteAPI("lms", "lms")

    // Define a data point to write
    p := influxdb2.NewPointWithMeasurement("lmsnfs").
		AddField("moodleid", moodleid).
		AddField("current_storage", dirinfo.Size).
		AddField("current_default", moodle.Packages.DocumentStorage * 1000000).
		AddField("max_storage", moodle.Packages.DocumentStorageExtraMax * 1000000).
		AddTag("moodleid_tag", moodleid).
        SetTime(time.Now())

    // Write the data point to InfluxDB
    writeAPI.WritePoint(p)

    // Close the client when you're done
    client.Close()
	return nil
}

