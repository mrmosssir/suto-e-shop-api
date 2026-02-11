package advertise

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"suto-e-shop-api/pkg/pagination"
)

// Handler holds the advertise service.
type Handler struct {
	service Service
}

// NewHandler creates a new advertise handler.
func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

// RegisterAdminRoutes registers the admin advertise routes to the router.
func (h *Handler) RegisterAdminRoutes(router *mux.Router) {
	adminRouter := router.PathPrefix("/advertise").Subrouter()

	adminRouter.HandleFunc("", h.AdminCreateAdvertise).Methods("POST")
	adminRouter.HandleFunc("", h.AdminGetAdvertises).Methods("GET")
	adminRouter.HandleFunc("/{id}", h.AdminGetAdvertise).Methods("GET")
	adminRouter.HandleFunc("/{id}", h.AdminUpdateAdvertise).Methods("PUT")
	adminRouter.HandleFunc("/{id}", h.AdminDeleteAdvertise).Methods("DELETE")
}

// RegisterClientRoutes registers the client advertise routes to the router.
func (h *Handler) RegisterClientRoutes(router *mux.Router) {
	router.HandleFunc("/advertises", h.GetAdvertises).Methods("GET")
}

func (h *Handler) AdminCreateAdvertise(w http.ResponseWriter, r *http.Request) {
	var advertise Advertise
	if err := json.NewDecoder(r.Body).Decode(&advertise); err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	createdAdvertise, err := h.service.AdminCreateAdvertise(r.Context(), advertise)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	RespondWithJSON(w, http.StatusCreated, Response{Data: createdAdvertise, Message: "success", Code: 0})
}

func (h *Handler) AdminGetAdvertises(w http.ResponseWriter, r *http.Request) {
	page, pageSize := pagination.GetPaginationParams(r)
	search := r.URL.Query().Get("search")

	advertises, totalCount, err := h.service.AdminGetAdvertises(r.Context(), page, pageSize, search)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	paginator := pagination.New(page, pageSize, totalCount)

	RespondWithJSON(w, http.StatusOK, PaginatedResponse{
		Data:       advertises,
		Pagination: paginator,
		Message:    "success",
		Code:       0,
	})
}

func (h *Handler) AdminGetAdvertise(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	advertise, err := h.service.AdminGetAdvertise(r.Context(), id)
	if err != nil {
		RespondWithError(w, http.StatusNotFound, "Advertise not found")
		return
	}

	RespondWithJSON(w, http.StatusOK, Response{Data: advertise, Message: "success", Code: 0})
}

func (h *Handler) AdminUpdateAdvertise(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var advertise Advertise
	if err := json.NewDecoder(r.Body).Decode(&advertise); err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	updatedAdvertise, err := h.service.AdminUpdateAdvertise(r.Context(), id, advertise)
	if err != nil {
		RespondWithError(w, http.StatusNotFound, "Advertise not found")
		return
	}

	RespondWithJSON(w, http.StatusOK, Response{Data: updatedAdvertise, Message: "success", Code: 0})
}

func (h *Handler) AdminDeleteAdvertise(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if err := h.service.AdminDeleteAdvertise(r.Context(), id); err != nil {
		RespondWithError(w, http.StatusNotFound, "Advertise not found")
		return
	}

	RespondWithJSON(w, http.StatusOK, Response{Message: "success", Code: 0})
}

func (h *Handler) GetAdvertises(w http.ResponseWriter, r *http.Request) {
	advertises, err := h.service.GetAdvertises(r.Context())
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	RespondWithJSON(w, http.StatusOK, Response{Data: advertises, Message: "success", Code: 0})
}
