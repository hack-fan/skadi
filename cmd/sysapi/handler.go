package main

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/hack-fan/skadi/service"
	"github.com/hack-fan/skadi/types"
)

type Handler struct {
	s *service.Service
}

func NewHandler(s *service.Service) *Handler {
	return &Handler{
		s: s,
	}
}

func (h *Handler) PostJob(c echo.Context) error {
	var req = new(types.JobInput)
	err := c.Bind(req)
	if err != nil {
		return err
	}
	err = h.s.JobPush(req)
	if err != nil {
		return err
	}
	return c.NoContent(204)
}

func (h *Handler) PutJobExpire(c echo.Context) error {
	id := c.Param("id")
	h.s.JobExpire(id)
	return c.NoContent(204)
}

func (h *Handler) PostAgent(c echo.Context) error {
	uid := c.Param("uid")
	basic := new(types.AgentBasic)
	err := c.Bind(basic)
	if err != nil {
		return err
	}
	agent, err := h.s.AgentAdd(uid, basic)
	if err != nil {
		return err
	}
	return c.JSON(201, agent)
}

// API status
func getStatus(c echo.Context) error {
	return c.NoContent(http.StatusNoContent)
}
