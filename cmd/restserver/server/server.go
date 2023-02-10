package server

import (
	"github.com/gin-gonic/gin"
)

func Start() {
	dependencies := initDependencies()
	router := gin.New()
	router.Use(MoodleMiddleware())
	router.Use(tracing())
	router.Use(logging())

	routes(router, dependencies)
	run(router)
}

func run(router *gin.Engine) {
	if err := router.Run(":5000"); err != nil {
		panic(err)
	}
}
