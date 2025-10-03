package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go/v4"
	"github.com/gorilla/mux"
	"suto-e-shop-api/auth"
	"suto-e-shop-api/coupon"
	"suto-e-shop-api/product"
)

// CORSMiddleware sets up the CORS headers for every request.
func CORSMiddleware(allowedOrigins []string) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")
			isAllowed := false

			// Check if the origin is in our allowed list
			for _, allowed := range allowedOrigins {
				if allowed == "*" || origin == allowed {
					isAllowed = true
					break
				}
			}

			if isAllowed {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Set("Access-Control-Allow-Credentials", "true")
			}

			w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
			w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, X-Auth-Token")

			// Continue to the next handler
			next.ServeHTTP(w, r)
		})
	}
}

func main() {
	ctx := context.Background()
	projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")
	if projectID == "" {
		log.Fatal("GOOGLE_CLOUD_PROJECT environment variable must be set.")
	}

	databaseID := os.Getenv("FIRESTORE_DATABASE_ID")
	if databaseID == "" {
		log.Fatal("FIRESTORE_DATABASE_ID environment variable must be set.")
	}

	// Initialize Firestore Client
	client, err := firestore.NewClientWithDatabase(ctx, projectID, databaseID)
	if err != nil {
		log.Fatalf("Failed to create Firestore client: %v", err)
	}
	defer client.Close()

	// Initialize Firebase App
	fbApp, err := firebase.NewApp(ctx, nil)
	if err != nil {
		log.Fatalf("Failed to create Firebase app: %v", err)
	}

	r := mux.NewRouter()
	r.StrictSlash(true)

	// Setup CORS
	allowedOrigins := []string{"http://localhost:5173"}
	r.Use(CORSMiddleware(allowedOrigins))

	// Add a handler for OPTIONS requests to handle preflight CORS requests.
	r.Methods("OPTIONS").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})

	// Create a subrouter for the admin routes that require authentication
	adminRouter := r.PathPrefix("/admin").Subrouter()
	adminRouter.Use(auth.FirebaseJWTMiddleware(fbApp))

	// Product routes
	productService := product.NewFirestoreService(client)
	productHandler := product.NewHandler(productService)
	productHandler.RegisterRoutes(adminRouter)

	// Coupon routes
	couponService := coupon.NewFirestoreService(client)
	couponHandler := coupon.NewHandler(couponService)
	couponHandler.RegisterRoutes(adminRouter)


	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("Server listening on port", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
