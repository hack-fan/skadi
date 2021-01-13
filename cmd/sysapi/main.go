package main

import (
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/hack-fan/config"
	"github.com/hack-fan/x/xdb"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/hack-fan/skadi/job"
)

func main() {
	// load config
	var settings = new(Settings)
	config.MustLoad(settings)

	// kv
	var kv = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", settings.Redis.Host, settings.Redis.Port),
		Password: settings.Redis.Password,
		DB:       settings.Redis.DB,
	})

	// db
	var db = xdb.New(settings.DB)

	// job service
	var js = job.NewService(kv, db, nil, nil)

	// handler
	var h = NewHandler(js)

	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.GET("/", getStatus)

	e.POST("/sys/jobs", h.PostJob)

	// Start server
	e.Logger.Fatal(e.Start(":1323"))

}
