package handler

import (
	"moodle/internal/core/domain"
	"moodle/internal/core/port"
	"net/http"
	"strconv"

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

// update vision image
func (h *MoodleHandler) UpdateVisionImage(request *gin.Context) {
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
	// update vision image
	resp, appErr := h.MoodleService.UpdateVisionImage(moodleID, file, header)
	if appErr != nil {
		request.JSON(appErr.StatusCode(), gin.H{"error": appErr.Message})
		return
	}
	request.JSON(http.StatusOK, resp)
}

// update banner image
func (h *MoodleHandler) UpdateBannerImage(request *gin.Context) {
	moodleID := request.Query("moodle_id")
	bannerID := request.Query("banner_id")
	bannerIDInt, err := strconv.Atoi(bannerID)
	if err != nil {
		request.JSON(http.StatusBadRequest, gin.H{"error": "banner_id must be a number"})
		return
	}
	// validate if banner not in range 1-3
	if bannerIDInt < 1 || bannerIDInt > 3 {
		request.JSON(http.StatusBadRequest, gin.H{"error": "banner_id must be in range 1-3"})
		return
	}
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
	// update banner image
	resp, appErr := h.MoodleService.UpdateBannerImage(moodleID, bannerIDInt, file, header)
	if appErr != nil {
		request.JSON(appErr.StatusCode(), gin.H{"error": appErr.Message})
		return
	}
	request.JSON(http.StatusOK, resp)
}

// update video url
func (h *MoodleHandler) UpdateVideoURL(request *gin.Context) {
	var  update domain.VideoURLTracking
	err := request.ShouldBindJSON(&update)
	if err != nil {
		request.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// chekc if video url is valid
	isValid := helpers.IsURL(update.VideoURL)
	if !isValid {
		request.JSON(http.StatusBadRequest, gin.H{"error": "video url is not valid"})
		return
	}
	// update video url
	resp, appErr := h.MoodleService.UpdateVideoURL(update.MoodleId, update.VideoURL)
	if appErr != nil {
		request.JSON(appErr.StatusCode(), gin.H{"error": appErr.Message})
		return
	}
	request.JSON(http.StatusOK, resp)
}

// update vision content
func (h *MoodleHandler) UpdateVisionContent(request *gin.Context) {
	var  update domain.VisionContentTracking
	err := request.ShouldBindJSON(&update)
	if err != nil {
		request.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// update vision content
	resp, appErr := h.MoodleService.UpdateVisionContent(update)
	if appErr != nil {
		request.JSON(appErr.StatusCode(), gin.H{"error": appErr.Message})
		return
	}
	request.JSON(http.StatusOK, resp)
}
