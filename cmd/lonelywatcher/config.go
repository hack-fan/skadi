package main

import (
	"github.com/hack-fan/skadi/event"
	"github.com/hack-fan/x/rdb"
	"github.com/hack-fan/x/xdb"
	"github.com/hack-fan/x/xlog"
)

// Settings will load from env and docker secret
type Settings struct {
	Debug    bool `default:"false"`
	Hostname string

	DB xdb.Config

	Redis rdb.Config

	Event event.Config

	Wework xlog.WeworkSender
}
