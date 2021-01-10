package job

import (
	"fmt"

	"github.com/rs/xid"
	"github.com/vmihailenco/msgpack/v5"

	"github.com/hack-fan/skadi/types"
)

func agentQueueKey(aid string) string {
	return fmt.Sprintf("aq:%s", aid)
}

func (s *service) Pop(aid string) (*types.JobBasic, error) {
	// pop from redis
	var job = new(types.JobBasic)
	data, err := s.kv.RPop(s.ctx, agentQueueKey(aid)).Bytes()
	if err != nil {
		return nil, fmt.Errorf("pop job from queue error: %w", err)
	}
	err = msgpack.Unmarshal(data, job)
	if err != nil {
		return nil, fmt.Errorf("msgpack unmarshal job basic error: %w", err)
	}
	// async change db status
	go s.SetSent(job.ID)

	return job, nil
}

func (s *service) Push(input *types.JobInput) error {
	// check agent status
	// gen id
	var id = xid.New().String()
	// save to kv
	data, err := msgpack.Marshal(&types.JobBasic{
		ID:      id,
		Message: input.Message,
	})
	if err != nil {
		return fmt.Errorf("msgpack marshal input error: %w", err)
	}
	err = s.kv.LPush(s.ctx, agentQueueKey(input.AgentID), data).Err()
	if err != nil {
		return fmt.Errorf("push input to agent queue error: %w", err)
	}
	// save to db
	go s.Store(id, input)

	return nil
}

// Store async store a job to db
func (s *service) Store(id string, input *types.JobInput) {
	var job = types.Job{
		ID:       id,
		JobInput: *input,
		Status:   "queuing",
	}
	err := s.db.Create(job).Error
	if err != nil {
		// TODO: notify back
		s.log.Error(err)
	}
}

func (s *service) SetSent(id string) {
	err := s.db.Model(&types.Job{}).Where("id = ?", id).Update("status", "sent").Error
	if err != nil {
		// TODO: notify back
		s.log.Error(err)
	}
}
