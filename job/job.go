package job

import (
	"fmt"
	"time"

	"github.com/rs/xid"
	"github.com/vmihailenco/msgpack/v5"

	"github.com/hack-fan/skadi/types"
)

func agentQueueKey(aid string) string {
	return fmt.Sprintf("aq:%s", aid)
}

func (s *Service) Pop(aid string) (*types.JobBasic, error) {
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
	go s.setSent(job.ID)

	return job, nil
}

func (s *Service) Push(input *types.JobInput) error {
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
	go s.store(id, input)

	return nil
}

func (s *Service) Succeed(id string, result string) {
	// change db
	err := s.db.Model(&types.Job{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":       "succeeded",
		"result":       result,
		"succeeded_at": time.Now(),
	}).Error
	if err != nil {
		s.notify(fmt.Errorf("save job %s succeeded status to db failed: %w", id, err))
		return
	}
	// callback
	s.callback(id)
}

func (s *Service) Fail(id string, result string) {
	// change db
	err := s.db.Model(&types.Job{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":    "failed",
		"result":    result,
		"failed_at": time.Now(),
	}).Error
	if err != nil {
		s.notify(fmt.Errorf("save job %s failed status to db failed: %w", id, err))
		return
	}
	// callback
	s.callback(id)
}

func (s *Service) Expire(id string) {
	// change db
	err := s.db.Model(&types.Job{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":     "expired",
		"expired_at": time.Now(),
	}).Error
	if err != nil {
		s.notify(fmt.Errorf("save job %s expired status to db failed: %w", id, err))
		return
	}
	// callback
	s.callback(id)
}

// store async store a job to db
func (s *Service) store(id string, input *types.JobInput) {
	s.log.Infow("new job", "id", id,
		"user", input.UserID, "agent", input.AgentID, "command", input.Message)
	var job = types.Job{
		ID:       id,
		JobInput: *input,
		Status:   "queuing",
	}
	err := s.db.Create(job).Error
	if err != nil {
		s.notify(fmt.Errorf("store new job %s to db failed: %w", id, err))
	}
}

func (s *Service) setSent(id string) {
	s.log.Infow("sent job to agent", "id", id)
	err := s.db.Model(&types.Job{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":  "sent",
		"sent_at": time.Now(),
	}).Error
	if err != nil {
		s.notify(fmt.Errorf("set job %s status to sent failed: %w", id, err))
	}
}

func (s *Service) callback(id string) {
	var job = new(types.Job)
	err := s.db.First(job, "id = ?", id).Error
	if err != nil {
		// TODO: notify back
		s.log.Error(err)
	}
	if job.Callback == "" {
		return
	}
	_, err = s.rest.R().SetBody(job).Post(job.Callback)
	if err != nil {
		s.notify(fmt.Errorf("job %s callback failed: %w", id, err))
	}
}
