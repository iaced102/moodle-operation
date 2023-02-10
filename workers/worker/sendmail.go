package worker

import (
	"context"
	"log"
	mongoAdapter "moodle/adapter/mongo"
	sendmailclient "moodle/pkg/client"

	"go.mongodb.org/mongo-driver/bson"
)


type SendmailWorker struct {
    mongoAdapter mongoAdapter.MongoAdapter
	sendmailclient sendmailclient.SendmailClient
}

func NewSendmailWorker(m mongoAdapter.MongoAdapter) *SendmailWorker{
    return &SendmailWorker{
		mongoAdapter: m,
		sendmailclient: *sendmailclient.NewSendmailClient(),
    }
}

type MailTracking struct {
	MoodleId 	  string `json:"moodle_id"`
	Email 		  string `json:"email"`
	LbStatus string `json:"lb_status"`
	K8sStatus string `json:"k8s_status"`
	MariadbStatus string `json:"mariadb_status"`
}

// get mail from mailtracking collection where isSent is false
func (worker *SendmailWorker) GetMailQueue() ([]MailTracking, error) {
	ctx := context.Background()
	var mails []MailTracking
	collection := worker.mongoAdapter.Collection("mail_tracking")
	filter := bson.M{"issent": false}
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return mails, err
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var mail MailTracking
		err := cursor.Decode(&mail)
		if err != nil {
			return mails, err
		}
		mails = append(mails, mail)
	}
	return mails, nil
}

// send mail
func (worker *SendmailWorker) SendMail(mailTracking MailTracking) error {
	to := []string{mailTracking.Email}
	subject := "Moodle Deployed"
	body := "Moodle Deployed"
	err := worker.sendmailclient.SendMail(to, subject, body)
	if err != nil {
		return err
	}
	log.Println("Mail sent to: ", mailTracking.Email)
	worker.UpdateIsSent(mailTracking.MoodleId)
	return nil
}

func (worker *SendmailWorker) UpdateIsSent(moodleId string) {
	ctx := context.Background()
	collection := worker.mongoAdapter.Collection("mail_tracking")
	filter := bson.M{"moodleid": moodleId}
	update := bson.M{"$set": bson.M{"issent": true}}
	_, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Println(err)
	}
}

// worker pool to send mail concurrency
func (worker *SendmailWorker) SendMailWorkerPool() {
	mails, err := worker.GetMailQueue()
	if err != nil {
		return
	}
	for _, mail := range mails {
		go worker.SendMail(mail)
	}
}
