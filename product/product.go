package product

import (
	"context"
	"fmt"
	"sort"
)

// 後台列表用
type Product struct {
	ID          string  `json:"id" firestore:"id"`
	Name        string  `json:"name" firestore:"name"`
	Category    string  `json:"category" firestore:"category"`
	Price       int32   `json:"price" firestore:"price"`
	OriginPrice int32   `json:"origin_price" firestore:"origin_price"`
	Unit        string  `json:"unit" firestore:"unit"`
	Description string  `json:"description" firestore:"description"`
	Content     string  `json:"content" firestore:"content"`
	IsEnabled   bool    `json:"is_enabled" firestore:"is_enabled"`
	ImageURL    string  `json:"image_url" firestore:"image_url"`
	Rating			float32 `json:"rating" firestore:"rating"`
}

// 給前台列表顯示用
type ProductSimple struct {
	ID          string  `json:"id" firestore:"id"`
	Category    string  `json:"category" firestore:"category"`
	Name        string  `json:"name" firestore:"name"`
	Price       int32   `json:"price" firestore:"price"`
	OriginPrice int32   `json:"origin_price" firestore:"origin_price"`
	ImageURL    string  `json:"image_url" firestore:"image_url"`
	Rating      float32 `json:"rating" firestore:"rating"`
}

// Service provides product CRUD operations.
type Service interface {
	AdminCreateProduct(ctx context.Context, product Product) (Product, error)
	AdminGetProducts(ctx context.Context, page, pageSize int, search string) ([]Product, int, error)
	AdminGetProduct(ctx context.Context, id string) (Product, error)
	AdminUpdateProduct(ctx context.Context, id string, product Product) (Product, error)
	AdminDeleteProduct(ctx context.Context, id string) error
	GetProducts(ctx context.Context, page, pageSize int, search string) ([]ProductSimple, int, error)
	GetProduct(ctx context.Context, id string) (Product, error)
}

// InMemoryService is an in-memory implementation of the product service.
type InMemoryService struct {
	products      map[string]Product
	nextProductID int
}

// NewInMemoryService creates a new in-memory product service.
func NewInMemoryService() *InMemoryService {
	return &InMemoryService{
		products:      make(map[string]Product),
		nextProductID: 1,
	}
}

func (s *InMemoryService) GetProducts(ctx context.Context, page, pageSize int, search string) ([]ProductSimple, int, error) {
	var productList []ProductSimple
	for _, p := range s.products {
		productList = append(productList, ProductSimple{
			ID:    p.ID,
			Name:  p.Name,
			Price: p.Price,
			OriginPrice: p.OriginPrice,
			ImageURL: p.ImageURL,
			Rating: p.Rating,
		})
	}

	// Sort by ID for consistent pagination
	sort.Slice(productList, func(i, j int) bool {
		return productList[i].ID < productList[j].ID
	})

	totalCount := len(productList)

	start := (page - 1) * pageSize
	end := start + pageSize

	if start > totalCount {
		return []ProductSimple{}, 0, nil
	}

	if end > totalCount {
		end = totalCount
	}

	return productList[start:end], totalCount, nil
}

func (s *InMemoryService) GetProduct(ctx context.Context, id string) (Product, error) {
	product, ok := s.products[id]
	if !ok {
		return Product{}, fmt.Errorf("product not found")
	}
	return product, nil
}

func (s *InMemoryService) AdminCreateProduct(ctx context.Context, product Product) (Product, error) {
	product.ID = fmt.Sprintf("%d", s.nextProductID)
	s.nextProductID++
	s.products[product.ID] = product
	return product, nil
}

func (s *InMemoryService) AdminGetProducts(ctx context.Context, page, pageSize int, search string) ([]Product, int, error) {
	var productList []Product
	for _, p := range s.products {
		productList = append(productList, p)
	}

	// Sort by ID for consistent pagination
	sort.Slice(productList, func(i, j int) bool {
		return productList[i].ID < productList[j].ID
	})

	totalCount := len(productList)

	start := (page - 1) * pageSize
	end := start + pageSize

	if start > totalCount {
		return []Product{}, totalCount, nil
	}

	if end > totalCount {
		end = totalCount
	}

	return productList[start:end], totalCount, nil
}

func (s *InMemoryService) AdminGetProduct(ctx context.Context, id string) (Product, error) {
	product, ok := s.products[id]
	if !ok {
		return Product{}, fmt.Errorf("product not found")
	}
	return product, nil
}

func (s *InMemoryService) AdminUpdateProduct(ctx context.Context, id string, product Product) (Product, error) {
	if _, ok := s.products[id]; !ok {
		return Product{}, fmt.Errorf("product not found")
	}
	product.ID = id
	s.products[id] = product
	return product, nil
}

func (s *InMemoryService) AdminDeleteProduct(ctx context.Context, id string) error {
	if _, ok := s.products[id]; !ok {
		return fmt.Errorf("product not found")
	}
	delete(s.products, id)
	return nil
}