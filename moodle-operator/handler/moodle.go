package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"time"

	mariaadapter "moodle/adapter/mariadb"
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

type MoodlePackages struct {
	Name string `json:"name"`
	Ccu int `json:"ccu"`
	AccountMax int `json:"account_max"`
	DocumentStorage int `json:"document_storage"`
	MoodleVersion string `json:"moodle_version"`
	BackupNum int `json:"backup_num"`
	CcuExtraMax int `json:"ccu_extra_max"`
	DocumentStorageExtraMax int `json:"document_storage_extra_max"`
}

var SmallPackages = MoodlePackages{
	Name : "100CCU",
	Ccu: 100,
	AccountMax: 4000,
	DocumentStorage: 50,
	MoodleVersion: "4.0.1",
	BackupNum: 4,
	CcuExtraMax: 400,
	DocumentStorageExtraMax: 2000,
}

var MediumPackages = MoodlePackages{
	Name : "200CCU",
	Ccu: 200,
	AccountMax: 8000,
	DocumentStorage: 100,
	MoodleVersion: "4.0.1",
	BackupNum: 4,
	CcuExtraMax: 800,
	DocumentStorageExtraMax: 4000,
}

var LargePackages = MoodlePackages{
	Name : "300CCU",
	Ccu: 300,
	AccountMax: 12000,
	DocumentStorage: 150,
	MoodleVersion: "4.0.1",
	BackupNum: 4,
	CcuExtraMax: 1200,
	DocumentStorageExtraMax: 6000,
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
	DbStatus string `json:"lb_status"`
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

type DeleteMoodlePayload struct {
	Id string `json:"id"`
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

type MoodleQueue struct {
	Id string `json:"id"`
}


// create moodle
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
	
	// check if email not exist then reponse user not found
	userCollection := h.db.Collection("users")
	var user User
	err = userCollection.FindOne(context.Background(), bson.M{"email": payload.Email}).Decode(&user)
	if err != nil {
		w.Write([]byte("user not found"))
		return
	}
	
	// check if website name is exist then response website name is exist
	exist := h.ValidateMoodleWebSiteName(payload.WebSiteName+".lms.bizflycloud.vn")
	if exist {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("website_name is exist"))
		return
	}

	switch payload.PakcagesName {
	case "100CCU":
		moodle.Packages = SmallPackages
	case "200CCU":
		moodle.Packages = MediumPackages
	case "300CCU":
		moodle.Packages = LargePackages
	default:
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("moodle-packages is not valid, please choose 100CCU, 200CCU, or 300CCU"))
		return
	}
	moodle.Id = uuid.New().String()
	moodle.Email = payload.Email
	moodle.Name = payload.WebSiteName
	moodle.LbName = "kube_service" + "_6o0cn9lv42livqek_" + moodle.Id + "_moodle-service"
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
		w.Write([]byte(err.Error()))
	}

	// write moodle id to maria_queue collection
	// moodleQueue := MoodleQueue{
	// 	Id: moodle.Id,
	// }
	// moodleQueueCollection := h.db.Collection("maria_queue")
	// _, err = moodleQueueCollection.InsertOne(context.Background(), moodleQueue)
	// if err != nil {
	// 	w.Write([]byte(err.Error()))
	// }

	// apply the service
	_, err = h.k8sclient.ApplyService(h.clientset, moodle.Id)
	if err != nil {
		w.Write([]byte(err.Error()))
	}
	lbTracking := LBTracking{
		MoodleId: moodle.Id,
		LbName: "kube_service" + "_6o0cn9lv42livqek_" + moodle.Id + "_moodle-service",
		VipAddress: "",
		LbStatus: "Creating",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	lbTrackingCollection := h.db.Collection("lb_tracking")
	_, err = lbTrackingCollection.InsertOne(context.Background(), lbTracking)
	if err != nil {
		w.Write([]byte(err.Error()))
	}

	mariaTracking := DBTracking{
		MoodleId: moodle.Id,
		DbName: "moodle_" + moodle.Id,
		DbStatus: "Creating",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	mariaTrackingCollection := h.db.Collection("maria_tracking")
	_, err = mariaTrackingCollection.InsertOne(context.Background(), mariaTracking)
	if err != nil {
		w.Write([]byte(err.Error()))
	}

	// create pvc before creating statefulset
	_, err = h.k8sclient.ApplyPVC(h.clientset, moodle.Id)
	if err != nil {
		w.Write([]byte(err.Error()))
	}

	// apply the statefulset
	_, err = h.k8sclient.ApplyStatefulSet(h.clientset, moodle.Id, "default")
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

type ChangeMoodlePackagesPayload struct {
	Id string `json:"id"`
	PackagesName string `json:"packages_name"`
}

// change moodle.packages by moodle.Id
func (h *Handler) ChangeMoodlePackages(w http.ResponseWriter, r *http.Request) {
	// get the payload
	var payload ChangeMoodlePackagesPayload
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// get the moodle
	moodle, err := h.GetMoodleById(payload.Id)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	switch payload.PackagesName {
	case "100CCU":
		moodle.Packages = SmallPackages
	case "200CCU":
		moodle.Packages = MediumPackages
	case "300CCU":
		moodle.Packages = LargePackages
	default:
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("moodle-packages is not valid, please choose [100CCU, 200CCU, or 300CCU]"))
		return
	}

	// update moodle
	moodleCollection := h.db.Collection("moodles")
	_, err = moodleCollection.UpdateOne(context.Background(), bson.M{"id": payload.Id}, bson.M{"$set": moodle})
	if err != nil {
		w.Write([]byte(err.Error()))
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(moodle)
}

type ChangeMoodleAutoScalePayload struct {
	Id string `json:"id"`
	AutoScale bool `json:"autoscale"`
}

// get moodle by moodle.Id
func (h *Handler) GetMoodleById(id string) (Moodle, error) {
	moodleCollection := h.db.Collection("moodles")
	var moodle Moodle
	err := moodleCollection.FindOne(context.Background(), bson.M{"id": id}).Decode(&moodle)
	if err != nil {
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
			w.Write([]byte(err.Error()))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(moodle)
		return
	}
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte("id is required"))
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
	moodle, err := h.GetMoodleById(payload.Id)
	if err != nil {
		w.Write([]byte("moodle not found"))
		return
	}

	// change moodle.autoscale
	moodle.AutoScale = payload.AutoScale

	// update moodle
	moodleCollection := h.db.Collection("moodles")
	_, err = moodleCollection.UpdateOne(context.Background(), bson.M{"id": payload.Id}, bson.M{"$set": moodle})
	if err != nil {
		w.Write([]byte(err.Error()))
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(moodle)
}

// change  moodle.documents_storage_extra by moodle.Id
type ChangeMoodleDocumentsStorageExtraPayload struct {
	Id string `json:"id"`
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
	moodle, err := h.GetMoodleById(payload.Id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("moodle not found"))
		return
	}

	if payload.DocumentsStorageExtra > moodle.Packages.DocumentStorageExtraMax || payload.DocumentsStorageExtra < 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("documents_storage_extra is not valid, please choose in range [0, " + strconv.Itoa(moodle.Packages.DocumentStorageExtraMax) + "]"))
		return
	}

	// change moodle.documents_storage_extra
	moodle.DocumentsStorageExtra = payload.DocumentsStorageExtra

	// update moodle
	moodleCollection := h.db.Collection("moodles")
	_, err = moodleCollection.UpdateOne(context.Background(), bson.M{"id": payload.Id}, bson.M{"$set": moodle})
	if err != nil {
		w.Write([]byte(err.Error()))
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(moodle)
}

