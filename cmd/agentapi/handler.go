package main

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/rs/xid"

	"github.com/hack-fan/skadi/service"
	"github.com/hack-fan/skadi/types"
)

type Handler struct {
	s  *service.Service
	ev types.EventCenter
}

func NewHandler(s *service.Service, ev types.EventCenter) *Handler {
	return &Handler{
		s:  s,
		ev: ev,
	}
}

func (h *Handler) GetJob(c echo.Context) error {
	aid := c.Get("aid").(string)
	ip := c.RealIP()
	// async set agent online
	go h.s.AgentOnline(aid, ip)
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
	h.s.JobSucceed(id, req.Result)
	return c.NoContent(204)
}

func (h *Handler) PutJobFail(c echo.Context) error {
	id := c.Param("id")
	req := types.JobResult{}
	err := c.Bind(&req)
	if err != nil {
		return err
	}
	h.s.JobFail(id, req.Result)
	return c.NoContent(204)
}

func (h *Handler) PostInfo(c echo.Context) error {
	aid := c.Get("aid").(string)
	uid := c.Get("uid").(string)
	input := new(types.EventInput)
	err := c.Bind(input)
	if err != nil {
		return err
	}
	err = h.ev.Pub(&types.Event{
		ID:        xid.New().String(),
		AgentID:   aid,
		UserID:    uid,
		Type:      types.EventTypeInfo,
		Message:   input.Message,
		CreatedAt: time.Now(),
	})
	if err != nil {
		return err
	}
	return c.NoContent(201)
}

func (h *Handler) PostWarning(c echo.Context) error {
	aid := c.Get("aid").(string)
	uid := c.Get("uid").(string)
	input := new(types.EventInput)
	err := c.Bind(input)
	if err != nil {
		return err
	}
	err = h.ev.Pub(&types.Event{
		ID:        xid.New().String(),
		AgentID:   aid,
		UserID:    uid,
		Type:      types.EventTypeWarning,
		Message:   input.Message,
		CreatedAt: time.Now(),
	})
	if err != nil {
		return err
	}
	return c.NoContent(201)
}

// API status
func getStatus(c echo.Context) error {
	return c.NoContent(http.StatusNoContent)
}
