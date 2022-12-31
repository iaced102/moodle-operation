package worker

import (
	"context"
	adapter "moodle/adapter/mongo"
	sendmailclient "moodle/client/sendmail"

	"go.mongodb.org/mongo-driver/bson"
)



type SendmailWorker struct {
    adapter adapter.MongoAdapter
	sendmailclient sendmailclient.SendmailClient
}

func NewSendmailWorker(m adapter.MongoAdapter) *SendmailWorker{
    return &SendmailWorker{
		adapter: m,
		sendmailclient: *sendmailclient.NewSendmailClient(),
    }
}

type MailQueue struct {
	MoodleId 	  string `json:"moodle_id"`
	Email 		  string `json:"email"`
	LbStatus string `json:"lb_status"`
	K8sStatus string `json:"k8s_status"`
	MariadbStatus string `json:"mariadb_status"`
}

// get mail queue from mail queue collection
func (worker *SendmailWorker) GetMailQueue() ([]MailQueue, error) {
	ctx := context.Background()
	var mailQueue []MailQueue
	collection := worker.adapter.Collection("mail_queue")
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	if err = cursor.All(ctx, &mailQueue); err != nil {
		return nil, err
	}
	return mailQueue, nil
}


// check if all status is true
func (worker *SendmailWorker) CheckStatus(mailQueue MailQueue) bool {
	if mailQueue.LbStatus == "true" && mailQueue.K8sStatus == "true" && mailQueue.MariadbStatus == "true" {
		return true
	}
	return false
}


// send mail
func (worker *SendmailWorker) SendMail(mailQueue MailQueue) error {
	to := []string{mailQueue.Email}
	subject := "Moodle Deployed"
	body := "Moodle Deployed"
	err := worker.sendmailclient.SendMail(to, subject, body)
	if err != nil {
		return err
	}
	return nil
}
