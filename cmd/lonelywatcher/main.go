package main

import (
	"context"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/hack-fan/config"
	"github.com/hack-fan/jq"
	"github.com/hack-fan/x/rdb"
	"github.com/hack-fan/x/xdb"
	"github.com/hack-fan/x/xlog"

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
	go db.AutoMigrate(&types.Job{}) // nolint

	// http client
	var rest = resty.New().SetRetryCount(3).
		SetRetryWaitTime(5 * time.Second).
		SetRetryMaxWaitTime(60 * time.Second)

	// events
	var evm = jq.NewQueue("skadi:event:"+types.EventMessage, kv)

	// service
	var s = service.New(kv, db, rest, log, evm)

	// watch redis expire events
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()
	pubsub := kv.Subscribe(ctx, "__keyevent@0__:expired")
	log.Info("start watching redis key expired event...")

	log.Warnf("Skadi Watcher started: %s", settings.Hostname)
LOOP:
	for {
		select {
		case msg := <-pubsub.Channel():
			key := msg.Payload
			log.Debugw("redis key expired", "key", key)
			if strings.HasPrefix(key, "job:wait:") {
				jid := strings.TrimPrefix(key, "job:wait:")
				go s.JobExpire(jid)
			} else if strings.HasPrefix(key, "agent:online:") {
				aid := strings.TrimPrefix(key, "agent:online:")
				go s.AgentOffline(aid)
			}
		case <-ctx.Done():
			stop()
			break LOOP
		}
	}
	log.Warnf("Skadi Watcher Stopped: %s", settings.Hostname)
}
