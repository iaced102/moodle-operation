package worker

import (
	"context"
	"log"
	adapter "moodle/adapter/mongo"
	k8sclient "moodle/client/k8s"
	"moodle/config"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"k8s.io/client-go/kubernetes"
)

func checkMariaInstanceStatus(name string) (string, error) {
	// mongoclient
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.NewClient(options.Client().ApplyURI(config.MONGOURI))
	if err != nil {
		log.Fatal(err)
	}
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)
	// check instance status by name
	collection := client.Database("moodle").Collection("maria_instances")
	filter := bson.D{{Key: "name", Value: name}}
	var result bson.M
	err = collection.FindOne(context.Background(), filter).Decode(&result)
	if err != nil {
		log.Fatal(err)
	}
	return result["status"].(string), nil
}


// k8s worker
type K8sWorker struct {
    adapter adapter.MongoAdapter
	// clientset *kubernetes.Clientset
	k8sclient k8sclient.K8sClient
}

func NewK8sWorker(m adapter.MongoAdapter) *K8sWorker{
	k8sclient := k8sclient.NewK8sClient()
	// clientset := k8sclient.NewClientSet()
    return &K8sWorker{
		adapter: m,
		k8sclient: *k8sclient,
		// clientset: clientset,
    }
}



// list maria id from maria_queue
func (k *K8sWorker) ListMariaQueue() ([]string, error) {
	ctx := context.Background()
	// collection := mongoClient.Database("moodle").Collection("maria_queue")
	collection := k.adapter.Collection("maria_queue")
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
	collection := k.adapter.Collection("maria_instances")
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
				log.Fatal(err)
			}
			if result["status"].(string) == "ACTIVE" {
				results = append(results, result["id"].(string))
			}
		}
	}

	return results, nil
}

// get theme from moodles collection by moodle id
func (k *K8sWorker) GetTheme(moodleId string) (string, error) {
	collection := k.adapter.Collection("moodles")
	filter := bson.M{"name": moodleId}
	var result bson.M
	err := collection.FindOne(context.Background(), filter).Decode(&result)
	if err != nil {
		log.Fatal(err)
	}
	// if result is empty, return default theme
	if result == nil {
		return "", nil
	}
	return result["theme"].(string), nil
}


//  before apply statefulset, prepare statefulset template


// apply statefulset from active maria id
func (k *K8sWorker) ApplyStatefulSet(clientset *kubernetes.Clientset, mariaId string, theme string) error {
	// apply the statefulset
	_, err := k.k8sclient.ApplyStatefulSet(clientset, mariaId, theme)
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
		if theme == "" {
			log.Fatal("theme is empty")
		}
		if err != nil {
			log.Fatal(err)
		}

		err = k.ApplyStatefulSet(clientset, activeMariaId, theme)
		if err != nil {
			log.Fatal(err)
		}
	}
	return nil
}
