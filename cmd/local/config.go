package main

import "github.com/code-gorilla-au/env"

type Config struct {
	DatabaseUrl string
}

func NewConfig() *Config {
	return &Config{
		DatabaseUrl: env.GetAsString("DATABASE_URL"),
	}
}
