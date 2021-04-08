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

var jobWaitingTime = time.Minute * 10

// call by agent, so will not return errors, notify the error to admin.
func (s *Service) JobPop(aid string) *types.JobBasic {
	// pop from redis
	var job = new(types.JobBasic)
	data, err := s.kv.RPop(s.ctx, agentQueueKey(aid)).Bytes()
	s.log.Debugw("pop", "data", string(data), "err", err)
	if err == redis.Nil {
		return nil
	} else if err != nil {
		s.log.Errorf("pop job from queue error: %s", err)
		return nil
	}
	err = msgpack.Unmarshal(data, job)
	if err != nil {
		s.log.Errorf("msgpack unmarshal job basic error: %s", err)
		return nil
	}
	// for expire count
	err = s.kv.Set(s.ctx, jobWaitingKey(job.ID), 0, jobWaitingTime).Err()
	if err != nil {
		s.log.Errorf("save job to redis for waiting error: %s", err)
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

func (s *Service) Job(id string) (*types.Job, error) {
	var job = new(types.Job)
	err := s.db.First(job, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return job, nil
}

func (s *Service) JobRunning(id string, message string) error {
	// check status
	cnt, err := s.kv.Exists(s.ctx, jobWaitingKey(id)).Result()
	if err != nil {
		s.log.Errorf("job running report redis error: %s", err)
		return err
	}
	if cnt <= 0 {
		return xerr.Newf(400, "JobNotFound", "job %s is not active", id)
	}
	// refresh ttl
	err = s.kv.Expire(s.ctx, jobWaitingKey(id), jobWaitingTime).Err()
	if err != nil {
		s.log.Errorf("refresh job ttl error: %s", err)
		return err
	}
	// save to db
	err = s.db.Model(&types.Job{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status": types.JobStatusRunning,
		"result": message,
	}).Error
	if err != nil {
		s.log.Errorf("save job %s running status to db failed: %s", id, err)
		return err
	}
	// callback
	s.jobCallback(id)
	return nil
}

func (s *Service) JobSucceed(id string, result string) error {
	// check status
	job, err := s.Job(id)
	if err != nil {
		s.log.Errorf("check job %s status failed: %s", id, err)
		return err
	}
	if job.Status == types.JobStatusSucceeded || job.Status == types.JobStatusFailed {
		return xerr.Newf(400, "JobFinished", "the job [%s] has been finished", job.Message)
	}
	// change db
	err = s.db.Model(&types.Job{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":       types.JobStatusSucceeded,
		"result":       result,
		"succeeded_at": time.Now(),
	}).Error
	if err != nil {
		s.log.Errorf("save job %s succeeded status to db failed: %s", id, err)
		return err
	}
	// rm waiting key
	err = s.kv.Del(s.ctx, jobWaitingKey(id)).Err()
	if err != nil {
		s.log.Errorf("delete job %s waiting key after succeeded error: %s", id, err)
		return err
	}
	// callback
	s.jobCallback(id)
	return nil
}

func (s *Service) JobFail(id string, result string) error {
	// check status
	job, err := s.Job(id)
	if err != nil {
		s.log.Errorf("check job %s status failed: %s", id, err)
		return err
	}
	if job.Status == types.JobStatusSucceeded || job.Status == types.JobStatusFailed {
		return xerr.Newf(400, "JobFinished", "the job [%s] has been finished", job.Message)
	}
	// change db
	err = s.db.Model(&types.Job{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":    types.JobStatusFailed,
		"result":    result,
		"failed_at": time.Now(),
	}).Error
	if err != nil {
		s.log.Errorf("save job %s failed status to db failed: %s", id, err)
		return err
	}
	// rm waiting key
	err = s.kv.Del(s.ctx, jobWaitingKey(id)).Err()
	if err != nil {
		s.log.Errorf("delete job %s waiting key after failed error: %s", id, err)
		return err
	}
	// callback
	s.jobCallback(id)
	return nil
}

func (s *Service) JobExpire(id string) {
	// change db
	err := s.db.Model(&types.Job{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":     types.JobStatusExpired,
		"expired_at": time.Now(),
	}).Error
	if err != nil {
		s.log.Errorf("save job %s expired status to db failed: %s", id, err)
		return
	}
	// callback
	s.jobCallback(id)
}

// Agent offline will cancel all job in queue
func (s *Service) JobCancel(id string) {
	// change db
	err := s.db.Model(&types.Job{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":      types.JobStatusCanceled,
		"canceled_at": time.Now(),
	}).Error
	if err != nil {
		s.log.Errorf("save job %s expired status to db failed: %s", id, err)
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
		s.log.Errorf("store new job %s to db failed: %s", id, err)
	}
}

func (s *Service) jobSent(id string) {
	s.log.Infow("sent job to agent", "id", id)
	err := s.db.Model(&types.Job{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":  types.JobStatusSent,
		"sent_at": time.Now(),
	}).Error
	if err != nil {
		s.log.Errorf("set job %s status to sent failed: %s", id, err)
	}
	// callback
	s.jobCallback(id)
}

func (s *Service) jobCallback(id string) {
	var job = new(types.Job)
	err := s.db.First(job, "id = ?", id).Error
	if err != nil {
		s.log.Errorf("job %s callback fetch job from db failed: %s", id, err)
		return
	}
	if job.Callback == "" {
		return
	}
	_, err = s.rest.R().SetBody(job).Post(job.Callback)
	if err != nil {
		s.log.Errorf("job %s callback failed: %s", id, err)
		return
	}
}