// change moodle.pre_installed_course by moodle.Id
type ChangeMoodlePreInstalledCoursePayload struct {
	Id string `json:"id"`
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
			w.Write([]byte("pre_installed_course is not valid, please choose in range [1, 8]"))
			return
		}
	}

	// get the moodle by moodle.Id
	moodle, err := h.GetMoodleById(payload.Id)
	if err != nil {
		w.Write([]byte("moodle not found"))
		return
	}

	// change moodle.pre_installed_course
	moodle.PreInstalledCourse = payload.PreInstalledCourse

	// update moodle
	moodleCollection := h.db.Collection("moodles")
	_, err = moodleCollection.UpdateOne(context.Background(), bson.M{"id": payload.Id}, bson.M{"$set": moodle})
	if err != nil {
		w.Write([]byte(err.Error()))
	}

	w.Header().Set("Content-Type", "application/json")
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
	LbTrackingCollection := h.db.Collection("lb_tracking")
	_, err = LbTrackingCollection.DeleteMany(context.Background(), bson.M{"moodleid": payload.Id})
	if err != nil {
		w.Write([]byte(err.Error()))
	}
	MariaTrackingCollection := h.db.Collection("maria_tracking")
	_, err = MariaTrackingCollection.DeleteMany(context.Background(), bson.M{"moodleid": payload.Id})
	if err != nil {
		w.Write([]byte(err.Error()))
	}
	w.WriteHeader(http.StatusNoContent)
	json.NewEncoder(w).Encode(resp)

}

type ListByIdPayload struct {
	Id string `json:"userId"`
}

