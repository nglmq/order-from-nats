package nats

import (
	"context"
	"encoding/json"
	"github.com/nats-io/nats.go"
	"github.com/nglmq/wildberries-0/internal/models"
	"github.com/nglmq/wildberries-0/internal/storage/cache"
	"log/slog"
	"os"
	"os/signal"
)

type OrderSaver interface {
	SaveOrder(ctx context.Context, orderID string, orderInfo models.Order) error
}

func NatsConnect(saver OrderSaver, cache *cache.Cache) error {
	var order models.Order

	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		return err
	}
	defer nc.Close()

	sub, err := nc.Subscribe("orders", func(msg *nats.Msg) {
		err := json.Unmarshal(msg.Data, &order)
		if err != nil {
			slog.Error("failed to unmarshal order")
			return
		}
		slog.Info("received new order")

		cache.SaveToCache(order.OrderID, order)

		err = saver.SaveOrder(context.Background(), order.OrderID, order)
		if err != nil {
			slog.Error("failed to save order")
			return
		}
	})
	if err != nil {
		return err
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)

	<-sig

	if err := sub.Unsubscribe(); err != nil {
		slog.Error("failed to unsubscribe: ", err)
		return err
	}

	return nil
}
