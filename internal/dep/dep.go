package dep

import (
	"moodle/pkg/mariadbiface"
	"moodle/pkg/mongodbiface"

	"moodle/internal/core/port"
	"moodle/internal/handler"
)

type Dep struct {
	MongoDB       mongodbiface.DB
	MariaDB       mariadbiface.MariaDB
	MoodleService    port.MoodleService
	MoodleHandler    *handler.MoodleHandler
	MoodleRepository port.MoodleRepository
	MariaRepository port.MariaRepository
	K8sRepository port.K8sRepository
}

