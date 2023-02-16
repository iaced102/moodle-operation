package server

import (
	"moodle/internal/dep"
	"net/http"

	"github.com/gin-gonic/gin"
)

func routes(router *gin.Engine, dependencies *dep.Dep) {
	router.GET("/health", func(request *gin.Context) {
		request.String(http.StatusOK, "OK")
	})
	v1 := router.Group("/api/v1")

	v1.GET("/moodles",  func(request *gin.Context) {
		get := request.Query("id")
		if get != "" {
			dependencies.MoodleHandler.Get(request)
			return
		}
		search := request.Query("search")
		if search != "" {
			dependencies.MoodleHandler.Search(request)
			return
		}
		dependencies.MoodleHandler.List(request)
	})
	v1.GET("/packages", dependencies.MoodleHandler.ListPackages)
	v1.GET("/courses", dependencies.MoodleHandler.ListCourses)
	v1.POST("/moodles", dependencies.MoodleHandler.Create)
	v1.POST("/logo", dependencies.MoodleHandler.UpdateLogo)
	v1.POST("/favicon", dependencies.MoodleHandler.UpdateFavicon)
}
