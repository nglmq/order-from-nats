package storage

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/nglmq/wildberries-0/internal/config"
	"github.com/nglmq/wildberries-0/internal/models"
	"github.com/nglmq/wildberries-0/internal/storage/cache"
)

type Storage struct {
	db *sql.DB
}

func New() (*Storage, error) {
	db, err := sql.Open("pgx", config.DatabaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS orders(
		    order_uid VARCHAR(255) PRIMARY KEY,
		    order_info JSONB NOT NULL);`)
	if err != nil {
		return nil, fmt.Errorf("failed to create order table: %w", err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) SaveOrder(ctx context.Context, orderID string, orderInfo models.Order) error {
	query := `INSERT INTO orders (order_uid, order_info) VALUES ($1, $2)`

	_, err := s.db.ExecContext(ctx, query, orderID, orderInfo)
	if err != nil {
		return fmt.Errorf("failed to save order to database: %w", err)
	}

	return nil
}

func (s *Storage) LoadToCache(ctx context.Context, c *cache.Cache) error {
	rows, err := s.db.QueryContext(ctx, "SELECT order_uid, order_info FROM orders")
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var orderID string
		var orderData []byte
		if err := rows.Scan(&orderID, &orderData); err != nil {
			return err
		}

		var order models.Order
		if err := json.Unmarshal(orderData, &order); err != nil {
			return err
		}

		c.SaveToCache(orderID, order)
	}

	return rows.Err()
}
