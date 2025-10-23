package product

import (
	"context"
	"log"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

// FirestoreService is a Firestore implementation of the product service.
type FirestoreService struct {
	client     *firestore.Client
	collection string
}

// NewFirestoreService creates a new Firestore-backed product service.
func NewFirestoreService(client *firestore.Client) *FirestoreService {
	return &FirestoreService{
		client:     client,
		collection: "products",
	}
}

func (s *FirestoreService) AdminCreateProduct(ctx context.Context, product Product) (Product, error) {
	ref := s.client.Collection(s.collection).NewDoc()
	product.ID = ref.ID
	_, err := ref.Set(ctx, product)
	if err != nil {
		log.Printf("Failed to create product: %v", err)
		return Product{}, err
	}
	return product, nil
}

func (s *FirestoreService) AdminGetProducts(ctx context.Context, page, pageSize int, search string) ([]Product, int, error) {
	var products []Product
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
			log.Printf("Failed to get products: %v", err)
			return nil, 0, err
		}
		var product Product
		doc.DataTo(&product)
		products = append(products, product)
	}

	totalCount := len(products)
	start := (page - 1) * pageSize
	end := start + pageSize

	if start > totalCount {
		return []Product{}, totalCount, nil
	}

	if end > totalCount {
		end = totalCount
	}

	return products[start:end], totalCount, nil
}

func (s *FirestoreService) AdminGetProduct(ctx context.Context, id string) (Product, error) {
	doc, err := s.client.Collection(s.collection).Doc(id).Get(ctx)
	if err != nil {
		log.Printf("Failed to get product: %v", err)
		return Product{}, err
	}
	var product Product
	doc.DataTo(&product)
	return product, nil
}

func (s *FirestoreService) AdminUpdateProduct(ctx context.Context, id string, product Product) (Product, error) {
	_, err := s.client.Collection(s.collection).Doc(id).Set(ctx, product)
	if err != nil {
		log.Printf("Failed to update product: %v", err)
		return Product{}, err
	}
	product.ID = id
	return product, nil
}

func (s *FirestoreService) AdminDeleteProduct(ctx context.Context, id string) error {
	_, err := s.client.Collection(s.collection).Doc(id).Delete(ctx)
	if err != nil {
		log.Printf("Failed to delete product: %v", err)
		return err
	}
	return nil
}

func (s *FirestoreService) GetProducts(ctx context.Context, page, pageSize int, search string) ([]ProductSimple, int, error) {
	var products []ProductSimple
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
			log.Printf("Failed to get products: %v", err)
			return nil, 0, err
		}
		var product ProductSimple
		doc.DataTo(&product)
		products = append(products, product)
	}

	totalCount := len(products)
	start := (page - 1) * pageSize
	end := start + pageSize

	if start > totalCount {
		return []ProductSimple{}, totalCount, nil
	}

	if end > totalCount {
		end = totalCount
	}

	return products[start:end], totalCount, nil
}

func (s *FirestoreService) GetProductsIds(ctx context.Context, ids []string) ([]Product, error) {
	var products []Product
	for _, id := range ids {
		product, err := s.GetProduct(ctx, id)
		if err != nil {
			log.Printf("Failed to get product: %v", err)
			continue
		}
		products = append(products, product)
	}
	return products, nil
}

func (s *FirestoreService) GetProduct(ctx context.Context, id string) (Product, error) {
	doc, err := s.client.Collection(s.collection).Doc(id).Get(ctx)
	if err != nil {
		log.Printf("Failed to get product: %v", err)
		return Product{}, err
	}
	var product Product
	doc.DataTo(&product)
	return product, nil
}