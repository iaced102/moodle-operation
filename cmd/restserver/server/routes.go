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

	v1.GET("/courses", dependencies.MoodleHandler.ListCourses)
	v1.GET("/moodles", dependencies.MoodleHandler.List)
	// router.POST("/users/:user_id/games", dependencies.MoodleHandler.Create)
	// router.GET("/users/:user_id/games", dependencies.MoodleHandler.GetAll)
	// router.GET("/users/:user_id/games/:game_id", dependencies.MoodleHandler.Get)
	// router.PUT("/users/:user_id/games/:game_id/actions/reveal", dependencies.MoodleHandler.Reveal)
	// router.PUT("/users/:user_id/games/:game_id/actions/mark", dependencies.MoodleHandler.Mark)
}
