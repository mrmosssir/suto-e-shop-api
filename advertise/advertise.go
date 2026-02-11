package advertise

import "context"

// Advertise defines the structure for an advertise.
type Advertise struct {
	ID        string `json:"id" firestore:"id"`
	Name      string `json:"name" firestore:"name"`
	Image     string `json:"image" firestore:"image"`
	Link      string `json:"link,omitempty" firestore:"link,omitempty"`
	IsEnabled bool   `json:"is_enabled" firestore:"is_enabled"`
}

// ClientAdvertise is for client API responses (without IsEnabled field)
type ClientAdvertise struct {
	ID    string `json:"id" firestore:"id"`
	Name  string `json:"name" firestore:"name"`
	Image string `json:"image" firestore:"image"`
	Link  string `json:"link,omitempty" firestore:"link,omitempty"`
}

// Service provides advertise operations.
type Service interface {
	// Admin operations
	AdminCreateAdvertise(ctx context.Context, advertise Advertise) (Advertise, error)
	AdminGetAdvertises(ctx context.Context, page, pageSize int, search string) ([]Advertise, int, error)
	AdminGetAdvertise(ctx context.Context, id string) (Advertise, error)
	AdminUpdateAdvertise(ctx context.Context, id string, advertise Advertise) (Advertise, error)
	AdminDeleteAdvertise(ctx context.Context, id string) error

	// Client operations
	GetAdvertises(ctx context.Context) ([]ClientAdvertise, error)
}
