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

// delete moodle
func (h *MoodleHandler) Delete(request *gin.Context) {
	var moodle domain.MariaTracking
	err := request.ShouldBindJSON(&moodle)
	if err != nil {
		request.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err = h.MoodleService.Delete(moodle.MoodleId)
	if err != nil {
		request.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	request.JSON(http.StatusNoContent, gin.H{"message": "moodle deleted successfully"})
}


// Update logo
func (h *MoodleHandler) UpdateLogo(request *gin.Context) {
	moodleID := request.Query("moodle_id")
	// get multipart file, multipart file header
	file, header, err := request.Request.FormFile("file")
	isValid := helpers.IsImage(file)
	if !isValid {
		request.JSON(http.StatusBadRequest, gin.H{"error": "file must be an image: jpg, png"})
		return
	}
	if err != nil {
		request.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// update logo
	resp, appErr := h.MoodleService.UpdateLogo(moodleID, file, header)
	if appErr != nil {
		request.JSON(appErr.StatusCode(), gin.H{"error": appErr.Message})
		return
	}
	request.JSON(http.StatusOK, resp)
}

// Update favicon
func (h *MoodleHandler) UpdateFavicon(request *gin.Context) {
	moodleID := request.Query("moodle_id")
	// get multipart file, multipart file header
	file, header, err := request.Request.FormFile("file")
	isValid := helpers.IsImage(file)
	if !isValid {
		request.JSON(http.StatusBadRequest, gin.H{"error": "file must be an image: jpg, png"})
		return
	}
	if err != nil {
		request.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// update favicon
	resp, appErr := h.MoodleService.UpdateFavicon(moodleID, file, header)
	if appErr != nil {
		request.JSON(appErr.StatusCode(), gin.H{"error": appErr.Message})
		return
	}
	request.JSON(http.StatusOK, resp)
}

// update pre_installed course
func (h *MoodleHandler) UpdatePreInstalledCourse(request *gin.Context) {
	var  update domain.Moodle
	err := request.ShouldBindJSON(&update)
	if err != nil {
		request.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// validate if all of pre_installed_course is in range 1-8
	for _, course := range update.PreInstalledCourse {
		if course < 1 || course > 8 {
			request.JSON(http.StatusBadRequest, gin.H{"error": "pre_installed_course must be in range 1-8"})
			return
		}
	}
	// update pre_installed course
	resp, appErr := h.MoodleService.UpdatePreInstalledCourse(update)
	if appErr != nil {
		request.JSON(appErr.StatusCode(), gin.H{"error": appErr.Message})
		return
	}
	request.JSON(http.StatusOK, resp)
}
