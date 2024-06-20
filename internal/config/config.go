package config

import (
	"flag"
	"os"
)

var (
	RunAddr     string
	DatabaseURL string
)

func ParseFlags() {
	flag.StringVar(&RunAddr, "a", "localhost:8080", "server run address")
	flag.StringVar(&DatabaseURL, "d", "", "postgres connection url")
	flag.Parse()

	endRunAddr := os.Getenv("SERVER_ADDRESS")
	if endRunAddr != "" {
		RunAddr = endRunAddr
	}

	envDBConnection := os.Getenv("DATABASE_URI")
	if envDBConnection != "" {
		DatabaseURL = envDBConnection
	}
}
