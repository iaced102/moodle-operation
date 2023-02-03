package handler

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"math"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"

	mariaadapter "moodle/adapter/mariadb"
	mongoadapter "moodle/adapter/mongo"
	k8sclient "moodle/client/k8s"
	"moodle/config"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"k8s.io/client-go/kubernetes"
)


type Handler struct {
	k8sclient k8sclient.K8sClient
	clientset *kubernetes.Clientset
	mariaclient *mariaadapter.MariaAdapter
    db mongoadapter.MongoAdapter
}

func New(d mongoadapter.MongoAdapter) *Handler {
	k8sclient := k8sclient.NewK8sClient()
	clientset := k8sclient.NewClientSet()
	mariaclient := mariaadapter.NewMariaAdapter(config.MARIAHOSTW, 3306, config.MARIAUSER, config.MARIAPASSWORD)
    return &Handler{
		k8sclient: *k8sclient,
		clientset: clientset,
		mariaclient: mariaclient,
		db: d,
    }
}

type MoodlePackages struct {
	Price int `json:"price"`
	Name string `json:"name"`
	Ccu int `json:"ccu"`
	AccountMax int `json:"account_max"`
	DocumentStorage int `json:"document_storage"`
	MoodleVersion string `json:"moodle_version"`
	BackupNum int `json:"backup_num"`
	CcuExtraMax int `json:"ccu_extra_max"`
	DocumentStorageExtraMax int `json:"document_storage_extra_max"`
}

var Package100 = MoodlePackages{
	Price: 8500000,
	Name : "100CCU",
	Ccu: 100,
	AccountMax: 4000,
	DocumentStorage: 50,
	MoodleVersion: "4.0.1",
	BackupNum: 4,
	CcuExtraMax: 400,
	DocumentStorageExtraMax: 2000,
}

var Package200 = MoodlePackages{
	Price: 12600000,
	Name : "200CCU",
	Ccu: 200,
	AccountMax: 8000,
	DocumentStorage: 100,
	MoodleVersion: "4.0.1",
	BackupNum: 4,
	CcuExtraMax: 800,
	DocumentStorageExtraMax: 4000,
}

var Package300 = MoodlePackages{
	Price: 16800000,
	Name : "300CCU",
	Ccu: 300,
	AccountMax: 12000,
	DocumentStorage: 150,
	MoodleVersion: "4.0.1",
	BackupNum: 4,
	CcuExtraMax: 1200,
	DocumentStorageExtraMax: 6000,
}

var Package400 = MoodlePackages{
	Price: 20900000,
	Name : "400CCU",
	Ccu: 400,
	AccountMax: 16000,
	DocumentStorage: 200,
	MoodleVersion: "4.0.1",
	BackupNum: 4,
	CcuExtraMax: 1600,
	DocumentStorageExtraMax: 8000,
}

var Package500 = MoodlePackages{
	Price: 0,
	Name : ">500CCU",
	Ccu: 0,
	AccountMax: 16000,
	DocumentStorage: 200,
	MoodleVersion: "4.0.1",
	BackupNum: 4,
	CcuExtraMax: 0,
	DocumentStorageExtraMax: 0,
}

