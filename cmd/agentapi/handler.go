package main

import (
	"net/http"
	"time"

	"github.com/hack-fan/jq"
	"github.com/labstack/echo/v4"
	"github.com/rs/xid"

	"github.com/hack-fan/skadi/service"
	"github.com/hack-fan/skadi/types"
)

type Handler struct {
	s  *service.Service
	ev *jq.Queue
}

func NewHandler(s *service.Service, ev *jq.Queue) *Handler {
	return &Handler{
		s:  s,
		ev: ev,
	}
}

func (h *Handler) GetJob(c echo.Context) error {
	aid := c.Get("aid").(string)
	uid := c.Get("uid").(string)
	ip := c.RealIP()
	// async set agent online
	go h.s.AgentOnline(aid, uid, ip)
	// pop a job
	resp := h.s.JobPop(aid)
	if resp != nil {
		return c.JSON(200, resp)
	}
	// no job found
	return c.NoContent(204)
}

func (h *Handler) PutJobSucceed(c echo.Context) error {
	id := c.Param("id")
	req := types.JobResult{}
	err := c.Bind(&req)
	if err != nil {
		return err
	}
	err = h.s.JobSucceed(id, req.Result)
	if err != nil {
		return err
	}
	return c.NoContent(204)
}

func (h *Handler) PutJobFail(c echo.Context) error {
	id := c.Param("id")
	req := types.JobResult{}
	err := c.Bind(&req)
	if err != nil {
		return err
	}
	err = h.s.JobFail(id, req.Result)
	if err != nil {
		return err
	}
	return c.NoContent(204)
}

func (h *Handler) PutJobRunning(c echo.Context) error {
	id := c.Param("id")
	req := types.JobResult{}
	err := c.Bind(&req)
	if err != nil {
		return err
	}
	err = h.s.JobRunning(id, req.Result)
	if err != nil {
		return err
	}
	return c.NoContent(204)
}

func (h *Handler) PostInfo(c echo.Context) error {
	aid := c.Get("aid").(string)
	uid := c.Get("uid").(string)
	input := new(types.MessageInput)
	err := c.Bind(input)
	if err != nil {
		return err
	}
	_, err = h.ev.Pub(&types.Message{
		ID:        xid.New().String(),
		AgentID:   aid,
		UserID:    uid,
		Type:      types.MessageTypeInfo,
		Message:   input.Message,
		CreatedAt: time.Now(),
	})
	if err != nil {
		return err
	}
	return c.NoContent(202)
}

func (h *Handler) PostWarning(c echo.Context) error {
	aid := c.Get("aid").(string)
	uid := c.Get("uid").(string)
	input := new(types.MessageInput)
	err := c.Bind(input)
	if err != nil {
		return err
	}
	_, err = h.ev.Pub(&types.Message{
		ID:        xid.New().String(),
		AgentID:   aid,
		UserID:    uid,
		Type:      types.MessageTypeWarning,
		Message:   input.Message,
		CreatedAt: time.Now(),
	})
	if err != nil {
		return err
	}
	return c.NoContent(202)
}

func (h *Handler) PostText(c echo.Context) error {
	aid := c.Get("aid").(string)
	uid := c.Get("uid").(string)
	input := new(types.MessageInput)
	err := c.Bind(input)
	if err != nil {
		return err
	}
	_, err = h.ev.Pub(&types.Message{
		ID:        xid.New().String(),
		AgentID:   aid,
		UserID:    uid,
		Type:      types.MessageTypeText,
		Message:   input.Message,
		CreatedAt: time.Now(),
	})
	if err != nil {
		return err
	}
	return c.NoContent(202)
}

// API status
func getStatus(c echo.Context) error {
	return c.NoContent(http.StatusNoContent)
}
