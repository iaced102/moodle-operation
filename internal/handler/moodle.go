package handler

import (
	"moodle/internal/core/domain"
	"moodle/internal/core/port"
	"net/http"

	"moodle/pkg/helpers"

	"github.com/gin-gonic/gin"
)


type MoodleHandler struct {
	MoodleService port.MoodleService
}

func NewMoodleHandler(moodleService port.MoodleService) *MoodleHandler {
	return &MoodleHandler{
		MoodleService: moodleService,
	}
}

// list moodle by email
func (h *MoodleHandler) List(request *gin.Context) {
	email := request.Query("email")
	moodles, err := h.MoodleService.List(email)
	if err != nil {
		request.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	page, _ := helpers.StringToInt(request.Query("page"))
	limit, _ := helpers.StringToInt(request.Query("limit"))
	pagination := helpers.Paginate(moodles, page, limit)
	request.JSON(http.StatusOK, pagination)
}

// get moodle by id
func (h *MoodleHandler) Get(request *gin.Context) {
	moodleID := request.Query("id")
	moodle, err := h.MoodleService.Get(moodleID)
	if err != nil {
		request.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	request.JSON(http.StatusOK, moodle)
}

// search by name from moodle list by email
func (h *MoodleHandler) Search(request *gin.Context) {
	email := request.Query("email")
	name := request.Query("search")
	moodles, err := h.MoodleService.Search(email, name)
	if err != nil {
		request.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	page, _ := helpers.StringToInt(request.Query("page"))
	limit, _ := helpers.StringToInt(request.Query("limit"))
	pagination := helpers.Paginate(moodles, page, limit)
	request.JSON(http.StatusOK, pagination)
}

// list courses
func (h *MoodleHandler) ListCourses(request *gin.Context) {
	courses, err := h.MoodleService.ListCourses()
	if err != nil {
		request.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	page, _ := helpers.StringToInt(request.Query("page"))
	limit, _ := helpers.StringToInt(request.Query("limit"))
	pagination := helpers.Paginate(courses, page, limit)
	request.JSON(http.StatusOK, pagination)
}

 // list packages
func (h *MoodleHandler) ListPackages(request *gin.Context) {
	packages, err := h.MoodleService.ListPackages()
	if err != nil {
		request.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	page, _ := helpers.StringToInt(request.Query("page"))
	limit, _ := helpers.StringToInt(request.Query("limit"))
	pagination := helpers.Paginate(packages, page, limit)
	request.JSON(http.StatusOK, pagination)
}

// create moodle

func (h *MoodleHandler) Create(request *gin.Context) {
	var moodle domain.Moodle
	err := request.ShouldBindJSON(&moodle)
	if err != nil {
		request.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	moodle_, appErr := h.MoodleService.Create(moodle)
	if appErr != nil {
		request.JSON(appErr.StatusCode(), gin.H{"error": appErr.Message})
		return
	}
	request.JSON(http.StatusCreated, moodle_)
}

