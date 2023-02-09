package handler

import (
	"moodle/internal/core/port"
	"net/http"

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
	request.JSON(http.StatusOK, moodles)
}


// list courses
func (h *MoodleHandler) ListCourses(request *gin.Context) {
	courses, err := h.MoodleService.ListCourses()
	if err != nil {
		request.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	request.JSON(http.StatusOK, courses)
}
