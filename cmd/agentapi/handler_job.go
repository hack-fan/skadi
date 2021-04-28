package main

import (
	"github.com/hack-fan/jq"
	"github.com/hack-fan/skadi/service"
	"github.com/hack-fan/skadi/types"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	s   *service.Service
	evm *jq.Queue
}

func NewHandler(s *service.Service, evm *jq.Queue) *Handler {
	return &Handler{
		s:   s,
		evm: evm,
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

func (h *Handler) PostJobSelf(c echo.Context) error {
	aid := c.Get("aid").(string)
	uid := c.Get("uid").(string)
	req := new(types.MessageInput)
	err := c.Bind(req)
	if err != nil {
		return err
	}
	err = h.s.JobPush(&types.JobInput{
		UserID:  uid,
		AgentID: aid,
		Message: req.Message,
		Source:  "self",
	})
	if err != nil {
		return err
	}
	return c.NoContent(201)
}

func (h *Handler) PostJobAdd(c echo.Context) error {
	aid := c.Get("aid").(string)
	uid := c.Get("uid").(string)
	req := new(types.MessageInput)
	err := c.Bind(req)
	if err != nil {
		return err
	}
	err = h.s.JobAdd(uid, req.Message, "agent:"+aid, "")
	if err != nil {
		return err
	}
	return c.NoContent(201)
}

func (h *Handler) PostJobDelayed(c echo.Context) error {
	aid := c.Get("aid").(string)
	uid := c.Get("uid").(string)
	req := new(types.DelayedJobInput)
	err := c.Bind(req)
	if err != nil {
		return err
	}
	err = h.s.DelayedJobAdd(aid, uid, req)
	if err != nil {
		return err
	}
	return c.NoContent(201)
}
