package order

import (
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
