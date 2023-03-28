package server

import (
	"moodle/internal/dep"
	"net/http"

	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	docs "moodle/cmd/docs"

	"github.com/gin-gonic/gin"
)

func routes(router *gin.Engine, dependencies *dep.Dep) {
	router.GET("/health", func(request *gin.Context) {
		request.String(http.StatusOK, "OK")
	})
	docs.SwaggerInfo.BasePath = "/api/v1"
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
	v1.DELETE("/moodles", dependencies.MoodleHandler.Delete)
	v1.GET("/logo", dependencies.MoodleHandler.GetLogo)
	v1.POST("/logo", dependencies.MoodleHandler.UpdateLogo)
	v1.GET("/favicon", dependencies.MoodleHandler.GetFavicon)
	v1.POST("/favicon", dependencies.MoodleHandler.UpdateFavicon)
	v1.POST("/pre-installed-course", dependencies.MoodleHandler.UpdatePreInstalledCourse)
	v1.POST("/sitename", dependencies.MoodleHandler.UpdateSiteName)
	v1.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

}
