package category

import "context"

// Category defines the structure for a category.
type Category struct {
	ID        string `json:"id" firestore:"id"`
	Name      string `json:"name" firestore:"name"`
	IsEnabled bool   `json:"is_enabled" firestore:"is_enabled"`
}

// Service provides category operations.
type Service interface {
	// Admin operations
	AdminCreateCategory(ctx context.Context, category Category) (Category, error)
	AdminGetCategories(ctx context.Context, page, pageSize int, search string) ([]Category, int, error)
	AdminGetCategory(ctx context.Context, id string) (Category, error)
	AdminUpdateCategory(ctx context.Context, id string, category Category) (Category, error)
	AdminDeleteCategory(ctx context.Context, id string) error

	// Client operations
	GetCategories(ctx context.Context) ([]Category, error)
}
