package coupon

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"suto-e-shop-api/pkg/pagination"
)

// Handler holds the coupon service.
type Handler struct {
	service Service
}

// NewHandler creates a new coupon handler.
func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

// RegisterRoutes registers the coupon routes to the router.
func (h *Handler) RegisterRoutes(router *mux.Router) {
	adminRouter := router.PathPrefix("/coupon").Subrouter()

	adminRouter.HandleFunc("", h.CreateCoupon).Methods("POST")
	adminRouter.HandleFunc("", h.GetCoupons).Methods("GET")
	adminRouter.HandleFunc("/{id}", h.GetCoupon).Methods("GET")
	adminRouter.HandleFunc("/{id}", h.UpdateCoupon).Methods("PUT")
	adminRouter.HandleFunc("/{id}", h.DeleteCoupon).Methods("DELETE")
}

func (h *Handler) CreateCoupon(w http.ResponseWriter, r *http.Request) {
	var coupon Coupon
	if err := json.NewDecoder(r.Body).Decode(&coupon); err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	createdCoupon, err := h.service.CreateCoupon(r.Context(), coupon)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	RespondWithJSON(w, http.StatusCreated, Response{Data: createdCoupon, Message: "success", Code: 0})
}

func (h *Handler) GetCoupons(w http.ResponseWriter, r *http.Request) {
	page, pageSize := pagination.GetPaginationParams(r)
	search := r.URL.Query().Get("search")

	coupons, totalCount, err := h.service.GetCoupons(r.Context(), page, pageSize, search)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	paginator := pagination.New(page, pageSize, totalCount)

	RespondWithJSON(w, http.StatusOK, PaginatedResponse{
		Data:       coupons,
		Pagination: paginator,
		Message:    "success",
		Code:       0,
	})
}

func (h *Handler) GetCoupon(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	coupon, err := h.service.GetCoupon(r.Context(), id)
	if err != nil {
		RespondWithError(w, http.StatusNotFound, "Coupon not found")
		return
	}

	RespondWithJSON(w, http.StatusOK, Response{Data: coupon, Message: "success", Code: 0})
}

func (h *Handler) UpdateCoupon(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var coupon Coupon
	if err := json.NewDecoder(r.Body).Decode(&coupon); err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	updatedCoupon, err := h.service.UpdateCoupon(r.Context(), id, coupon)
	if err != nil {
		RespondWithError(w, http.StatusNotFound, "Coupon not found")
		return
	}

	RespondWithJSON(w, http.StatusOK, Response{Data: updatedCoupon, Message: "success", Code: 0})
}

func (h *Handler) DeleteCoupon(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if err := h.service.DeleteCoupon(r.Context(), id); err != nil {
		RespondWithError(w, http.StatusNotFound, "Coupon not found")
		return
	}

	RespondWithJSON(w, http.StatusOK, Response{Message: "success", Code: 0})
}