type LBTracking struct {
	MoodleId string `json:"moodle_id"`
	LbName string `json:"lb_name"`
	VipAddress string `json:"vip_address"`
	LbStatus string `json:"lb_status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type DBTracking struct {
	MoodleId string `json:"moodle_id"`
	DbName string `json:"db_name"`
	SiteName string `json:"site_name"`
	SiteNameUpdate bool `json:"site_name_update"`
	DbStatus string `json:"lb_status"`
	FilePath string `json:"file_path"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type MoodleConfig struct {
	MoodleId string `json:"moodle_id"`
	DbName string `json:"db_name"`
	Ip string `json:"ip"`
	CreatedAt string `json:"created_at"`
	UpadatedAt string `json:"updated_at"`
}

type CreateMoodlePayload struct {
	Email string `json:"email"`
	WebSiteName string `json:"website_name"`
	PreInstalledCourse []int `json:"pre_installed_course"`
	PakcagesName string `json:"packages_name"`
	AutoScale bool `json:"autoscale"`
	DocumentsStorageExtra int `json:"documents_storage_extra"`
}

type MailTracking struct {
	MoodleId string `json:"moodle_id"`
	Email string `json:"email"`
	IsSent bool `json:"is_sent"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type DeleteMoodlePayload struct {
	MoodleId string `json:"moodle_id"`
}

type DeleteMoodleResponse struct {
	Id string `json:"id"`
	Status string `json:"status"`
}

type Moodle struct {
	Id  string `json:"id"`
	Email string `json:"email"`
	Ip string `json:"ip"`
	Name string `json:"name"`
	LbName string `json:"lb_name"`
	WebSiteName string `json:"website_name"`
	PreInstalledCourse []int `json:"pre_installed_course"`
	Packages MoodlePackages `json:"packages"`
	AutoScale bool `json:"autoscale"`
	DocumentsStorageExtra int `json:"documents_storage_extra"`
	Status string `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type MoodleErrorResponse struct {
	Error string `json:"error"`
}

// create moodle
func (h *Handler) CreateMoodle(w http.ResponseWriter, r *http.Request) {
	// get the payload
	var payload CreateMoodlePayload
	var moodle Moodle
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		moodleErrorResponse := MoodleErrorResponse{Error: "Error when decode payload"}
		json.NewEncoder(w).Encode(moodleErrorResponse)
		return
	}
	
	// check if website name is exist then response website name is exist
	exist := h.ValidateMoodleWebSiteName(payload.WebSiteName+".lms.bizflycloud.vn")
	if exist {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(MoodleErrorResponse{Error: "Website name is exist"})
		return
	}

	switch payload.PakcagesName {
	case "100CCU":
		moodle.Packages = Package100
	case "200CCU":
		moodle.Packages = Package200
	case "300CCU":
		moodle.Packages = Package300
	case "400CCU":
		moodle.Packages = Package400
	case "500CCU":
		moodle.Packages = Package500
	default:
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(MoodleErrorResponse{Error: "moodle-packages is not valid, please choose 100CCU, 200CCU, 300CCU, 400CCU or Liên Hệ"})
		return
	}
	moodle.Id = uuid.New().String()
	moodle.Email = payload.Email
	moodle.Name = payload.WebSiteName
	moodle.LbName = "kube_service" + "_" + config.CLUSTERID + "_" + moodle.Id + "_moodle-service"
	moodle.Ip = "Provisioning"
	moodle.WebSiteName = payload.WebSiteName + ".lms.bizflycloud.vn"
	moodle.PreInstalledCourse = payload.PreInstalledCourse
	moodle.AutoScale = payload.AutoScale
	moodle.DocumentsStorageExtra = payload.DocumentsStorageExtra
	moodle.Status = "Creating"
	moodle.CreatedAt = time.Now()
	moodle.UpdatedAt = time.Now()

	// create namespace before creating statefulset
	_, err = h.k8sclient.CreateNamespace(h.clientset, moodle.Id)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			moodleErrorResponse := MoodleErrorResponse{Error: err.Error()}
			json.NewEncoder(w).Encode(moodleErrorResponse)
			return
	}

	lbTracking := LBTracking{
		MoodleId: moodle.Id,
		LbName: "kube_service" + "_" + config.CLUSTERID + "_" + moodle.Id + "_moodle-service",
		VipAddress: "",
		LbStatus: "Creating",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// create mariadb for moodle
	moodleDbName := strings.ReplaceAll(moodle.Id, "-", "_")
	err = h.mariaclient.CreateDatabase(moodleDbName)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		moodleErrorResponse := MoodleErrorResponse{Error: err.Error()}
		json.NewEncoder(w).Encode(moodleErrorResponse)
		return
	}

	mariaTracking := DBTracking{
		MoodleId: moodle.Id,
		DbName: moodleDbName,
		SiteName: moodle.Name,
		SiteNameUpdate: false,
		DbStatus: "Creating",
		FilePath: "$HOME/gits/moodle-operator/docker/moodle_seded.sql",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}


	mailTracking := MailTracking{
		MoodleId: moodle.Id,
		Email: moodle.Email,
		IsSent: false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// create pvc before creating statefulset
	_, err = h.k8sclient.ApplyPVC(h.clientset, moodle.Id)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		moodleErrorResponse := MoodleErrorResponse{Error: err.Error()}
		json.NewEncoder(w).Encode(moodleErrorResponse)
		return
	}

	// apply the statefulset
	_, err = h.k8sclient.ApplyStatefulSet(h.clientset, moodle.Id, "default")
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		moodleErrorResponse := MoodleErrorResponse{Error: err.Error()}
		json.NewEncoder(w).Encode(moodleErrorResponse)
		return
	}

	// apply the service
	_, err = h.k8sclient.ApplyService(h.clientset, moodle.Id)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		moodleErrorResponse := MoodleErrorResponse{Error: err.Error()}
		json.NewEncoder(w).Encode(moodleErrorResponse)
		return
	}

	// write to db

	mailTrackingCollection := h.db.Collection("mail_tracking")
	_, err = mailTrackingCollection.InsertOne(context.Background(), mailTracking)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		moodleErrorResponse := MoodleErrorResponse{Error: err.Error()}
		json.NewEncoder(w).Encode(moodleErrorResponse)
		return
	}

	lbTrackingCollection := h.db.Collection("lb_tracking")
	_, err = lbTrackingCollection.InsertOne(context.Background(), lbTracking)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		moodleErrorResponse := MoodleErrorResponse{Error: err.Error()}
		json.NewEncoder(w).Encode(moodleErrorResponse)
		return
	}

	mariaTrackingCollection := h.db.Collection("maria_tracking")
	_, err = mariaTrackingCollection.InsertOne(context.Background(), mariaTracking)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		moodleErrorResponse := MoodleErrorResponse{Error: err.Error()}
		json.NewEncoder(w).Encode(moodleErrorResponse)
		return
	}

	moodleCollection := h.db.Collection("moodles")
	_, err = moodleCollection.InsertOne(context.Background(), moodle)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			moodleErrorResponse := MoodleErrorResponse{Error: err.Error()}
			json.NewEncoder(w).Encode(moodleErrorResponse)
			return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(moodle)

}

