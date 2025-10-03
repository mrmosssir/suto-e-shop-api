package order

import (
	"context"
	"log"
	"strings"

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
