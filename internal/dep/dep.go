package dep

import (
	"moodle/pkg/mariadbiface"
	"moodle/pkg/mongodbiface"

	"moodle/internal/core/port"
	"moodle/internal/handler"
)

type Dep struct {
	MongoDB       mongodbiface.MongoDB
	MariaDB       mariadbiface.MariaDB
	MoodleService    port.MoodleService
	MoodleHandler    *handler.MoodleHandler
	MoodleRepository port.MoodleRepository
	MariaRepository port.MoodleRepository
}

