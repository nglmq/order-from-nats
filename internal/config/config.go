package config

import (
	"flag"
	"os"
)

var (
	DatabaseURL string
)

func ParseFlags() {
	flag.StringVar(&DatabaseURL, "d", "", "postgres connection url")
	flag.Parse()

	envDBConnection := os.Getenv("DATABASE_URI")
	if envDBConnection != "" {
		DatabaseURL = envDBConnection
	}
}
