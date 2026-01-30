package category

import (
	"context"
	"log"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)


// FirestoreService is a Firestore implementation of the category service.
type FirestoreService struct {
	client     *firestore.Client
	collection string
}

// NewFirestoreService creates a new Firestore-backed category service.
func NewFirestoreService(client *firestore.Client) *FirestoreService {
	return &FirestoreService{
		client:     client,
		collection: "category",
	}
}

func (s *FirestoreService) AdminCreateCategory(ctx context.Context, category Category) (Category, error) {
	ref := s.client.Collection(s.collection).NewDoc()
	category.ID = ref.ID

	_, err := ref.Set(ctx, category)
	if err != nil {
		log.Printf("Failed to create category: %v", err)
		return Category{}, err
	}
	return category, nil
}

func (s *FirestoreService) AdminGetCategories(ctx context.Context, page, pageSize int, search string) ([]Category, int, error) {
	var categories []Category
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
			log.Printf("Failed to get categories: %v", err)
			return nil, 0, err
		}
		var category Category
		doc.DataTo(&category)
		categories = append(categories, category)
	}

	totalCount := len(categories)
	start := (page - 1) * pageSize
	end := start + pageSize

	if start > totalCount {
		return []Category{}, totalCount, nil
	}

	if end > totalCount {
		end = totalCount
	}

	return categories[start:end], totalCount, nil
}

func (s *FirestoreService) AdminGetCategory(ctx context.Context, id string) (Category, error) {
	doc, err := s.client.Collection(s.collection).Doc(id).Get(ctx)
	if err != nil {
		return Category{}, err
	}

	var category Category
	if err := doc.DataTo(&category); err != nil {
		return Category{}, err
	}
	return category, nil
}

func (s *FirestoreService) AdminUpdateCategory(ctx context.Context, id string, category Category) (Category, error) {
	category.ID = id
	_, err := s.client.Collection(s.collection).Doc(id).Set(ctx, category)
	if err != nil {
		return Category{}, err
	}
	return category, nil
}

func (s *FirestoreService) AdminDeleteCategory(ctx context.Context, id string) error {
	_, err := s.client.Collection(s.collection).Doc(id).Delete(ctx)
	return err
}

func (s *FirestoreService) GetCategories(ctx context.Context) ([]ClientCategory, error) {
	var categories []ClientCategory
	iter := s.client.Collection(s.collection).Where("is_enabled", "==", true).Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}

		var category ClientCategory
		if err := doc.DataTo(&category); err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}
	return categories, nil
}
