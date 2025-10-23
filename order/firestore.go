package order

import (
	"context"
	"log"
	"strconv"
	"strings"
	"time"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

// FirestoreService is a Firestore implementation of the order service.
type FirestoreService struct {
	client     *firestore.Client
	collection string
}

// NewFirestoreService creates a new Firestore-backed order service.
func NewFirestoreService(client *firestore.Client) *FirestoreService {
	return &FirestoreService{
		client:     client,
		collection: "orders",
	}
}

func (s *FirestoreService) GetOrders(ctx context.Context, page, pageSize int, search string) ([]Order, int, error) {
	var orders []Order

	query := s.client.Collection(s.collection).Query

	iter := query.Documents(ctx)
	defer iter.Stop()

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Printf("Failed to get orders: %v", err)
			return nil, 0, err
		}
		var order Order
		doc.DataTo(&order)

		if search != "" {
			if strings.Contains(order.Name, search) || strings.Contains(order.Mail, search) {
				orders = append(orders, order)
			}
		} else {
			orders = append(orders, order)
		}
	}

	totalCount := len(orders)
	start := (page - 1) * pageSize
	end := start + pageSize

	if start > totalCount {
		return []Order{}, totalCount, nil
	}

	if end > totalCount {
		end = totalCount
	}

	return orders[start:end], totalCount, nil
}

func (s *FirestoreService) UpdateOrder(ctx context.Context, id string, data map[string]interface{}) (Order, error) {
	docRef := s.client.Collection(s.collection).Doc(id)

	// Get the original document
	doc, err := docRef.Get(ctx)
	if err != nil {
		log.Printf("Failed to get order for update: %v", err)
		return Order{}, err
	}
	var originalOrder Order
	doc.DataTo(&originalOrder)

	// Prepare updates
	var updates []firestore.Update
	for key, value := range data {
		updates = append(updates, firestore.Update{Path: key, Value: value})

		// Check for is_paid update
		if key == "is_paid" {
			isPaid, ok := value.(bool)
			if ok && isPaid && !originalOrder.IsPaid {
				updates = append(updates, firestore.Update{Path: "paid_at", Value: strconv.FormatInt(time.Now().Unix(), 10)})
			}
		}

		// Check for is_picked update
		if key == "is_picked" {
			isPicked, ok := value.(bool)
			if ok && isPicked && !originalOrder.IsPicked {
				updates = append(updates, firestore.Update{Path: "picked_at", Value: strconv.FormatInt(time.Now().Unix(), 10)})
			}
		}

		if (key == "is_enabled") {
			isEnabled, ok := value.(bool)
			if ok && !isEnabled {
				updates = append(updates, firestore.Update{Path: "disabled_at", Value: strconv.FormatInt(time.Now().Unix(), 10)})
			}
		}
	}

	_, err = docRef.Update(ctx, updates)
	if err != nil {
		log.Printf("Failed to update order: %v", err)
		return Order{}, err
	}

	// Get the updated document
	updatedDoc, err := docRef.Get(ctx)
	if err != nil {
		log.Printf("Failed to get updated order: %v", err)
		return Order{}, err
	}

	var updatedOrder Order
	updatedDoc.DataTo(&updatedOrder)
	return updatedOrder, nil
}

func (s *FirestoreService) CreateOrder(ctx context.Context, req CreateOrderRequest) (Order, error) {
	ref := s.client.Collection(s.collection).NewDoc()

	totalPrice := 0
	for _, p := range req.Products {
		totalPrice += p.Price * p.Count
	}

	order := Order{
		ID:         ref.ID,
		Name:       req.Name,
		Mail:       req.Mail,
		Products:   req.Products,
		TotalPrice: totalPrice,
		IsEnabled:  true,
		CreatedAt:  strconv.FormatInt(time.Now().Unix(), 10),
	}

	_, err := ref.Set(ctx, order)
	if err != nil {
		log.Printf("Failed to create order: %v", err)
		return Order{}, err
	}
	return order, nil
}
