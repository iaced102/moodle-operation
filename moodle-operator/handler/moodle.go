package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	mongoadapter "moodle/adapter/mongo"
	k8sclient "moodle/client/k8s"
	mariaclient "moodle/client/maria"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"k8s.io/client-go/kubernetes"
)




type Handler struct {
	k8sclient k8sclient.K8sClient
	clientset *kubernetes.Clientset
	mariaclient mariaclient.MariaClient
    db mongoadapter.MongoAdapter
}

func New(d mongoadapter.MongoAdapter) *Handler {
	k8sclient := k8sclient.NewK8sClient()
	clientset := k8sclient.NewClientSet()
	mariaclient := mariaclient.NewMariaClient()
    return &Handler{
		k8sclient: *k8sclient,
		clientset: clientset,
		mariaclient: *mariaclient,
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

type MoodleQueue struct {
	Id string `json:"id"`
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
	uuid := uuid.New()

	name := payload.Name+"-"+uuid.String()
	moodle.Name = name[:31]
	moodle.Id = uuid.String()
	moodle.Ccu = payload.Ccu


	// create namespace before creating statefulset
	_, err = h.k8sclient.CreateNamespace(h.clientset, uuid.String()[:31])
		if err != nil {
		w.Write([]byte(err.Error()))
	}

	// write moodle id to maria_queue collection
	moodleQueue := MoodleQueue{
		Id: moodle.Id[:31],
	}
	moodleQueueCollection := h.db.Collection("maria_queue")
	_, err = moodleQueueCollection.InsertOne(context.Background(), moodleQueue)
	if err != nil {
		w.Write([]byte(err.Error()))
	}


	// create db before creating statefulset
	backupID := "0fb8221a-a31c-4dce-bf87-18bf918298f0"
	err = h.mariaclient.CreateInstanceFromBackup(moodle.Id[:31], backupID)
	if err != nil {
		w.Write([]byte(err.Error()))
	}

	// apply the service
	_, err = h.k8sclient.ApplyService(h.clientset, moodle.Id[:31])
		if err != nil {
		w.Write([]byte(err.Error()))
	}

	// create pvc before creating statefulset
	_, err = h.k8sclient.ApplyPVC(h.clientset, moodle.Id[:31])
		if err != nil {
		w.Write([]byte(err.Error()))
	}

	// apply the statefulset
	// statefulset, err := h.k8sclient.ApplyStatefulSet(h.clientset, moodle.Id[:31], payload.Theme)

	// replicas := statefulset.Spec.Replicas
	// moodle.Replicas = *replicas
	moodle.Replicas = 2
	moodle.Status = "Provisioning"
	moodle.Cpu = "2"
	moodle.Memory = "2Gi"
	moodle.Userid = payload.UserId


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
