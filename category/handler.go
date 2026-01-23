package category

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

// Handler holds the category service.
type Handler struct {
	service Service
}

// NewHandler creates a new category handler.
func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

// RegisterAdminRoutes registers the admin category routes to the router.
func (h *Handler) RegisterAdminRoutes(router *mux.Router) {
	adminRouter := router.PathPrefix("/category").Subrouter()

	adminRouter.HandleFunc("", h.AdminCreateCategory).Methods("POST")
	adminRouter.HandleFunc("", h.AdminGetCategories).Methods("GET")
	adminRouter.HandleFunc("/{id}", h.AdminGetCategory).Methods("GET")
	adminRouter.HandleFunc("/{id}", h.AdminUpdateCategory).Methods("PUT")
	adminRouter.HandleFunc("/{id}", h.AdminDeleteCategory).Methods("DELETE")
}

// RegisterClientRoutes registers the client category routes to the router.
func (h *Handler) RegisterClientRoutes(router *mux.Router) {
	router.HandleFunc("/categories", h.GetCategories).Methods("GET")
}

func (h *Handler) AdminCreateCategory(w http.ResponseWriter, r *http.Request) {
	var category Category
	if err := json.NewDecoder(r.Body).Decode(&category); err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	createdCategory, err := h.service.AdminCreateCategory(r.Context(), category)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	RespondWithJSON(w, http.StatusCreated, Response{Data: createdCategory, Message: "success", Code: 0})
}

func (h *Handler) AdminGetCategories(w http.ResponseWriter, r *http.Request) {
	categories, err := h.service.AdminGetCategories(r.Context())
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	RespondWithJSON(w, http.StatusOK, Response{Data: categories, Message: "success", Code: 0})
}

func (h *Handler) AdminGetCategory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	category, err := h.service.AdminGetCategory(r.Context(), id)
	if err != nil {
		RespondWithError(w, http.StatusNotFound, "Category not found")
		return
	}

	RespondWithJSON(w, http.StatusOK, Response{Data: category, Message: "success", Code: 0})
}

func (h *Handler) AdminUpdateCategory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var category Category
	if err := json.NewDecoder(r.Body).Decode(&category); err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	updatedCategory, err := h.service.AdminUpdateCategory(r.Context(), id, category)
	if err != nil {
		RespondWithError(w, http.StatusNotFound, "Category not found")
		return
	}

	RespondWithJSON(w, http.StatusOK, Response{Data: updatedCategory, Message: "success", Code: 0})
}

func (h *Handler) AdminDeleteCategory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if err := h.service.AdminDeleteCategory(r.Context(), id); err != nil {
		RespondWithError(w, http.StatusNotFound, "Category not found")
		return
	}

	RespondWithJSON(w, http.StatusOK, Response{Message: "success", Code: 0})
}

func (h *Handler) GetCategories(w http.ResponseWriter, r *http.Request) {
	categories, err := h.service.GetCategories(r.Context())
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	RespondWithJSON(w, http.StatusOK, Response{Data: categories, Message: "success", Code: 0})
}
