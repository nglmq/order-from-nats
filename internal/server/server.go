package server

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/nglmq/wildberries-0/internal/config"
	"github.com/nglmq/wildberries-0/internal/handlers"
	"github.com/nglmq/wildberries-0/internal/nats"
	"github.com/nglmq/wildberries-0/internal/storage"
	"github.com/nglmq/wildberries-0/internal/storage/cache"
	"log/slog"
	"net/http"
	"time"
)

func Start() (http.Handler, error) {
	config.ParseFlags()

	store, err := storage.New()
	if err != nil {
		slog.Error("failed to init db: ", err)
		return nil, err
	}

	newCache := cache.NewCache()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	err = store.LoadToCache(ctx, newCache)
	if err != nil {
		slog.Error("failed to load data to cache: ", err)
		return nil, err
	}

	go func() {
		err = nats.NatsConnect(store, newCache)
		if err != nil {
			slog.Error("failed to connect to NATS: ", err)
			return
		}
	}()

	r := chi.NewRouter()

	r.Route("/", func(r chi.Router) {
		r.Get("/", handlers.TemplateHandler())
		r.Post("/", handlers.GetOrderHandler(newCache))
	})

	return r, nil
}
