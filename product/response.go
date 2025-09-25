package product

import (
	"encoding/json"
	"net/http"

	"suto-e-shop-api/pkg/pagination"
)

// Response is the standardized API response format.
type Response struct {
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message"`
	Code    int         `json:"code"`
}

// PaginatedResponse is the standardized API response format for paginated data.
type PaginatedResponse struct {
	Data       interface{}           `json:"data,omitempty"`
	Pagination *pagination.Pagination `json:"pagination,omitempty"`
	Message    string                `json:"message"`
	Code       int                   `json:"code"`
}

// RespondWithError writes an error response.
func RespondWithError(w http.ResponseWriter, code int, message string) {
	RespondWithJSON(w, code, Response{Message: message, Code: code})
}

// RespondWithJSON writes a JSON response.
func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}