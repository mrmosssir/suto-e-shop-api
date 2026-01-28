package upload

import (
	"encoding/json"
	"net/http"
)

// Response is a standard JSON response.
type Response struct {
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message"`
	Code    int         `json:"code"`
}

// RespondWithError sends an error response.
func RespondWithError(w http.ResponseWriter, code int, message string) {
	RespondWithJSON(w, code, Response{Message: message, Code: code})
}

// RespondWithJSON sends a JSON response.
func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(payload)
}
