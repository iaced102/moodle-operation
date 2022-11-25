package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	k8sclient "moodle/client/k8s"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"k8s.io/client-go/kubernetes"
)



type Database interface{
	 Collection(name string, opts ...*options.CollectionOptions) *mongo.Collection
	 CreateCollection(ctx context.Context, name string, opts ...*options.CreateCollectionOptions) error
}

type Handler struct {
	k8sclient k8sclient.K8sClient
	clientset *kubernetes.Clientset
    db Database
}

func New(d Database) *Handler {
	client := k8sclient.NewK8sClient()
	clientset := client.NewClientSet()
    return &Handler{
		k8sclient: *client,
		clientset: clientset,
        db: d,
    }
}

type CreateMoodlePayload struct {
	Name string `json:"name"`
	Ccu  int    `json:"ccu"`
	UserId string `json:"userId"`
	Theme string `json:"theme"`
}

type DeleteMoodlePayload struct {
	Id string `json:"id"`
}

type DeleteMoodleResponse struct {
	Id string `json:"id"`
	Status string `json:"status"`
}

type Moodle struct {
	Id  string `json:"id"`
	Userid string `json:"userId"`
	Name string `json:"name"`
	Ccu int    `json:"ccu"`
	Status string `json:"status"`
	Replicas int32 `json:"replicas"`
	Cpu string `json:"cpu"`
	Memory string `json:"memory"`
}


// Apply statefulset
func (h *Handler) CreateMoodle(w http.ResponseWriter, r *http.Request) {
	// get the payload
	var payload CreateMoodlePayload
	var moodle Moodle
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Println(payload)
	// create namespace before creating statefulset
	uuid := uuid.New()
	_, err = h.k8sclient.CreateNamespace(h.clientset, payload.Name+"-"+uuid.String())
		if err != nil {
		w.Write([]byte(err.Error()))
	}


	moodle.Id = payload.Name+"-"+uuid.String()
	moodle.Name = payload.Name
	moodle.Ccu = payload.Ccu

	// create db before creating statefulset


	// create pvc before creating statefulset
	_, err = h.k8sclient.ApplyPVC(h.clientset, moodle.Id)
		if err != nil {
		w.Write([]byte(err.Error()))
	}


	// apply the statefulset
	statefulset, err := h.k8sclient.ApplyStatefulSet(h.clientset, moodle.Id, payload.Theme)

	replicas := statefulset.Spec.Replicas
	moodle.Status = "Provisioning"
	moodle.Replicas = *replicas
	moodle.Cpu = "2"
	moodle.Memory = "2Gi"
	moodle.Userid = payload.UserId

	// apply the service
	_, err = h.k8sclient.ApplyService(h.clientset, moodle.Id)
		if err != nil {
		w.Write([]byte(err.Error()))
	}


	// write to db
	moodleCollection := h.db.Collection("moodles")
	_, err = moodleCollection.InsertOne(context.Background(), moodle)
		if err != nil {
		w.Write([]byte(err.Error()))
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(moodle)

}


// Delete statefulset
func (h *Handler) DeleteMoodle(w http.ResponseWriter, r *http.Request) {
	// get the payload
	var payload DeleteMoodlePayload
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var resp DeleteMoodleResponse
	resp.Id = payload.Id
	resp.Status = "Deleted"
	// delete the statefulset
	err = h.k8sclient.DeleteNamespace(h.clientset, payload.Id)
		if err != nil {
		w.Write([]byte(err.Error()))
	}

	// delete moodle from db
	moodleCollection := h.db.Collection("moodles")
	_, err = moodleCollection.DeleteMany(context.Background(), payload)
	if err != nil {
		w.Write([]byte(err.Error()))
	}
	w.WriteHeader(http.StatusNoContent)
	json.NewEncoder(w).Encode(resp)

}

type ListByIdPayload struct {
	Id string `json:"userId"`
}

// list all moodles by user id
func (h *Handler) ListMoodle(w http.ResponseWriter, r *http.Request) {
	// get id from params
	v := r.URL.Query()
	id := v.Get("userid")
	var moodles []Moodle
	moodleCollection := h.db.Collection("moodles")
	cursor, err := moodleCollection.Find(context.Background(), bson.M{"userid": id})
	if err != nil {
		w.Write([]byte(err.Error()))
	}
	for cursor.Next(context.Background()) {
		var moodle Moodle
		cursor.Decode(&moodle)
		moodles = append(moodles, moodle)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(moodles)
}


type ScaleMoodlePayload struct {
	Id string `json:"id"`
	Replicas int32 `json:"replicas"`
}
type ScaleMoodleResponse struct {
	Id string `json:"id"`
	Replicas int32 `json:"replicas"`
	Status string `json:"status"`
}

// scale statefulset
func (h *Handler) ScaleMoodle(w http.ResponseWriter, r *http.Request) {
	// get the payload
	var payload ScaleMoodlePayload
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// scale the statefulset
	err = h.k8sclient.ScaleStatefulSet(h.clientset, payload.Id, payload.Replicas)
		if err != nil {
		w.Write([]byte(err.Error()))
	}

	// update db
	moodleCollection := h.db.Collection("moodles")
	_, err = moodleCollection.UpdateOne(context.Background(), bson.M{"id": payload.Id}, bson.M{"$set": bson.M{"replicas": payload.Replicas}})
	if err != nil {
		w.Write([]byte(err.Error()))
	}
	var resp ScaleMoodleResponse
	resp.Id = payload.Id
	resp.Replicas = payload.Replicas
	resp.Status = "Scaled"
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

