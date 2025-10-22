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

// RegisterAdminRoutes registers the product routes to the router.
func (h *Handler) RegisterAdminRoutes(router *mux.Router) {

	adminRouter := router.PathPrefix("/product").Subrouter()

	adminRouter.HandleFunc("", h.AdminCreateProduct).Methods("POST")
	adminRouter.HandleFunc("", h.AdminGetProducts).Methods("GET")
	adminRouter.HandleFunc("/{id}", h.AdminGetProduct).Methods("GET")
	adminRouter.HandleFunc("/{id}", h.AdminUpdateProduct).Methods("PUT")
	adminRouter.HandleFunc("/{id}", h.AdminDeleteProduct).Methods("DELETE")
}

// RegisterClientRoutes registers the product routes to the router.
func (h *Handler) RegisterClientRoutes(router *mux.Router) {
	router.HandleFunc("/products", h.GetProducts).Methods("GET")
	router.HandleFunc("/product/{id}", h.GetProduct).Methods("GET")
}

func (h *Handler) AdminCreateProduct(w http.ResponseWriter, r *http.Request) {
	var product Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	createdProduct, err := h.service.AdminCreateProduct(r.Context(), product)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	RespondWithJSON(w, http.StatusCreated, Response{Data: createdProduct, Message: "success", Code: 0})
}

func (h *Handler) AdminGetProducts(w http.ResponseWriter, r *http.Request) {
	page, pageSize := pagination.GetPaginationParams(r)
	search := r.URL.Query().Get("search")

	products, totalCount, err := h.service.AdminGetProducts(r.Context(), page, pageSize, search)
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

func (h *Handler) AdminGetProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	product, err := h.service.AdminGetProduct(r.Context(), id)
	if err != nil {
		RespondWithError(w, http.StatusNotFound, "Product not found")
		return
	}

	RespondWithJSON(w, http.StatusOK, Response{Data: product, Message: "success", Code: 0})
}

func (h *Handler) AdminUpdateProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var product Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	updatedProduct, err := h.service.AdminUpdateProduct(r.Context(), id, product)
	if err != nil {
		RespondWithError(w, http.StatusNotFound, "Product not found")
		return
	}

	RespondWithJSON(w, http.StatusOK, Response{Data: updatedProduct, Message: "success", Code: 0})
}

func (h *Handler) AdminDeleteProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if err := h.service.AdminDeleteProduct(r.Context(), id); err != nil {
		RespondWithError(w, http.StatusNotFound, "Product not found")
		return
	}

	RespondWithJSON(w, http.StatusOK, Response{Message: "success", Code: 0})
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