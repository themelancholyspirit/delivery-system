package main

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Storer interface {
	InsertOrder(ctx context.Context, order Order) (Order, error)
	GetOrder(ctx context.Context, id int) (Order, error)
	UpdateOrder(ctx context.Context, order Order) error
	GetOrders(ctx context.Context, limit, offset int) ([]Order, error)
}

type PostgreStorer struct {
	DB *gorm.DB
}

// for testing purposes, sqlite is used.
type SQLiteStorer struct {
	db *gorm.DB
}

// NewSQLiteStorer initializes a new SQLite database connection and returns a SQLiteStorer.
func NewSQLiteStorer() (*SQLiteStorer, error) {
	db, err := gorm.Open(sqlite.Open("orders.db"), &gorm.Config{})

	if err != nil {
		return nil, err
	}

	// Automatically migrate the schema, creating the table if it doesn't exist
	if err := db.AutoMigrate(&Order{}); err != nil {
		return nil, err
	}

	return &SQLiteStorer{db: db}, nil
}

func NewPostgreStorer() (*PostgreStorer, error) {
	dsn := "host=postgres user=myuser password=mypassword dbname=mydatabase port=5432 sslmode=disable"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		return nil, err
	}

	// Auto-migrate the Order schema
	if err := db.AutoMigrate(&Order{}); err != nil {
		return nil, err
	}

	return &PostgreStorer{DB: db}, nil

}

// InsertOrder inserts a new order into the database

func (ps *PostgreStorer) InsertOrder(ctx context.Context, order Order) (Order, error) {
	order.Version = 1
	if err := ps.DB.WithContext(ctx).Create(&order).Error; err != nil {
		return Order{}, err
	}
	return order, nil
}

// GetOrder retrieves an order by its ID
func (ps *PostgreStorer) GetOrder(ctx context.Context, id int) (Order, error) {
	var order Order
	if err := ps.DB.WithContext(ctx).First(&order, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return Order{}, fmt.Errorf("order with ID %d not found", id)
		}
		return Order{}, err
	}
	return order, nil
}

// ModifyOrder updates an existing order with new data

func (ps *PostgreStorer) UpdateOrder(ctx context.Context, order Order) error {
	result := ps.DB.Model(&order).
		Where("id = ? AND version = ?", order.ID, order.Version).
		Updates(map[string]interface{}{
			"status":  order.Status,
			"version": gorm.Expr("version + 1"),
		})

	if result.RowsAffected == 0 {
		return errors.New("record not found") // Handle the case where the record has changed
	}
	return result.Error
}

func (ps *PostgreStorer) GetOrders(ctx context.Context, limit, offset int) ([]Order, error) {
	var orders []Order
	// Fetch orders with limit and offset
	if err := ps.DB.Offset(offset).Limit(limit).Find(&orders).Error; err != nil {
		return nil, err // Return any error encountered
	}
	return orders, nil // Return the fetched orders
}

// InsertOrder inserts a new order into the database and returns the created order.
func (s *SQLiteStorer) InsertOrder(ctx context.Context, order Order) (Order, error) {
	if err := s.db.WithContext(ctx).Create(&order).Error; err != nil {
		return Order{}, err
	}
	return order, nil
}

// GetOrder retrieves an order by its ID from the database.
func (s *SQLiteStorer) GetOrder(ctx context.Context, id int) (Order, error) {
	var order Order
	if err := s.db.WithContext(ctx).First(&order, id).Error; err != nil {
		return Order{}, err
	}
	return order, nil
}

// UpdateOrder updates an existing order in the database.
func (s *SQLiteStorer) UpdateOrder(ctx context.Context, order Order) error {
	if err := s.db.WithContext(ctx).Save(&order).Error; err != nil {
		return err
	}
	return nil
}

// GetOrders retrieves a list of orders with pagination.
func (s *SQLiteStorer) GetOrders(ctx context.Context, limit, offset int) ([]Order, error) {
	var orders []Order
	if err := s.db.WithContext(ctx).Limit(limit).Offset(offset).Find(&orders).Error; err != nil {
		return nil, err
	}
	return orders, nil
}
