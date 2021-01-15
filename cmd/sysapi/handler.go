package main

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/hack-fan/skadi/service"
	"github.com/hack-fan/skadi/types"
)

type Handler struct {
	js *service.Service
}

func NewHandler(js *service.Service) *Handler {
	return &Handler{
		js: js,
	}
}

func (h *Handler) PostJob(c echo.Context) error {
	var req = new(types.JobInput)
	err := c.Bind(req)
	if err != nil {
		return err
	}
	err = h.js.JobPush(req)
	if err != nil {
		return err
	}
	return c.NoContent(204)
}

func (h *Handler) PutJobExpire(c echo.Context) error {
	id := c.Param("id")
	h.js.JobExpire(id)
	return c.NoContent(204)
}

func (h *Handler) PostAgent(c echo.Context) error {
	uid := c.Param("uid")
	agent, err := h.js.AgentAdd(uid)
	if err != nil {
		return err
	}
	return c.JSON(201, agent)
}

// API status
func getStatus(c echo.Context) error {
	return c.NoContent(http.StatusNoContent)
}
