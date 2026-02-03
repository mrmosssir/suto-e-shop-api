package banner

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"suto-e-shop-api/pkg/pagination"
)

// Handler holds the banner service.
type Handler struct {
	service Service
}

// NewHandler creates a new banner handler.
func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

// RegisterAdminRoutes registers the admin banner routes to the router.
func (h *Handler) RegisterAdminRoutes(router *mux.Router) {
	adminRouter := router.PathPrefix("/banner").Subrouter()

	adminRouter.HandleFunc("", h.AdminCreateBanner).Methods("POST")
	adminRouter.HandleFunc("", h.AdminGetBanners).Methods("GET")
	adminRouter.HandleFunc("/{id}", h.AdminGetBanner).Methods("GET")
	adminRouter.HandleFunc("/{id}", h.AdminUpdateBanner).Methods("PUT")
	adminRouter.HandleFunc("/{id}", h.AdminDeleteBanner).Methods("DELETE")
}

// RegisterClientRoutes registers the client banner routes to the router.
func (h *Handler) RegisterClientRoutes(router *mux.Router) {
	router.HandleFunc("/banners", h.GetBanners).Methods("GET")
}

func (h *Handler) AdminCreateBanner(w http.ResponseWriter, r *http.Request) {
	var banner Banner
	if err := json.NewDecoder(r.Body).Decode(&banner); err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	createdBanner, err := h.service.AdminCreateBanner(r.Context(), banner)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	RespondWithJSON(w, http.StatusCreated, Response{Data: createdBanner, Message: "success", Code: 0})
}

func (h *Handler) AdminGetBanners(w http.ResponseWriter, r *http.Request) {
	page, pageSize := pagination.GetPaginationParams(r)
	search := r.URL.Query().Get("search")

	banners, totalCount, err := h.service.AdminGetBanners(r.Context(), page, pageSize, search)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	paginator := pagination.New(page, pageSize, totalCount)

	RespondWithJSON(w, http.StatusOK, PaginatedResponse{
		Data:       banners,
		Pagination: paginator,
		Message:    "success",
		Code:       0,
	})
}

func (h *Handler) AdminGetBanner(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	banner, err := h.service.AdminGetBanner(r.Context(), id)
	if err != nil {
		RespondWithError(w, http.StatusNotFound, "Banner not found")
		return
	}

	RespondWithJSON(w, http.StatusOK, Response{Data: banner, Message: "success", Code: 0})
}

func (h *Handler) AdminUpdateBanner(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var banner Banner
	if err := json.NewDecoder(r.Body).Decode(&banner); err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	updatedBanner, err := h.service.AdminUpdateBanner(r.Context(), id, banner)
	if err != nil {
		RespondWithError(w, http.StatusNotFound, "Banner not found")
		return
	}

	RespondWithJSON(w, http.StatusOK, Response{Data: updatedBanner, Message: "success", Code: 0})
}

func (h *Handler) AdminDeleteBanner(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if err := h.service.AdminDeleteBanner(r.Context(), id); err != nil {
		RespondWithError(w, http.StatusNotFound, "Banner not found")
		return
	}

	RespondWithJSON(w, http.StatusOK, Response{Message: "success", Code: 0})
}

func (h *Handler) GetBanners(w http.ResponseWriter, r *http.Request) {
	banners, err := h.service.GetBanners(r.Context())
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	RespondWithJSON(w, http.StatusOK, Response{Data: banners, Message: "success", Code: 0})
}