type ChangeMoodlePackagesPayload struct {
	MoodleId string `json:"moodle_id"`
	PackagesName string `json:"packages_name"`
}

// change moodle.packages by moodle.Id
func (h *Handler) ChangeMoodlePackages(w http.ResponseWriter, r *http.Request) {
	// get the payload
	var payload ChangeMoodlePackagesPayload
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		moodleErrorResponse := MoodleErrorResponse{Error: "Not implemented"}
		json.NewEncoder(w).Encode(moodleErrorResponse)
		return
	}

	// get the moodle
	moodle, err := h.GetMoodleById(payload.MoodleId)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		moodleErrorResponse := MoodleErrorResponse{Error: err.Error()}
		json.NewEncoder(w).Encode(moodleErrorResponse)
		return
	}
	switch payload.PackagesName {
	case "100CCU":
		moodle.Packages = Package100
	case "200CCU":
		moodle.Packages = Package200
	case "300CCU":
		moodle.Packages = Package300
	case "400CCU":
		moodle.Packages = Package400
	case "500CCU":
		moodle.Packages = Package500
	default:
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		moodleErrorResponse := MoodleErrorResponse{Error: "moodle-packages is not valid, please choose [100CCU, 200CCU, 300CCU, 400CCU, Liên Hệ]"}
		json.NewEncoder(w).Encode(moodleErrorResponse)
		return
	}

	// update moodle
	moodleCollection := h.db.Collection("moodles")
	_, err = moodleCollection.UpdateOne(context.Background(), bson.M{"id": payload.MoodleId}, bson.M{"$set": moodle})
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		moodleErrorResponse := MoodleErrorResponse{Error: "Not implemented"}
		json.NewEncoder(w).Encode(moodleErrorResponse)
		return

	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(moodle)
}

type ChangeMoodleAutoScalePayload struct {
	MoodleId string `json:"moodle_id"`
	AutoScale bool `json:"autoscale"`
}

// get moodle by moodle.Id
func (h *Handler) GetMoodleById(id string) (Moodle, error) {
	moodleCollection := h.db.Collection("moodles")
	var moodle Moodle
	err := moodleCollection.FindOne(context.Background(), bson.M{"id": id}).Decode(&moodle)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return moodle, errors.New(id + " is not found")
		}
		return moodle, err
	}
	return moodle, nil
}

// get moodle
func (h *Handler) GetMoodle(w http.ResponseWriter, r *http.Request) {
	// get the query
	query := r.URL.Query()
	id := query.Get("id")
	// get the moodle
	if id != "" {
		moodle, err := h.GetMoodleById(id)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusNotFound)
			moodleErrorResponse := MoodleErrorResponse{Error: err.Error()}
			json.NewEncoder(w).Encode(moodleErrorResponse)
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(moodle)
		return
	}

	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(MoodleErrorResponse{Error: "id is required"})
}

