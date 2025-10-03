package coupon

import (
	"context"
)

// Coupon defines the coupon data structure.
type Coupon struct {
	ID        string `json:"id" firestore:"id"`
	Name      string `json:"name" firestore:"name"`
	Code      string `json:"code" firestore:"code"`
	Percent   int    `json:"percent" firestore:"percent"`
	StartTime int64  `json:"start_time" firestore:"start_time"`
	EndTime   int64  `json:"end_time" firestore:"end_time"`
	IsEnabled bool   `json:"is_enabled" firestore:"is_enabled"`
}

// Service provides coupon CRUD operations.
type Service interface {
	CreateCoupon(ctx context.Context, coupon Coupon) (Coupon, error)
	GetCoupons(ctx context.Context, page, pageSize int, search string) ([]Coupon, int, error)
	GetCoupon(ctx context.Context, id string) (Coupon, error)
	UpdateCoupon(ctx context.Context, id string, coupon Coupon) (Coupon, error)
	DeleteCoupon(ctx context.Context, id string) error
}
