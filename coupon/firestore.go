package coupon

import (
	"context"
	"log"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

// FirestoreService is a Firestore implementation of the coupon service.
type FirestoreService struct {
	client     *firestore.Client
	collection string
}

// NewFirestoreService creates a new Firestore-backed coupon service.
func NewFirestoreService(client *firestore.Client) *FirestoreService {
	return &FirestoreService{
		client:     client,
		collection: "coupons",
	}
}

func (s *FirestoreService) CreateCoupon(ctx context.Context, coupon Coupon) (Coupon, error) {
	ref := s.client.Collection(s.collection).NewDoc()
	coupon.ID = ref.ID
	_, err := ref.Set(ctx, coupon)
	if err != nil {
		log.Printf("Failed to create coupon: %v", err)
		return Coupon{}, err
	}
	return coupon, nil
}

func (s *FirestoreService) GetCoupons(ctx context.Context, page, pageSize int, search string) ([]Coupon, int, error) {
	var coupons []Coupon
	// For more advanced search capabilities, consider using a dedicated search service like Algolia or Elasticsearch.
	query := s.client.Collection(s.collection).Query
	if search != "" {
		query = query.Where("name", ">=", search).Where("name", "<=", search+"\uf8ff")
	}

	iter := query.Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Printf("Failed to get coupons: %v", err)
			return nil, 0, err
		}
		var coupon Coupon
		doc.DataTo(&coupon)
		coupons = append(coupons, coupon)
	}

	totalCount := len(coupons)
	start := (page - 1) * pageSize
	end := start + pageSize

	if start > totalCount {
		return []Coupon{}, totalCount, nil
	}

	if end > totalCount {
		end = totalCount
	}

	return coupons[start:end], totalCount, nil
}

func (s *FirestoreService) GetCoupon(ctx context.Context, id string) (Coupon, error) {
	doc, err := s.client.Collection(s.collection).Doc(id).Get(ctx)
	if err != nil {
		log.Printf("Failed to get coupon: %v", err)
		return Coupon{}, err
	}
	var coupon Coupon
	doc.DataTo(&coupon)
	return coupon, nil
}

func (s *FirestoreService) UpdateCoupon(ctx context.Context, id string, coupon Coupon) (Coupon, error) {
	_, err := s.client.Collection(s.collection).Doc(id).Set(ctx, coupon)
	if err != nil {
		log.Printf("Failed to update coupon: %v", err)
		return Coupon{}, err
	}
	coupon.ID = id
	return coupon, nil
}

func (s *FirestoreService) DeleteCoupon(ctx context.Context, id string) error {
	_, err := s.client.Collection(s.collection).Doc(id).Delete(ctx)
	if err != nil {
		log.Printf("Failed to delete coupon: %v", err)
		return err
	}
	return nil
}
