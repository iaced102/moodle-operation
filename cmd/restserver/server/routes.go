package server

import (
	"moodle/internal/dep"
	"net/http"

	"github.com/gin-gonic/gin"
)

func routes(router *gin.Engine, dependencies *dep.Dep) {
	router.GET("/ping", func(request *gin.Context) {
		request.String(http.StatusOK, "pong")
	})

	// list moodle by email
	router.GET("/moodles", dependencies.MoodleHandler.List)
	router.GET("/courses", dependencies.MoodleHandler.ListCourses)
	// router.POST("/users/:user_id/games", dependencies.MoodleHandler.Create)
	// router.GET("/users/:user_id/games", dependencies.MoodleHandler.GetAll)
	// router.GET("/users/:user_id/games/:game_id", dependencies.MoodleHandler.Get)
	// router.PUT("/users/:user_id/games/:game_id/actions/reveal", dependencies.MoodleHandler.Reveal)
	// router.PUT("/users/:user_id/games/:game_id/actions/mark", dependencies.MoodleHandler.Mark)
}
