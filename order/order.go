package order

import (
	"context"
)

type Product struct {
	Name  string `json:"name" firestore:"name"`
	Count int    `json:"count" firestore:"count"`
	Price int    `json:"price" firestore:"price"`
}

// Order defines the order data structure.
type Order struct {
	ID         string    `json:"id" firestore:"id"`
	Products   []Product `json:"products" firestore:"products"`
	Name       string    `json:"name" firestore:"name"`
	Mail       string    `json:"mail" firestore:"mail"`
	Note       string    `json:"note" firestore:"note"`
	TotalPrice int       `json:"total_price" firestore:"total_price"`
	IsPaid     bool      `json:"is_paid" firestore:"is_paid"`
	IsPicked   bool      `json:"is_picked" firestore:"is_picked"`
	IsEnabled  bool      `json:"is_enabled" firestore:"is_enabled"`
	PaidAt     string    `json:"paid_at" firestore:"paid_at"`
	PickedAt   string    `json:"picked_at" firestore:"picked_at"`
	CreatedAt  string    `json:"created_at" firestore:"created_at"`
}

type CreateOrderRequest struct {
	Mail     string    `json:"mail"`
	Name     string    `json:"name"`
	Products []Product `json:"products"`
}

// Service provides order operations.
type Service interface {
	GetOrders(ctx context.Context, page, pageSize int, search string) ([]Order, int, error)
	UpdateOrder(ctx context.Context, id string, data map[string]interface{}) (Order, error)
	CreateOrder(ctx context.Context, req CreateOrderRequest) (Order, error)
}
