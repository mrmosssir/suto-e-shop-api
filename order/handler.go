package order

import (
	"encoding/json"
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

// RegisterRoutes registers the order routes to the router.
func (h *Handler) RegisterRoutes(router *mux.Router) {
	adminRouter := router.PathPrefix("/order").Subrouter()
	adminRouter.HandleFunc("", h.GetOrders).Methods("GET")
	adminRouter.HandleFunc("/{id}", h.UpdateOrder).Methods("PUT")
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
