package main

import (
	"encoding/json"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/nats-io/nats.go"
	"github.com/nglmq/wildberries-0/internal/models"
	"log/slog"
	"os"
	"os/signal"
	"time"
)

func main() {
	var order models.Order

	sc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		slog.Error("unable to connect to NATS", err)
		return
	}

	go func() {
		for {
			time.Sleep(time.Second * 2)

			err = gofakeit.Struct(&order)
			if err != nil {
				slog.Error("error marshalling fake order", err)
				return
			}

			jsonToSend, err := json.MarshalIndent(order, "", " ")
			if err != nil {
				slog.Error("error marshalling json", err)
				return
			}

			err = sc.Publish("orders", jsonToSend)
			if err != nil {
				slog.Error("error while publishing to NATS channel", err)
				return
			}
			slog.Info("send order successfully: ", order.OrderID)
		}
	}()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)

	<-sig
	slog.Info("Stop generating orders...")

	sc.Close()
}
