package banner

import (
	"context"
	"log"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

// FirestoreService is a Firestore implementation of the banner service.
type FirestoreService struct {
	client     *firestore.Client
	collection string
}

// NewFirestoreService creates a new Firestore-backed banner service.
func NewFirestoreService(client *firestore.Client) *FirestoreService {
	return &FirestoreService{
		client:     client,
		collection: "banners",
	}
}

func (s *FirestoreService) AdminCreateBanner(ctx context.Context, banner Banner) (Banner, error) {
	ref := s.client.Collection(s.collection).NewDoc()
	banner.ID = ref.ID

	_, err := ref.Set(ctx, banner)
	if err != nil {
		log.Printf("Failed to create banner: %v", err)
		return Banner{}, err
	}
	return banner, nil
}

func (s *FirestoreService) AdminGetBanners(ctx context.Context, page, pageSize int, search string) ([]Banner, int, error) {
	var banners []Banner
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
			log.Printf("Failed to get banners: %v", err)
			return nil, 0, err
		}
		var banner Banner
		doc.DataTo(&banner)
		banners = append(banners, banner)
	}

	totalCount := len(banners)
	start := (page - 1) * pageSize
	end := start + pageSize

	if start > totalCount {
		return []Banner{}, totalCount, nil
	}

	if end > totalCount {
		end = totalCount
	}

	return banners[start:end], totalCount, nil
}

func (s *FirestoreService) AdminGetBanner(ctx context.Context, id string) (Banner, error) {
	doc, err := s.client.Collection(s.collection).Doc(id).Get(ctx)
	if err != nil {
		return Banner{}, err
	}

	var banner Banner
	if err := doc.DataTo(&banner); err != nil {
		return Banner{}, err
	}
	return banner, nil
}

func (s *FirestoreService) AdminUpdateBanner(ctx context.Context, id string, banner Banner) (Banner, error) {
	banner.ID = id
	_, err := s.client.Collection(s.collection).Doc(id).Set(ctx, banner)
	if err != nil {
		return Banner{}, err
	}
	return banner, nil
}

func (s *FirestoreService) AdminDeleteBanner(ctx context.Context, id string) error {
	_, err := s.client.Collection(s.collection).Doc(id).Delete(ctx)
	return err
}

func (s *FirestoreService) GetBanners(ctx context.Context) ([]ClientBanner, error) {
	var banners []ClientBanner
	query := s.client.Collection(s.collection).Where("is_enabled", "==", true)
	iter := query.Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Printf("Failed to get banners: %v", err)
			return nil, err
		}
		var banner ClientBanner
		doc.DataTo(&banner)
		banners = append(banners, banner)
	}
	return banners, nil
}
