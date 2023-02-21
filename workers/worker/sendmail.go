package worker

import (
	"log"
	"moodle/internal/core/domain"
	repo "moodle/internal/repository/moodle"
	mongo "moodle/pkg/mongodbiface"
	"moodle/pkg/client"
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

func (worker *SendmailWorker) SendMail(mailTracking domain.MailTracking) error {
	// get moodle info
	var moodle domain.Moodle
	moodle, err := worker.mongoRepo.Get(mailTracking.MoodleId)
	if err != nil {
		return err
	}
	packages := moodle.Packages

	email := moodle.Email
	webSiteName := moodle.WebSiteName
	ip := moodle.Ip
	packageName := packages.Name
	username := "manager"
	password := "Z}6a@7Dybf<l"
	err = worker.UpdateIsSent(mailTracking.MoodleId, true)
	if err != nil {
		return err
	}
	err = client.SendMail(email, webSiteName, ip, packageName, username, password)
	if err != nil {
		err = worker.UpdateIsSent(mailTracking.MoodleId, false)
		if err != nil {
			return err
		}
		return err
	}
	return nil
}

func (worker *SendmailWorker) UpdateIsSent(moodleId string, isSent bool) error {
	err := worker.mongoRepo.UpdateMailTracking(moodleId, isSent)
	if err != nil {
		log.Println(err)
	}
	return nil
}

func (worker *SendmailWorker) SendMailWorkerPool() {
	createQueue, err := worker.GetMailCreateQueue()
	if err != nil {
		return
	}
	for _, mail := range createQueue{
		go worker.SendMail(mail)
	}
}
