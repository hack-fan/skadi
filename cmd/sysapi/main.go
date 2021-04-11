// this package is used by local debug only
// in production there would be other upstream service
package main

import (
	"io/ioutil"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/hack-fan/config"
	"github.com/hack-fan/skadi/event"
	"github.com/hack-fan/x/rdb"
	"github.com/hack-fan/x/xdb"
	"github.com/hack-fan/x/xecho"
	"github.com/hack-fan/x/xlog"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/hack-fan/skadi/service"
	"github.com/hack-fan/skadi/types"
)

func main() {
	// load config
	var settings = new(Settings)
	config.MustLoad(settings)

	// logger
	var logger = xlog.New(settings.Debug, settings.Wework)
	defer logger.Sync() // nolint
	var log = logger.Sugar()

	// kv
	rdb.SetLogger(log)
	var kv = rdb.New(settings.Redis)

	// db
	xdb.SetLogger(log)
	var db = xdb.New(settings.DB)
	if settings.Debug {
		db = db.Debug()
	}
	// auto create table
	go db.AutoMigrate(&types.Job{}, &types.Agent{}) // nolint

	// http client
	var rest = resty.New().SetRetryCount(3).
		SetRetryWaitTime(5 * time.Second).
		SetRetryMaxWaitTime(60 * time.Second)

	// default event center is just log the events
	// if you have event worker, set it to redis in settings
	var ev = event.NewEventCenter(log, settings.Event)

	// service
	var s = service.New(kv, db, rest, log, ev)

	// handler
	var h = NewHandler(s)

	// Echo instance
	e := echo.New()
	if settings.Debug {
		e.Debug = true
	}
	// Real IP
	e.IPExtractor = echo.ExtractIPFromXFFHeader()
	// Error handler
	e.HTTPErrorHandler = xecho.NewErrorHandler(logger)
	// Disable echo logs, error handler above will log the error
	e.Logger.SetOutput(ioutil.Discard)
	// Middleware
	e.Use(xecho.ZapLogger(logger))
	e.Use(middleware.Recover())

	// Routes
	e.GET("/", getStatus)

	e.POST("/sys/jobs", h.PostJob)
	e.PUT("/sys/jobs/:id/expire", h.PutJobExpire)

	e.POST("/sys/users/:uid/agents", h.PostAgent)

	// Start server
	e.Logger.Fatal(e.Start(settings.ListenAddr))

}