// list all moodles by email
func (h *Handler) ListMoodle(w http.ResponseWriter, r *http.Request) {
	// result := h.SelectDataFromMariaDB("moodle")
	// fmt.Println(result)
	
	// get id from params
	v := r.URL.Query()
	email := v.Get("email")
	var moodles []Moodle
	moodleCollection := h.db.Collection("moodles")
	cursor, err := moodleCollection.Find(context.Background(), bson.M{"email": email})
	if err != nil {
		w.Write([]byte(err.Error()))
	}
	for cursor.Next(context.Background()) {
		var moodle Moodle
		cursor.Decode(&moodle)
		moodles = append(moodles, moodle)
	}

	type ListMoodleResponse struct {
		Total int `json:"total"`
		Pages int `json:"pages"`
		Page int `json:"page"`
		Limit int `json:"limit"`
		Moodles []Moodle `json:"moodles"`
	}

	if len(moodles) == 0 {
		moodles = []Moodle{} 
		resp := ListMoodleResponse{
			Total: 0,
			Pages: 0,
			Page: 0,
			Limit: 0,
			Moodles: moodles,
		}
		json.NewEncoder(w).Encode(resp)
		return
	}
	// pagination
	var page int
	var limit int
	var err1 error
	var err2 error

	if v.Get("page") != "" {
		page, err1 = strconv.Atoi(v.Get("page"))
	}
	if v.Get("limit") != "" {
		limit, err2 = strconv.Atoi(v.Get("limit"))
	}
	if err1 != nil || err2 != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("page or limit is not valid"))
		return
	}
	if page < 1 || limit < 1 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("page or limit is not valid"))
		return
	}
	var total int = len(moodles)
	var pages int = int(math.Ceil(float64(total) / float64(limit)))
	if page > pages {
		page = pages
	}
	var start int = (page - 1) * limit
	var end int = start + limit
	if end > total {
		end = total
	}

	var moodlesPage []Moodle = moodles[start:end]
	var resp ListMoodleResponse
	resp.Moodles = moodlesPage
	resp.Total = total
	resp.Pages = pages
	resp.Page = page
	resp.Limit = limit
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)

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

// search engine for moodle
func (h *Handler) SearchMoodle(w http.ResponseWriter, r *http.Request) {
	// get id from params
	v := r.URL.Query()
	search := v.Get("search")
	var moodles []Moodle
	moodleCollection := h.db.Collection("moodles")
	// serach moodle.name by regex
	cursor, err := moodleCollection.Find(context.Background(), bson.M{"websitename": bson.M{"$regex": search}})

	// cursor, err := moodleCollection.Find(context.Background(), bson.M{"$text": bson.M{"$search": search}})
	if err != nil {
		w.Write([]byte(err.Error()))
	}
	for cursor.Next(context.Background()) {
		var moodle Moodle
		cursor.Decode(&moodle)
		moodles = append(moodles, moodle)
	}
	if len(moodles) == 0 {
		// write not found
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// pagination
	var page int
	var limit int
	var err1 error
	var err2 error

	if v.Get("page") != "" {
		page, err1 = strconv.Atoi(v.Get("page"))
	}
	if v.Get("limit") != "" {
		limit, err2 = strconv.Atoi(v.Get("limit"))
	}
	if err1 != nil || err2 != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("page or limit is not valid"))
		return
	}
	if page < 1 || limit < 1 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("page or limit is not valid"))
		return
	}
	var total int = len(moodles)
	var pages int = int(math.Ceil(float64(total) / float64(limit)))
	if page > pages {
		page = pages
	}
	var start int = (page - 1) * limit
	var end int = start + limit
	if end > total {
		end = total
	}

	type ListMoodleResponse struct {
		Total int `json:"total"`
		Pages int `json:"pages"`
		Page int `json:"page"`
		Limit int `json:"limit"`
		Moodles []Moodle `json:"moodles"`
	}

	var moodlesPage []Moodle = moodles[start:end]
	var resp ListMoodleResponse
	resp.Moodles = moodlesPage
	resp.Total = total
	resp.Pages = pages
	resp.Page = page
	resp.Limit = limit
	w.Header().Set("Content-Type", "application/json")
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


// select data from mariadb use maria worker
func (h *Handler) SelectDataFromMariaDB(dbname string) map[string]string {
	newMariaAdapter := mariaadapter.NewMariaAdapter("45.124.94.39", 3306, "duy", "4Yk7741J2JVWPTQkT9eKkcbAaTUs5XzTvIFL", dbname)
	result := newMariaAdapter.Select("SELECT * FROM mdl_course")
	return result
}
