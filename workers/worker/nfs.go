package worker

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"moodle/internal/core/domain"
	repo "moodle/internal/repository/moodle"
	nfs "moodle/pkg/client"
	nfs4 "moodle/pkg/client/nfs4"
	helpers "moodle/pkg/helpers"
	mongo "moodle/pkg/mongodbiface"
	"os"
)


type NFSWorker struct {
	mongoRepo *repo.MongoDB
	NFSClient	 *nfs4.NfsClient
}

func NewNFSWorker(mongo mongo.DB) *NFSWorker {
	nfsclient := nfs.NewNFSClient()
	mongoRepo := repo.NewMongoDB(mongo)
	return &NFSWorker{
		mongoRepo: mongoRepo,
		NFSClient: nfsclient,
	}
}

// GetFileList get file list from nfs server
func (n *NFSWorker) ValidateFilePath(moodleid string) (string, error) {
	// moodleid := "3f952434-112c-45fc-b141-84a0ffdd3c85"
	files, err := n.NFSClient.GetFileList("/moodle")
	if err != nil {
		return "", err
	}
	filePath := ""
	for _, file := range files {
		if helpers.ContainsString(file.Name, moodleid) {
			filePath = fmt.Sprintf("/moodle/%s/filedir/2b/83/2b831e652219e350659e4a71af9c4b9c4c99411c", file.Name)
		}
	}
	log.Printf("File path: %s", filePath)
	return filePath, nil
}

// get filepath from logo_tracking collection
func (n *NFSWorker) UpdateLogo() error {
	// get all logo_tracking
	logoTrackings, err := n.mongoRepo.GetAllLogoTracking()
	if err != nil {
		return err
	}
	for logoTrackings.Next(context.Background()) {
		var logoTracking domain.LogoTracking
		err := logoTrackings.Decode(&logoTracking)
		if err != nil {
			return err
		}
		n.Update(logoTracking.MoodleId, logoTracking.FilePath)
	}
	return nil
}

// update
func (n *NFSWorker) Update(moodleid, src string) error {
	var logoTracking domain.LogoTracking
	// get file path
	filePath, err := n.ValidateFilePath(moodleid)
	if err != nil {
		return err
	}

	file, err := os.Open(src)
	if err != nil {
		// handle error
	}
	defer file.Close()

	reader := bufio.NewReader(file)

	// read file from src then write to filePath
	offset, err := getFileSize(src)
	if err != nil {
		return err
	}
	_, err = n.NFSClient.WriteFile(filePath, false, uint64(offset), reader)

	// update logo_tracking
	log.Println("Updating logo")
	logoTracking.FilePath = filePath
	logoTracking.Status = "updated"
	err = n.mongoRepo.UpdateLogoTracking(logoTracking)
	if err != nil {
		return err
	}
	return nil
}

// get filesize
func getFileSize(path string) (int64, error) {
    fileInfo, err := os.Stat(path)
    if err != nil {
        return 0, err
    }
    return fileInfo.Size(), nil
}
