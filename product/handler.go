package product

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"suto-e-shop-api/pkg/pagination"
)

// Handler holds the product service.
type Handler struct {
	service Service
}

// NewHandler creates a new product handler.
func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

// RegisterRoutes registers the product routes to the router.
func (h *Handler) RegisterRoutes(router *mux.Router) {
	adminRouter := router.PathPrefix("/product").Subrouter()

	adminRouter.HandleFunc("", h.CreateProduct).Methods("POST")
	adminRouter.HandleFunc("", h.GetProducts).Methods("GET")
	adminRouter.HandleFunc("/{id}", h.GetProduct).Methods("GET")
	adminRouter.HandleFunc("/{id}", h.UpdateProduct).Methods("PUT")
	adminRouter.HandleFunc("/{id}", h.DeleteProduct).Methods("DELETE")
}

func (h *Handler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var product Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	createdProduct, err := h.service.CreateProduct(r.Context(), product)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	RespondWithJSON(w, http.StatusCreated, Response{Data: createdProduct, Message: "success", Code: 0})
}

func (h *Handler) GetProducts(w http.ResponseWriter, r *http.Request) {
	page, pageSize := pagination.GetPaginationParams(r)
	search := r.URL.Query().Get("search")

	products, totalCount, err := h.service.GetProducts(r.Context(), page, pageSize, search)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	paginator := pagination.New(page, pageSize, totalCount)

	RespondWithJSON(w, http.StatusOK, PaginatedResponse{
		Data:       products,
		Pagination: paginator,
		Message:    "success",
		Code:       0,
	})
}

func (h *Handler) GetProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	product, err := h.service.GetProduct(r.Context(), id)
	if err != nil {
		RespondWithError(w, http.StatusNotFound, "Product not found")
		return
	}

	RespondWithJSON(w, http.StatusOK, Response{Data: product, Message: "success", Code: 0})
}

func (h *Handler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var product Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	updatedProduct, err := h.service.UpdateProduct(r.Context(), id, product)
	if err != nil {
		RespondWithError(w, http.StatusNotFound, "Product not found")
		return
	}

	RespondWithJSON(w, http.StatusOK, Response{Data: updatedProduct, Message: "success", Code: 0})
}

func (h *Handler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if err := h.service.DeleteProduct(r.Context(), id); err != nil {
		RespondWithError(w, http.StatusNotFound, "Product not found")
		return
	}

	RespondWithJSON(w, http.StatusOK, Response{Message: "success", Code: 0})
}