package dep

import (
	"moodle/internal/core/port"
	"moodle/internal/handler"
)

type Dep struct {
	MoodleService    port.MoodleService
	MoodleHandler    *handler.MoodleHandler
	MongoRepository port.MongoRepository
	MariaRepository port.MariaRepository
	K8sRepository port.K8sRepository
}

