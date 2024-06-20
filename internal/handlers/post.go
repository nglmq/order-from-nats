package handlers

import (
	"encoding/json"
	"github.com/nglmq/wildberries-0/internal/models"
	"net/http"
)

type OrderGetter interface {
	GetFromCache(orderID string) (models.Order, bool)
}

func GetOrderHandler(orderGetter OrderGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		orderID := r.PostFormValue("order_id")

		order, exists := orderGetter.GetFromCache(orderID)
		if !exists {
			http.Error(w, "Order not found", http.StatusInternalServerError)
			return
		}

		orderJSON, err := json.MarshalIndent(order, "", "  ")
		if err != nil {
			http.Error(w, "Error marshalling order", http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(orderJSON)
	}
}
