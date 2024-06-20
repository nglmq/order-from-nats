package handlers

import (
	"encoding/json"
	"github.com/nglmq/wildberries-0/internal/models"
	"net/http"
)

type OrderGetter interface {
	GetOrder(orderID string) (models.Order, error)
}

func GetOrderHandler(orderGetter OrderGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		orderID := r.URL.Query().Get("orderID")

		order, err := orderGetter.GetOrder(orderID)
		if err != nil {
			http.Error(w, "Order not found", http.StatusNotFound)
			return
		}

		orderJSON, err := json.Marshal(order)
		if err != nil {
			http.Error(w, "Error marshalling order", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(orderJSON)
	}
}
