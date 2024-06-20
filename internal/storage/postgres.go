package storage

import (
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/nglmq/wildberries-0/internal/config"
	"github.com/nglmq/wildberries-0/internal/models"
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
		    order_uid TEXT NOT NULL PRIMARY KEY,
		    track_number TEXT NOT NULL UNIQUE,
		    entry TEXT NOT NULL,
		    locale TEXT NOT NULL,
		    internal_signature TEXT,
		    customer_id TEXT NOT NULL,
		    delivery_service TEXT NOT NULL,
		    shardkey TEXT NOT NULL,
		    sm_id INTEGER NOT NULL,
		    date_created TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		    oof_shard TEXT NOT NULL);`)
	if err != nil {
		return nil, fmt.Errorf("failed to create orders table: %w", err)
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS delivery(
		    order_uid TEXT NOT NULL PRIMARY KEY,
		    name TEXT NOT NULL,
		    phone TEXT NOT NULL,
		    zip TEXT NOT NULL,
		    city TEXT NOT NULL,
		    address TEXT NOT NULL,
		    region TEXT NOT NULL,
		    email TEXT NOT NULL,
		    FOREIGN KEY (order_uid) REFERENCES orders(order_uid)
		);`)
	if err != nil {
		return nil, fmt.Errorf("failed to create delivery table: %w", err)
	}

	_, err = db.Exec(`
  	CREATE TABLE IF NOT EXISTS payment(
  	    transaction TEXT NOT NULL PRIMARY KEY,
  	    order_uid TEXT NOT NULL,
		request_id TEXT,
		currency TEXT NOT NULL,
		provider TEXT NOT NULL,
		amount INTEGER NOT NULL,
		payment_dt TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		bank TEXT NOT NULL,
		delivery_cost INTEGER NOT NULL,
		goods_total INTEGER NOT NULL,
		custom_fee INTEGER NOT NULL DEFAULT 0,
		FOREIGN KEY (order_uid) REFERENCES orders(order_uid)
  	    );`)
	if err != nil {
		return nil, fmt.Errorf("failed to create payment table: %w", err)
	}

	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS items(
	    chrt_id INTEGER NOT NULL PRIMARY KEY,
	    order_uid TEXT NOT NULL,
	    track_number TEXT NOT NULL UNIQUE,
	    price INTEGER NOT NULL,
	    rid TEXT NOT NULL,
	    name TEXT NOT NULL,
	    sale INTEGER NOT NULL DEFAULT 0,
	    size TEXT NOT NULL DEFAULT '0',
	    total_price INTEGER NOT NULL,
	    nm_id INTEGER NOT NULL,
	    brand TEXT NOT NULL,
	    status INTEGER NOT NULL,
	    FOREIGN KEY (track_number) REFERENCES orders(track_number),
	    FOREIGN KEY (order_uid) REFERENCES orders(order_uid)
	    );`)
	if err != nil {
		return nil, fmt.Errorf("failed to create items table: %w", err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) GetOrder(orderID string) (*models.Order, error) {
	var order models.Order
	// Query the order details
	orderQuery := `SELECT order_uid, track_number, entry, locale, internal_signature, customer_id,
                    delivery_service, shardkey, sm_id, date_created, oof_shard
                   FROM orders WHERE order_uid = $1`
	row := s.db.QueryRow(orderQuery, orderID)
	err := row.Scan(&order.OrderID, &order.TrackNumber, &order.Entry, &order.Locale, &order.InternalSignature,
		&order.CustomerID, &order.DeliveryService, &order.ShardKey, &order.SmID,
		&order.DateCreated, &order.OofShard)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch order: %w", err)
	}

	// Query the delivery details
	deliveryQuery := `SELECT name, phone, zip, city, address, region, email FROM delivery WHERE order_uid = $1`
	err = s.db.QueryRow(deliveryQuery, orderID).Scan(&order.Delivery.Name, &order.Delivery.Phone, &order.Delivery.Zip,
		&order.Delivery.City, &order.Delivery.Address, &order.Delivery.Region,
		&order.Delivery.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch delivery details: %w", err)
	}

	// Query the payment details
	paymentQuery := `SELECT transaction, request_id, currency, provider, amount, payment_dt,
                     bank, delivery_cost, goods_total, custom_fee FROM payment WHERE order_uid = $1`
	err = s.db.QueryRow(paymentQuery, orderID).Scan(&order.Payment.Transaction, &order.Payment.RequestID, &order.Payment.Currency,
		&order.Payment.Provider, &order.Payment.Amount, &order.Payment.PaymentDt,
		&order.Payment.Bank, &order.Payment.DeliveryCost, &order.Payment.GoodTotal,
		&order.Payment.CustomFee)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch payment details: %w", err)
	}

	// Query the items details
	itemsQuery := `SELECT chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status FROM items WHERE order_uid = $1`
	rows, err := s.db.Query(itemsQuery, orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch items details: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var item models.Item
		err = rows.Scan(&item.ChrtID, &item.TrackNumber, &item.Price, &item.Rid, &item.Name, &item.Sale,
			&item.Size, &item.TotalPrice, &item.NmID, &item.Brand, &item.Status)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch item details: %w", err)
		}
		order.Items = append(order.Items, item)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return &order, nil
}
