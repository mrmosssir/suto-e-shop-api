package banner

import "context"

// Banner defines the structure for a banner.
type Banner struct {
	ID        string `json:"id" firestore:"id"`
	Name      string `json:"name" firestore:"name"`
	Image     string `json:"image" firestore:"image"`
	IsEnabled bool   `json:"is_enabled" firestore:"is_enabled"`
}

// ClientBanner is for client API responses (without IsEnabled field)
type ClientBanner struct {
	ID    string `json:"id" firestore:"id"`
	Name  string `json:"name" firestore:"name"`
	Image string `json:"image" firestore:"image"`
}

// Service provides banner operations.
type Service interface {
	// Admin operations
	AdminCreateBanner(ctx context.Context, banner Banner) (Banner, error)
	AdminGetBanners(ctx context.Context, page, pageSize int, search string) ([]Banner, int, error)
	AdminGetBanner(ctx context.Context, id string) (Banner, error)
	AdminUpdateBanner(ctx context.Context, id string, banner Banner) (Banner, error)
	AdminDeleteBanner(ctx context.Context, id string) error

	// Client operations
	GetBanners(ctx context.Context) ([]ClientBanner, error)
}
