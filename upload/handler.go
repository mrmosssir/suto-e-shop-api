package upload

import (
	"io"
	"net/http"

	"github.com/gorilla/mux"
)

// Handler holds the upload service.
type Handler struct {
	service Service
}

// NewHandler creates a new upload handler.
func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

// RegisterAdminRoutes registers the admin upload routes to the router.
func (h *Handler) RegisterAdminRoutes(router *mux.Router) {
	router.HandleFunc("/upload", h.UploadImage).Methods("POST")
}

// UploadImage handles the image upload request.
func (h *Handler) UploadImage(w http.ResponseWriter, r *http.Request) {
	// Parse multipart form with max 10MB
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		RespondWithError(w, http.StatusBadRequest, "Failed to parse form: file too large or invalid format")
		return
	}

	// Get the type parameter
	uploadType := r.FormValue("type")
	if uploadType == "" {
		RespondWithError(w, http.StatusBadRequest, "Missing required field: type")
		return
	}

	// Get the file from the form
	file, header, err := r.FormFile("file")
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Missing required field: file")
		return
	}
	defer file.Close()

	// Validate content type
	contentType := header.Header.Get("Content-Type")
	if !isValidImageType(contentType) {
		RespondWithError(w, http.StatusBadRequest, "Invalid file type. Allowed types: jpeg, png, gif, webp, svg")
		return
	}

	// Read file data
	fileData, err := io.ReadAll(file)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Failed to read file")
		return
	}

	// Upload the image
	result, err := h.service.UploadImage(r.Context(), fileData, contentType, uploadType)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	RespondWithJSON(w, http.StatusCreated, Response{
		Data:    result,
		Message: "success",
		Code:    0,
	})
}

// isValidImageType checks if the content type is a valid image type
func isValidImageType(contentType string) bool {
	validTypes := map[string]bool{
		"image/jpeg":    true,
		"image/png":     true,
		"image/gif":     true,
		"image/webp":    true,
		"image/svg+xml": true,
	}
	return validTypes[contentType]
}
