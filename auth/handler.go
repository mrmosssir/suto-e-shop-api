package auth

import (
	"encoding/json"
	"net/http"
)

// LogoutHandler handles user logout.
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	// For stateless JWT, logout is typically handled on the client-side by
	// deleting the token. This endpoint is provided for completeness but
	// doesn't perform any server-side action like invalidating the token.
	// A more complex implementation could involve a token blocklist.
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Logout successful"})
}
