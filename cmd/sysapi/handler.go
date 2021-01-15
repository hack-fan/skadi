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
	err = h.js.Push(req)
	if err != nil {
		return err
	}
	return c.NoContent(204)
}

func (h *Handler) PutJobExpire(c echo.Context) error {
	id := c.Param("id")
	h.js.Expire(id)
	return c.NoContent(204)
}

// API status
func getStatus(c echo.Context) error {
	return c.NoContent(http.StatusNoContent)
}
