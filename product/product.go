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
	Rating      float32 `json:"rating" firestore:"rating"`
	IsNew       bool    `json:"is_new" firestore:"is_new"`
	IsHot       bool    `json:"is_hot" firestore:"is_hot"`
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
	GetProductsIds(ctx context.Context, ids []string) ([]Product, error)
	GetProduct(ctx context.Context, id string) (Product, error)
	GetNewProducts(ctx context.Context) ([]ProductSimple, error)
	GetHotProducts(ctx context.Context) ([]ProductSimple, error)
	CountNewProducts(ctx context.Context) (int, error)
	CountHotProducts(ctx context.Context) (int, error)
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

func (s *InMemoryService) GetProductsIds(ctx context.Context, ids []string) ([]Product, error) {
	var productList []Product
	for _, id := range ids {
		product, ok := s.products[id]
		if !ok {
			return nil, fmt.Errorf("product not found")
		}
		productList = append(productList, product)
	}
	return productList, nil
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

func (s *InMemoryService) GetNewProducts(ctx context.Context) ([]ProductSimple, error) {
	var products []ProductSimple
	for _, p := range s.products {
		if p.IsNew && p.IsEnabled {
			products = append(products, ProductSimple{
				ID:          p.ID,
				Category:    p.Category,
				Name:        p.Name,
				Price:       p.Price,
				OriginPrice: p.OriginPrice,
				ImageURL:    p.ImageURL,
				Rating:      p.Rating,
			})
		}
	}
	return products, nil
}

func (s *InMemoryService) GetHotProducts(ctx context.Context) ([]ProductSimple, error) {
	var products []ProductSimple
	for _, p := range s.products {
		if p.IsHot && p.IsEnabled {
			products = append(products, ProductSimple{
				ID:          p.ID,
				Category:    p.Category,
				Name:        p.Name,
				Price:       p.Price,
				OriginPrice: p.OriginPrice,
				ImageURL:    p.ImageURL,
				Rating:      p.Rating,
			})
		}
	}
	return products, nil
}

func (s *InMemoryService) CountNewProducts(ctx context.Context) (int, error) {
	count := 0
	for _, p := range s.products {
		if p.IsNew {
			count++
		}
	}
	return count, nil
}

func (s *InMemoryService) CountHotProducts(ctx context.Context) (int, error) {
	count := 0
	for _, p := range s.products {
		if p.IsHot {
			count++
		}
	}
	return count, nil
}