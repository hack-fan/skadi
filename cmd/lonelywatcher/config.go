package main

import (
	"github.com/hack-fan/x/xdb"
	"github.com/hack-fan/x/xlog"
)

// Settings will load from env and docker secret
type Settings struct {
	Debug    bool `default:"false"`
	Hostname string

	DB xdb.Config

	Redis struct {
		Host     string `default:"redis"`
		Port     string `default:"6379"`
		Password string
		DB       int `default:"0"`
	}

	Wework xlog.WeworkSender
}
