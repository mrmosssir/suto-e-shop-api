package order

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gorilla/mux"
	"suto-e-shop-api/pkg/pagination"
)

// Handler holds the order service.
type Handler struct {
	service Service
}

// NewHandler creates a new order handler.
func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

// RegisterAdminRoutes registers the admin order routes to the router.
func (h *Handler) RegisterAdminRoutes(router *mux.Router) {
	adminRouter := router.PathPrefix("/order").Subrouter()
	adminRouter.HandleFunc("", h.GetOrders).Methods("GET")
	adminRouter.HandleFunc("/{id}", h.UpdateOrder).Methods("PUT")
}

// RegisterClientRoutes registers the client order routes to the router.
func (h *Handler) RegisterClientRoutes(router *mux.Router) {
	router.HandleFunc("/order", h.SearchOrders).Methods("GET")
	router.HandleFunc("/order", h.CreateOrder).Methods("POST")
}

func (h *Handler) GetOrders(w http.ResponseWriter, r *http.Request) {
	page, pageSize := pagination.GetPaginationParams(r)
	search := r.URL.Query().Get("search")

	orders, totalCount, err := h.service.GetOrders(r.Context(), page, pageSize, search)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	paginator := pagination.New(page, pageSize, totalCount)

	RespondWithJSON(w, http.StatusOK, PaginatedResponse{
		Data:       orders,
		Pagination: paginator,
		Message:    "success",
		Code:       0,
	})
}

func (h *Handler) UpdateOrder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var data map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Validate the keys in the payload
	allowedKeys := map[string]bool{
		"is_enable": true,
		"is_picked": true,
		"is_paid":   true,
	}
	for key := range data {
		if !allowedKeys[key] {
			RespondWithError(w, http.StatusBadRequest, "Invalid field in request payload: "+key)
			return
		}
	}

	updatedOrder, err := h.service.UpdateOrder(r.Context(), id, data)
	if err != nil {
		RespondWithError(w, http.StatusNotFound, "Order not found")
		return
	}

	RespondWithJSON(w, http.StatusOK, Response{Data: updatedOrder, Message: "success", Code: 0})
}

func (h *Handler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var req CreateOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if err := validateCreateOrderRequest(req); err != nil {
		RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	order, err := h.service.CreateOrder(r.Context(), req)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	RespondWithJSON(w, http.StatusCreated, Response{Data: order, Message: "success", Code: 0})
}

func (h *Handler) SearchOrders(w http.ResponseWriter, r *http.Request) {
	search := r.URL.Query().Get("search")

	if search == "" {
		RespondWithError(w, http.StatusBadRequest, "search parameter is required")
		return
	}

	orders, err := h.service.SearchOrders(r.Context(), search)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	RespondWithJSON(w, http.StatusOK, Response{
		Data:    orders,
		Message: "success",
		Code:    0,
	})
}

func validateCreateOrderRequest(req CreateOrderRequest) error {
	if req.Mail == "" {
		return errors.New("mail is required")
	}
	if req.Name == "" {
		return errors.New("name is required")
	}
	if len(req.Products) == 0 {
		return errors.New("products are required")
	}
	for _, p := range req.Products {
		if p.Name == "" {
			return errors.New("product name is required")
		}
		if p.Count <= 0 {
			return errors.New("product count must be positive")
		}
		if p.Price <= 0 {
			return errors.New("product price must be positive")
		}
	}
	return nil
}
