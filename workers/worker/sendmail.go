package worker

import (
	"log"
	"moodle/internal/core/domain"
	repo "moodle/internal/repository/moodle"
	"moodle/pkg/client"
	mongo "moodle/pkg/mongodbiface"
	"strconv"
	"time"
)


type SendmailWorker struct {
    mongoRepo *repo.MongoDB
}

func NewSendmailWorker(m mongo.DB) *SendmailWorker{
	mongoRepo := repo.NewMongoDB(m)
    return &SendmailWorker{
		mongoRepo: mongoRepo,
    }
}

func (worker *SendmailWorker) GetMailCreateQueue() ([]domain.MailTracking, error) {
	mails, err := worker.mongoRepo.GetAllMailCreateTracking()
	if err != nil {
		return mails, err
	}
	return mails, nil
}

func (worker *SendmailWorker) GetMailDeleteQueue() ([]domain.MailTracking, error) {
	mails, err := worker.mongoRepo.GetAllMailDeleteTracking()
	if err != nil {
		return mails, err
	}
	return mails, nil
}

func (worker *SendmailWorker) SendMailCreate(mailTracking domain.MailTracking) error {
	// get moodle info
	var moodle domain.Moodle
	moodle, err := worker.mongoRepo.Get(mailTracking.MoodleId)
	if err != nil {
		return err
	}
	packages := moodle.Packages

	email := moodle.Email
	webSiteName := moodle.Name
	webSiteNameAddress := moodle.WebSiteName
	ip := moodle.Ip
	packageName := packages.Name + " / " + strconv.Itoa(packages.DocumentStorage) + "GB" + " / " + strconv.Itoa(packages.BackupNum) + " backups"
	username := "manager"
	password := "Z}6a@7Dybf<l"
	createdat := mailTracking.CreatedAt.Format(time.RFC3339)
	
	err = worker.UpdateIsSentCreate(mailTracking.MoodleId, true)
	if err != nil {
		return err
	}
	log.Println("Sending create mail to: ", email)
	err = client.SendMailCreate(email, webSiteName, webSiteNameAddress , ip, packageName, username, password, createdat)
	if err != nil {
		log.Println(err)
		err = worker.UpdateIsSentCreate(mailTracking.MoodleId, false)
		if err != nil {
			return err
		}
		return err
	}
	log.Println("Sent mail create to: ", email, " success")
	return nil
}

func (worker *SendmailWorker) SendMailDelete(mailTracking domain.MailTracking) error {
	var moodle domain.Moodle
	moodle, err := worker.mongoRepo.Get(mailTracking.MoodleId)
	if err != nil {
		return err
	}
	packages := moodle.Packages

	email := moodle.Email
	webSiteName := moodle.Name
	webSiteNameAddress := moodle.WebSiteName
	ip := moodle.Ip
	packageName := packages.Name + " / " + strconv.Itoa(packages.DocumentStorage) + "GB" + " / " + strconv.Itoa(packages.BackupNum) + " backups"
	createdat := mailTracking.CreatedAt.Format(time.RFC3339)
	
	err = worker.UpdateIsSentDelete(mailTracking.MoodleId, true)
	if err != nil {
		return err
	}
	log.Println("Sending delete mail to: ", email)
	err = client.SendMailDelete(email, webSiteName, webSiteNameAddress, ip, packageName,createdat)
	if err != nil {
		log.Println(err)
		err = worker.UpdateIsSentDelete(mailTracking.MoodleId,false)
		if err != nil {
			return err
		}
		return err
	}
	err = worker.mongoRepo.Delete(mailTracking.MoodleId)
	if err != nil {
		return err
	}

	log.Println("Sent mail delete to: ", email, " success")
	return nil
}

func (worker *SendmailWorker) UpdateIsSentCreate(moodleId string, isSent bool) error {
	err := worker.mongoRepo.UpdateCreateMailTracking(moodleId, isSent)
	if err != nil {
		log.Println(err)
	}
	return nil
}

func (worker *SendmailWorker) UpdateIsSentDelete(moodleId string, isSent bool) error {
	err := worker.mongoRepo.UpdateDeleteMailTracking(moodleId, isSent)
	if err != nil {
		log.Println(err)
	}
	return nil
}

func (worker *SendmailWorker) SendMailCreateWorkerPool() {
	createQueue, err := worker.GetMailCreateQueue()
	if err != nil {
		return
	}
	for _, mail := range createQueue{
		go worker.SendMailCreate(mail)
	}
}

func (worker *SendmailWorker) SendMailDeleteWorkerPool() {
	deleteQueue, err := worker.GetMailDeleteQueue()
	if err != nil {
		return
	}
	for _, mail := range deleteQueue{
		go worker.SendMailDelete(mail)
	}
}
