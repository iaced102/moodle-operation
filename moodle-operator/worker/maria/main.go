package worker

import (
	adapter "moodle/adapter/mariadb"
)

// new mariadb worker
type MariaWorker struct {
	Adapter *adapter.MariaAdapter
}

// new mariadb worker
func NewMariaWorker(adapter *adapter.MariaAdapter) *MariaWorker {
	return &MariaWorker{
		Adapter: adapter,
	}
}

