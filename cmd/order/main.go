package main

import (
	"github.com/nglmq/wildberries-0/internal/config"
	"github.com/nglmq/wildberries-0/internal/server"
	"log"
	"log/slog"
	"net/http"
)

func main() {
	r, err := server.Start()
	if err != nil {
		log.Fatal(err)
	}

	slog.Info("Listening on http://" + config.RunAddr)
	log.Fatal(http.ListenAndServe(config.RunAddr, r))
}
