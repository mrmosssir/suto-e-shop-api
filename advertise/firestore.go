package advertise

import (
	"context"
	"log"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

// FirestoreService is a Firestore implementation of the advertise service.
type FirestoreService struct {
	client     *firestore.Client
	collection string
}

// NewFirestoreService creates a new Firestore-backed advertise service.
func NewFirestoreService(client *firestore.Client) *FirestoreService {
	return &FirestoreService{
		client:     client,
		collection: "advertises",
	}
}

func (s *FirestoreService) AdminCreateAdvertise(ctx context.Context, advertise Advertise) (Advertise, error) {
	ref := s.client.Collection(s.collection).NewDoc()
	advertise.ID = ref.ID

	_, err := ref.Set(ctx, advertise)
	if err != nil {
		log.Printf("Failed to create advertise: %v", err)
		return Advertise{}, err
	}
	return advertise, nil
}

func (s *FirestoreService) AdminGetAdvertises(ctx context.Context, page, pageSize int, search string) ([]Advertise, int, error) {
	var advertises []Advertise
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
			log.Printf("Failed to get advertises: %v", err)
			return nil, 0, err
		}
		var advertise Advertise
		doc.DataTo(&advertise)
		advertises = append(advertises, advertise)
	}

	totalCount := len(advertises)
	start := (page - 1) * pageSize
	end := start + pageSize

	if start > totalCount {
		return []Advertise{}, totalCount, nil
	}

	if end > totalCount {
		end = totalCount
	}

	return advertises[start:end], totalCount, nil
}

func (s *FirestoreService) AdminGetAdvertise(ctx context.Context, id string) (Advertise, error) {
	doc, err := s.client.Collection(s.collection).Doc(id).Get(ctx)
	if err != nil {
		return Advertise{}, err
	}

	var advertise Advertise
	if err := doc.DataTo(&advertise); err != nil {
		return Advertise{}, err
	}
	return advertise, nil
}

func (s *FirestoreService) AdminUpdateAdvertise(ctx context.Context, id string, advertise Advertise) (Advertise, error) {
	advertise.ID = id
	_, err := s.client.Collection(s.collection).Doc(id).Set(ctx, advertise)
	if err != nil {
		return Advertise{}, err
	}
	return advertise, nil
}

func (s *FirestoreService) AdminDeleteAdvertise(ctx context.Context, id string) error {
	_, err := s.client.Collection(s.collection).Doc(id).Delete(ctx)
	return err
}

func (s *FirestoreService) GetAdvertises(ctx context.Context) ([]ClientAdvertise, error) {
	var advertises []ClientAdvertise
	query := s.client.Collection(s.collection).Where("is_enabled", "==", true)
	iter := query.Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Printf("Failed to get advertises: %v", err)
			return nil, err
		}
		var advertise ClientAdvertise
		doc.DataTo(&advertise)
		advertises = append(advertises, advertise)
	}
	return advertises, nil
}
