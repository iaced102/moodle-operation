package worker

import (
	"context"
	"log"
	mongo "moodle/pkg/mongodbiface"
	lbclient "moodle/pkg/client"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)


type TrackingWorker struct {
    mongo mongo.DB
}

func NewTrackingWorker(m mongo.DB) *TrackingWorker{
    return &TrackingWorker{
		mongo: m,
    }
}


// update lbstatus and vipaddress from lb_instances collection to lb_tracking colection
func (w *TrackingWorker) UpdateLbStatus() error {
	lbInstancesCollection := w.mongo.Collection("lb_instances")
	lbTrackingCollection := w.mongo.Collection("lb_tracking")
	moodlesCollection := w.mongo.Collection("moodles")

	// get all lb instances
	lbInstances, err := lbInstancesCollection.Find(context.Background(), bson.M{})
	if err != nil {
		log.Println(err)
		return err
	}

	// update lbstatus and vipaddress from lb_instances collection to lb_tracking colection
	for lbInstances.Next(context.Background()) {
		var lbInstance lbclient.LoadBalancer
		err := lbInstances.Decode(&lbInstance)
		if err != nil {
			log.Println(err)
			return err
		}

		filter := bson.M{"lbname": lbInstance.Name}
		update1 := bson.M{"$set": bson.M{"lbstatus": lbInstance.ProvisioningStatus, "vipaddress": lbInstance.VipAddress, "updatedat": time.Now()}}
		update2 := bson.M{"$set": bson.M{"status": lbInstance.ProvisioningStatus, "ip": lbInstance.VipAddress, "updatedat": time.Now()}}
		_, err = lbTrackingCollection.UpdateOne(context.Background(), filter, update1)
		if err != nil {
			log.Println(err)
			return err
		}
		_, err = moodlesCollection.UpdateOne(context.Background(), filter, update2)
		if err != nil {
			log.Println(err)
			return err
		}
	}

	return nil
}
