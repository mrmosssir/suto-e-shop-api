package auth

import (
	"context"
	"net/http"
	"strings"

	firebase "firebase.google.com/go/v4"
)

// FirebaseJWTMiddleware is a middleware function to protect routes using Firebase ID tokens.
func FirebaseJWTMiddleware(fbApp *firebase.App) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// A helper function to send CORS-compliant errors
			sendError := func(message string, code int) {
				// The CORSMiddleware should have already set this, but we ensure it for error responses.
				w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
				http.Error(w, message, code)
			}

			authHeader := r.Header.Get("X-Auth-Token")
			if authHeader == "" {
				sendError("X-Auth-Token header required", http.StatusUnauthorized)
				return
			}

			tokenString := strings.TrimPrefix(authHeader, "Bearer ")
			if tokenString == authHeader {
				sendError("Could not find bearer token in X-Auth-Token header", http.StatusUnauthorized)
				return
			}

			client, err := fbApp.Auth(context.Background())
			if err != nil {
				sendError("Failed to create Firebase auth client", http.StatusInternalServerError)
				return
			}

			token, err := client.VerifyIDToken(context.Background(), tokenString)
			if err != nil {
				sendError("Invalid Firebase ID token", http.StatusUnauthorized)
				return
			}

			// You can add the decoded token to the request context if needed by other handlers
			ctx := context.WithValue(r.Context(), "firebaseToken", token)
			r = r.WithContext(ctx)

			// Token is valid, proceed to the next handler
			next.ServeHTTP(w, r)
		})
	}
}
