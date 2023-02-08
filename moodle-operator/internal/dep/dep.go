package dep

import (
	"moodle/pkg/mariadbiface"
	"moodle/pkg/mongodbiface"

	"github.com/matiasvarela/minesweeper-API/internal/core/port"
	"github.com/matiasvarela/minesweeper-API/internal/handler"
)

type Dep struct {
	MongoDB       mongodbiface.MongoDB
	MariaDB       mariadbiface.MariaDB
	GameService    port.GameService
	GameHandler    *handler.GameHandler
	MoodleRepository port.MoodleRepository
}

