package server

import (
	"github.com/go-chi/chi/v5"
	"github.com/nglmq/wildberries-0/internal/config"
	"github.com/nglmq/wildberries-0/internal/handlers"
	"github.com/nglmq/wildberries-0/internal/storage"
	"log/slog"
	"net/http"
)

func Start() (http.Handler, error) {
	config.ParseFlags()

	storage, err := storage.New()
	if err != nil {
		slog.Error("failed to init db")
		return nil, err
	}

	r := chi.NewRouter()

	r.Route("/order", func(r chi.Router) {
		r.Get("/{orderID}", handlers.GetOrderHandler(storage))
	})

	return nil
}