// change moodle.autoscale by moodle.Id
func (h *Handler) ChangeMoodleAutoScale(w http.ResponseWriter, r *http.Request) {
	// get the payload
	var payload ChangeMoodleAutoScalePayload
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// get the moodle by moodle.Id
	moodle, err := h.GetMoodleById(payload.MoodleId)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		moodleErrorResponse := MoodleErrorResponse{Error: err.Error()}
		json.NewEncoder(w).Encode(moodleErrorResponse)
		return
	}

	// change moodle.autoscale
	moodle.AutoScale = payload.AutoScale

	// update moodle
	moodleCollection := h.db.Collection("moodles")
	_, err = moodleCollection.UpdateOne(context.Background(), bson.M{"id": payload.MoodleId}, bson.M{"$set": moodle})
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		moodleErrorResponse := MoodleErrorResponse{Error: err.Error()}
		json.NewEncoder(w).Encode(moodleErrorResponse)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(moodle)
}

// change  moodle.documents_storage_extra by moodle.Id
type ChangeMoodleDocumentsStorageExtraPayload struct {
	MoodleId string `json:"moodle_id"`
	DocumentsStorageExtra int `json:"documents_storage_extra"`
}

func (h *Handler) ChangeMoodleDocumentsStorageExtra(w http.ResponseWriter, r *http.Request) {
	// get the payload
	var payload ChangeMoodleDocumentsStorageExtraPayload
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// get the moodle by moodle.Id
	moodle, err := h.GetMoodleById(payload.MoodleId)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		moodleErrorResponse := MoodleErrorResponse{Error: err.Error()}
		json.NewEncoder(w).Encode(moodleErrorResponse)
		return
	}

	if payload.DocumentsStorageExtra > moodle.Packages.DocumentStorageExtraMax || payload.DocumentsStorageExtra < 0 {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err)
		moodleErrorResponse := MoodleErrorResponse{Error: "documents_storage_extra is not valid, please choose in range [0, " + strconv.Itoa(moodle.Packages.DocumentStorageExtraMax) + "]"}
		json.NewEncoder(w).Encode(moodleErrorResponse)
		return
	}

	// change moodle.documents_storage_extra
	moodle.DocumentsStorageExtra = payload.DocumentsStorageExtra

	// update moodle
	moodleCollection := h.db.Collection("moodles")
	_, err = moodleCollection.UpdateOne(context.Background(), bson.M{"id": payload.MoodleId}, bson.M{"$set": moodle})
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		moodleErrorResponse := MoodleErrorResponse{Error: err.Error()}
		json.NewEncoder(w).Encode(moodleErrorResponse)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(moodle)
}

// change moodle.pre_installed_course by moodle.Id
type ChangeMoodlePreInstalledCoursePayload struct {
	MoodleId string `json:"moodle_id"`
	PreInstalledCourse []int `json:"pre_installed_course"`
}

func (h *Handler) ChangeMoodlePreInstalledCourse(w http.ResponseWriter, r *http.Request) {
	// get the payload
	var payload ChangeMoodlePreInstalledCoursePayload
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// validate the payload any course is not in range [1,8]
	for _, course := range payload.PreInstalledCourse {
		if course < 1 || course > 8 {
			w.WriteHeader(http.StatusBadRequest)
			log.Println(err)
			moodleErrorResponse := MoodleErrorResponse{Error: "pre_installed_course is not valid, please choose in range [1, 8]"}
			json.NewEncoder(w).Encode(moodleErrorResponse)
			return
		}
	}

	// get the moodle by moodle.Id
	moodle, err := h.GetMoodleById(payload.MoodleId)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		moodleErrorResponse := MoodleErrorResponse{Error: err.Error()}
		json.NewEncoder(w).Encode(moodleErrorResponse)
		return
	}

	// change moodle.pre_installed_course
	moodle.PreInstalledCourse = payload.PreInstalledCourse

	// update moodle
	moodleCollection := h.db.Collection("moodles")
	_, err = moodleCollection.UpdateOne(context.Background(), bson.M{"id": payload.MoodleId}, bson.M{"$set": moodle})
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		moodleErrorResponse := MoodleErrorResponse{Error: err.Error()}
		json.NewEncoder(w).Encode(moodleErrorResponse)
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(moodle)
}



