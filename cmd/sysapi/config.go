package main

import "github.com/hack-fan/x/xdb"

// Settings will load from env and docker secret
type Settings struct {
	Debug bool `default:"false"`

	DB xdb.Config

	Redis struct {
		Host     string `default:"redis"`
		Port     string `default:"6379"`
		Password string
		DB       int `default:"0"`
	}
}
