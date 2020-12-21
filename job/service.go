package job

import "github.com/hack-fan/serverfan/types"

type Service interface {
	// Pop a job to server
	Pop(sid string) (*types.ServerJob, error)
	// Push a job by client
	Push(job *types.Job) error
}

type service struct {
}

// NewService create a job service instance
func NewService() Service {
	var s = &service{}
	return s
}