// Delete moodle
func (h *Handler) DeleteMoodle(w http.ResponseWriter, r *http.Request) {
	// get the payload
	var payload DeleteMoodlePayload
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var resp DeleteMoodleResponse
	resp.Id = payload.MoodleId
	resp.Status = "Deleted"
	// delete the namespace
	err = h.k8sclient.DeleteNamespace(h.clientset, payload.MoodleId)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			moodleErrorResponse := MoodleErrorResponse{Error: err.Error()}
			json.NewEncoder(w).Encode(moodleErrorResponse)
		}

	mariaTrackingCollection := h.db.Collection("maria_tracking")
	// update dbstatus to mariaTrackingCollection
	_, err = mariaTrackingCollection.UpdateOne(context.Background(), bson.M{"moodleid": payload.MoodleId}, bson.M{"$set": bson.M{"dbstatus": "Deleting"}})
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		moodleErrorResponse := MoodleErrorResponse{Error: err.Error()}
		json.NewEncoder(w).Encode(moodleErrorResponse)
	}

	// delete moodle from db
	moodleCollection := h.db.Collection("moodles")
	_, err = moodleCollection.DeleteMany(context.Background(), bson.M{"id": payload.MoodleId})
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		moodleErrorResponse := MoodleErrorResponse{Error: err.Error()}
		json.NewEncoder(w).Encode(moodleErrorResponse)
	}
	LbTrackingCollection := h.db.Collection("lb_tracking")
	_, err = LbTrackingCollection.DeleteMany(context.Background(), bson.M{"moodleid": payload.MoodleId})
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		moodleErrorResponse := MoodleErrorResponse{Error: err.Error()}
		json.NewEncoder(w).Encode(moodleErrorResponse)
	}

	// this for clean worker

	// err = h.mariaclient.DropDatabase(strings.ReplaceAll(payload.MoodleId, "-", "_"))
	// if err != nil {
	// 	log.Println(err)
	// 	w.WriteHeader(http.StatusBadRequest)
	// 	moodleErrorResponse := MoodleErrorResponse{Error: err.Error()}
	// 	json.NewEncoder(w).Encode(moodleErrorResponse)
	// }
	// MariaTrackingCollection := h.db.Collection("maria_tracking")
	// _, err = MariaTrackingCollection.DeleteMany(context.Background(), bson.M{"moodleid": payload.MoodleId})
	// if err != nil {
	// 	log.Println(err)
	// 	w.WriteHeader(http.StatusBadRequest)
	// 	moodleErrorResponse := MoodleErrorResponse{Error: err.Error()}
	// 	json.NewEncoder(w).Encode(moodleErrorResponse)
	// }

	w.WriteHeader(http.StatusNoContent)
	json.NewEncoder(w).Encode(resp)

}

type ListByIdPayload struct {
	Id string `json:"userId"`
}

// list all moodles by email
func (h *Handler) ListMoodle(w http.ResponseWriter, r *http.Request) {
	// result := h.mariaclient.Select("moodle", "SELECT * FROM mdl_course")
	// log.Println(result)

	// err := h.mariaclient.CreateDatabase("moodle002")
	// if err != nil {
	// 	log.Println(err)
	// }

	// err := h.mariaclient.RestoreDatabase("moodle002", "/home/duy/gits/moodle-operator/docker/moodle21122022.sql")
	// if err != nil {
	// 	log.Println(err)
	// }
	
	// get id from params
	v := r.URL.Query()
	email := v.Get("email")
	search := v.Get("search")

	var moodles []Moodle
	moodleCollection := h.db.Collection("moodles")
	// search if have search params else list moodle
	if search != "" {
		// search moodle by name
		cursor, err := moodleCollection.Find(context.Background(), bson.M{"name": bson.M{"$regex": search, "$options": "i"}})
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			moodleErrorResponse := MoodleErrorResponse{Error: "Email is NotFound"}
			json.NewEncoder(w).Encode(moodleErrorResponse)
		}
		defer cursor.Close(context.Background())
		for cursor.Next(context.Background()) {
			var moodle Moodle
			err := cursor.Decode(&moodle)
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusBadRequest)
				moodleErrorResponse := MoodleErrorResponse{Error: "Email is NotFound"}
				json.NewEncoder(w).Encode(moodleErrorResponse)
			}
			moodles = append(moodles, moodle)
		}
	} else {
		// list all moodle
		cursor, err := moodleCollection.Find(context.Background(), bson.M{"email": email})
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			moodleErrorResponse := MoodleErrorResponse{Error: err.Error()}
			json.NewEncoder(w).Encode(moodleErrorResponse)
		}
		defer cursor.Close(context.Background())
		for cursor.Next(context.Background()) {
			var moodle Moodle
			err := cursor.Decode(&moodle)
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusBadRequest)
				moodleErrorResponse := MoodleErrorResponse{Error: err.Error()}
				json.NewEncoder(w).Encode(moodleErrorResponse)
			}
			moodles = append(moodles, moodle)
		}
	}

	page, err1 := strconv.Atoi(v.Get("page"))
	limit, err2 := strconv.Atoi(v.Get("limit"))
	if err1 != nil || err2 != nil || page < 1 || limit < 1 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(MoodleErrorResponse{Error: "page or limit is not valid"})
		return
	}
	// resp := h.PaginationMoodle(moodles, page, limit)
	resp := PaginateList(moodles, page, limit)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

