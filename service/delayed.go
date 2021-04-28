package service

import (
	"fmt"
	"time"

	"github.com/hack-fan/skadi/types"
	"github.com/hack-fan/x/xerr"
	"github.com/rs/xid"
)

func (s *Service) DelayedJobAdd(aid, uid string, job *types.DelayedJobInput) error {
	// check input
	err := s.validate.Struct(job)
	if err != nil {
		return xerr.Newf(400, "InvalidMessage", "add delayed job failed: %s", err)
	}
	if job.Duration == "" && job.Days == 0 && job.Hours == 0 && job.Minutes == 0 {
		return xerr.New(400, "InvalidDuration", "duration is required")
	}
	var du time.Duration
	if job.Duration != "" {
		du, err = time.ParseDuration(job.Duration)
		if err != nil {
			return xerr.Newf(400, "InvalidDuration", "bad duration format: %s", err)
		}
	} else {
		du = (time.Hour * 24 * time.Duration(job.Days)) +
			(time.Hour * time.Duration(job.Hours)) +
			(time.Minute * time.Duration(job.Minutes))
	}
	// save input
	dj := &types.DelayedJob{
		ID: xid.New().String(),
		JobInput: types.JobInput{
			UserID:   uid,
			AgentID:  aid,
			Message:  job.Message,
			Source:   "delayed",
			Callback: job.Callback,
		},
		ActiveAt: time.Now().Add(du),
	}
	err = s.db.Create(dj).Error
	if err != nil {
		return fmt.Errorf("save delayed job to db failed: %w", err)
	}
	return nil
}
