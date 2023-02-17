package worker

import (
	repo "moodle/internal/repository/moodle"
)

// k8s worker
type K8sWorker struct {
	mongo *repo.MongoDB
	k8s   *repo.K8sClient

}

func NewK8sWorker(mongo *repo.MongoDB, k8s *repo.K8sClient) *K8sWorker {
	return &K8sWorker{
		mongo: mongo,
		k8s:   k8s,
	}
}