// return clientset
func (h *Handler) GetClientset() *kubernetes.Clientset {
	return h.clientset
}

// validate if moodle.WebSiteName is existed
func (h *Handler) ValidateMoodleWebSiteName(webSiteName string) bool {
	moodleCollection := h.db.Collection("moodles")
	var moodle Moodle
	err := moodleCollection.FindOne(context.Background(), bson.M{"websitename": webSiteName}).Decode(&moodle)
	if err != nil {
		return false
	}
	return true
}


type Course struct {
	Id int `json:"id"`
	Name string `json:"name"`
	Content string `json:"content"`
	Thumb_Url string `json:"thumb_url"`
	Highlight []string `json:"highlight"`
	Routine []string `json:"routine"`
	Course_Preview []string `json:"course_preview"`
}

type  Meta struct {
	Total int `json:"total"`
	Pages int `json:"pages"`
	Page int `json:"page"`
	Limit int `json:"limit"`
}

type CourseErrorResponse struct {
	Error string `json:"error"`
}

// list course from course collection with pagination
func (h *Handler) ListCourse(w http.ResponseWriter, r *http.Request) {
	v := r.URL.Query()
	var courses []Course
	courseCollection := h.db.Collection("courses")
	cursor, err := courseCollection.Find(context.Background(), bson.M{})
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusNotFound)
		courseErrorResponse := CourseErrorResponse{Error: "Not Found"}
		json.NewEncoder(w).Encode(courseErrorResponse)
		return
	}
	for cursor.Next(context.Background()) {
		var course Course
		cursor.Decode(&course)
		courses = append(courses, course)
	}

	page, err1 := strconv.Atoi(v.Get("page"))
	limit, err2 := strconv.Atoi(v.Get("limit"))

	if err1 != nil || err2 != nil || page < 1 || limit < 1 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(MoodleErrorResponse{Error: "page or limit is not valid"})
		return
	}

	resp := PaginateList(courses, page, limit)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

func (h *Handler) ListPackage(w http.ResponseWriter, r *http.Request) {
	var packages []MoodlePackages
	packages = append(packages, Package100, Package200, Package300, Package400, Package500)
	v := r.URL.Query()

	page, err1 := strconv.Atoi(v.Get("page"))
	limit, err2 := strconv.Atoi(v.Get("limit"))

	if err1 != nil || err2 != nil || page < 1 || limit < 1 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(MoodleErrorResponse{Error: "page or limit is not valid"})
		return
	}
	resp := PaginateList(packages, page, limit)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}


// pagiante
type Pagination interface {
	GetPage() int
	GetLimit() int
	GetTotal() int
	GetPages() int
	GetData() interface{}
}

type pagination struct {
	Page    int         `json:"page"`
	Limit   int         `json:"limit"`
	Total   int         `json:"total"`
	Pages   int         `json:"pages"`
	Data    interface{} `json:"data"`
}

func (p *pagination) GetPage() int {
	return p.Page
}

func (p *pagination) GetLimit() int {
	return p.Limit
}

func (p *pagination) GetTotal() int {
	return p.Total
}

func (p *pagination) GetPages() int {
	return p.Pages
}

func (p *pagination) GetData() interface{} {
	return p.Data
}

func PaginateList(list interface{}, page int, limit int) Pagination {
	total := reflect.ValueOf(list).Len()
	if total == 0 {
		return &pagination{
			Page:    1,
			Limit:   0,
			Total:   0,
			Pages:   1,
			Data:    []interface{}{},
		}
	}

	pages := int(math.Ceil(float64(total) / float64(limit)))
	if page > pages {
		page = pages
	}
	start := (page - 1) * limit
	end := start + limit
	if end > total {
		end = total
	}
	data := reflect.ValueOf(list).Slice(start, end).Interface()

	return &pagination{
		Page:    page,
		Limit:   limit,
		Total:   total,
		Pages:   pages,
		Data:    data,
	}
}
