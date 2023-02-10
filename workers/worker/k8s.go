package worker

import (
	"context"
	"log"
	mongo "moodle/pkg/mongodbiface"
	k8sclient "moodle/pkg/client"

	"go.mongodb.org/mongo-driver/bson"
	"k8s.io/client-go/kubernetes"
)

// k8s worker
type K8sWorker struct {
    mongo mongo.DB
	k8sclient k8sclient.K8sClient
}

func NewK8sWorker(m mongo.DB) *K8sWorker{
	k8sclient := k8sclient.NewK8sClient()
	// clientset := k8sclient.NewClientSet()
    return &K8sWorker{
		mongo: m,
		k8sclient: *k8sclient,
    }
}

// list maria id from maria_queue
func (k *K8sWorker) ListMariaQueue() ([]string, error) {
	ctx := context.Background()
	// collection := mongoClient.Database("moodle").Collection("maria_queue")
	collection := k.mongo.Collection("maria_queue")
	cursor, err := collection.Find(context.Background(), bson.D{})
	if err != nil {
		log.Fatal(err)
	}
	var results []string
	for cursor.Next(ctx) {
		var result bson.M
		err := cursor.Decode(&result)
		if err != nil {
			log.Fatal(err)
		}
		results = append(results, result["id"].(string))
	}
	return results, nil
}


// get list of ACTIVE maria instances from maria id list
func (k *K8sWorker) GetActiveMariaInstances(mariaIds []string) ([]string, error) {
	ctx := context.Background()
	collection := k.mongo.Collection("maria_instances")
	count, err := collection.CountDocuments(ctx, bson.D{})
	if err != nil {
		log.Fatal(err)
	}
	if count == 0 {
		return []string{}, nil
	}

	var results []string
	if len(mariaIds) > 0 {
		for _, mariaId := range mariaIds {
			filter := bson.M{"name": mariaId}
			var result bson.M
			err := collection.FindOne(context.Background(), filter).Decode(&result)
			if err != nil {
				continue
			}

			if result["status"] == "ACTIVE" {
				results = append(results, result["name"].(string))
			}
		}
	}

	return results, nil
}

// get theme from moodles collection by moodle id
func (k *K8sWorker) GetTheme(moodleId string) (string, error) {
	collection := k.mongo.Collection("moodles")
	filter := bson.M{"name": moodleId}
	var result bson.M
	err := collection.FindOne(context.Background(), filter).Decode(&result)
	if err != nil {
		if err.Error() == "mongo: no documents in result" {
			return "", nil
		} else {
			log.Fatal(err)
		}
	}
	return result["theme"].(string), nil
}


//  before apply statefulset, prepare statefulset template


// apply statefulset from active maria id
func (k *K8sWorker) ApplyStatefulSet(clientset *kubernetes.Clientset, mariaId string, theme string) error {
	// apply the statefulset
	err := k.k8sclient.ApplyStatefulSet(mariaId)
	if err != nil {
		log.Fatal(err)
	}

	return nil
}

// worker to apply statefulset
func (k *K8sWorker) ApplyStatefulSetWorker(clientset *kubernetes.Clientset) error {
	// get list of maria id from maria_queue
	mariaIds, err := k.ListMariaQueue()
	if err != nil {
		log.Fatal(err)
	}
	
	if len(mariaIds) == 0 {
		return nil
	}

	// get list of ACTIVE maria instances from maria id list
	activeMariaIds, err := k.GetActiveMariaInstances(mariaIds)
	if err != nil {
		log.Fatal(err)
	}
	if len(activeMariaIds) == 0 {
		return nil
	}

	// apply statefulset from active maria id
	for _, activeMariaId := range activeMariaIds {
		// get theme from moodles collection by moodle id
		theme, err := k.GetTheme(activeMariaId)
		if err != nil {
			log.Fatal(err)
		}
		if theme == "" {
			return nil
		}

		err = k.ApplyStatefulSet(clientset, activeMariaId, theme)
		if err != nil {
			log.Fatal(err)
		}
	}
	return nil
}
