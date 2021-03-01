package service

import (
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/hack-fan/x/xerr"
	"github.com/rs/xid"
	"github.com/vmihailenco/msgpack/v5"

	"github.com/hack-fan/skadi/types"
)

func agentQueueKey(aid string) string {
	return fmt.Sprintf("aq:%s", aid)
}

func jobWaitingKey(id string) string {
	return fmt.Sprintf("job:wait:%s", id)
}

// call by agent, so will not return errors, notify the error to admin.
func (s *Service) JobPop(aid string) *types.JobBasic {
	// pop from redis
	var job = new(types.JobBasic)
	data, err := s.kv.RPop(s.ctx, agentQueueKey(aid)).Bytes()
	s.log.Debugw("pop", "data", string(data), "err", err)
	if err == redis.Nil {
		return nil
	} else if err != nil {
		go s.notify(fmt.Errorf("pop job from queue error: %w", err))
		return nil
	}
	err = msgpack.Unmarshal(data, job)
	if err != nil {
		s.notify(fmt.Errorf("msgpack unmarshal job basic error: %w", err))
		return nil
	}
	// for expire count
	err = s.kv.Set(s.ctx, jobWaitingKey(job.ID), "", 10*time.Minute).Err()
	if err != nil {
		s.notify(fmt.Errorf("save job to redis for waiting error: %w", err))
		return nil
	}
	// async change db status
	go s.jobSent(job.ID)

	return job
}

// call by upstream system
func (s *Service) JobPush(input *types.JobInput) error {
	// check agent status
	exists, err := s.kv.Exists(s.ctx, agentOnlineKey(input.AgentID)).Result()
	if err != nil {
		return fmt.Errorf("msgpack marshal input error: %w", err)
	}
	if exists <= 0 {
		return xerr.New(400, "InvalidAgent", "target agent is offline or invalid")
	}
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
	go s.jobStore(id, input)

	return nil
}

func (s *Service) JobSucceed(id string, result string) {
	// change db
	err := s.db.Model(&types.Job{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":       types.JobStatusSucceeded,
		"result":       result,
		"succeeded_at": time.Now(),
	}).Error
	if err != nil {
		s.notify(fmt.Errorf("save job %s succeeded status to db failed: %w", id, err))
		return
	}
	// callback
	s.jobCallback(id)
}

func (s *Service) JobFail(id string, result string) {
	// change db
	err := s.db.Model(&types.Job{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":    types.JobStatusFailed,
		"result":    result,
		"failed_at": time.Now(),
	}).Error
	if err != nil {
		s.notify(fmt.Errorf("save job %s failed status to db failed: %w", id, err))
		return
	}
	// callback
	s.jobCallback(id)
}

func (s *Service) JobExpire(id string) {
	// change db
	err := s.db.Model(&types.Job{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":     types.JobStatusExpired,
		"expired_at": time.Now(),
	}).Error
	if err != nil {
		s.notify(fmt.Errorf("save job %s expired status to db failed: %w", id, err))
		return
	}
	// callback
	s.jobCallback(id)
}

// store async store a job to db
func (s *Service) jobStore(id string, input *types.JobInput) {
	s.log.Infow("new job", "id", id,
		"user", input.UserID, "agent", input.AgentID, "command", input.Message)
	var job = &types.Job{
		ID:       id,
		JobInput: *input,
		Status:   types.JobStatusQueuing,
	}
	err := s.db.Create(job).Error
	if err != nil {
		s.notify(fmt.Errorf("store new job %s to db failed: %w", id, err))
	}
}

func (s *Service) jobSent(id string) {
	s.log.Infow("sent job to agent", "id", id)
	err := s.db.Model(&types.Job{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":  types.JobStatusSent,
		"sent_at": time.Now(),
	}).Error
	if err != nil {
		s.notify(fmt.Errorf("set job %s status to sent failed: %w", id, err))
	}
}

func (s *Service) jobCallback(id string) {
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
