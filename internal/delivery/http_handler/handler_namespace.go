package http_handler

import (
	"app/backendv1/internal/domain"
	"app/backendv1/internal/usecase"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

type handler struct {
	uc usecase.NamespaceUsecase
}

// NewHandler creates a new namespace handler
// @Summary Create namespace handler
func NewHandler(uc usecase.NamespaceUsecase) *handler {
	return &handler{uc: uc}
}

func (h *handler) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/namespaces", h.Create).Methods("POST")
	r.HandleFunc("/namespaces", h.GetAll).Methods("GET")
	r.HandleFunc("/namespaces/{code}", h.GetByCode).Methods("GET")
	r.HandleFunc("/namespaces/{code}", h.Update).Methods("PUT")
	r.HandleFunc("/namespaces/{code}", h.Delete).Methods("DELETE")
}

// Create creates a new namespace
// @Summary Create a namespace
// @Description Create a new namespace with the input payload
// @Tags namespaces
// @Accept  json
// @Produce  json
// @Param namespace body domain.Namespace true "Namespace data"
// @Success 201 {object} domain.Namespace
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /namespaces [post]
func (h *handler) Create(w http.ResponseWriter, r *http.Request) {
	var namespace domain.Namespace
	if err := json.NewDecoder(r.Body).Decode(&namespace); err != nil {
		http.Error(w, "bad request: "+err.Error(), http.StatusBadRequest)
		return
	}
	if err := h.uc.Create(&namespace); err != nil {
		http.Error(w, "error creating"+err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(namespace)
}

// GetAll gets all namespaces
// @Summary Get all namespaces
// @Description Get list of all namespaces
// @Tags namespaces
// @Produce  json
// @Success 200 {array} domain.Namespace
// @Failure 500 {string} string "Internal Server Error"
// @Router /namespaces [get]
func (h *handler) GetAll(w http.ResponseWriter, r *http.Request) {
	namespaces, err := h.uc.GetAll()
	if err != nil {
		http.Error(w, "error loading all", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(namespaces)
}

// GetByCode gets a namespace by code
// @Summary Get namespace by code
// @Description Get namespace details by its code
// @Tags namespaces
// @Produce  json
// @Param code path string true "Namespace code"
// @Success 200 {object} domain.Namespace
// @Failure 404 {string} string "Not Found"
// @Failure 500 {string} string "Internal Server Error"
// @Router /namespaces/{code} [get]
func (h *handler) GetByCode(w http.ResponseWriter, r *http.Request) {
	code := mux.Vars(r)["namespace"]
	namespace, err := h.uc.GetByCode(code)
	if err != nil {
		http.Error(w, "error loading by code", http.StatusInternalServerError)
		return
	}
	if namespace == nil {
		http.NotFound(w, r)
		return
	}
	json.NewEncoder(w).Encode(namespace)
}

// Update updates a namespace
// @Summary Update namespace
// @Description Update namespace details
// @Tags namespaces
// @Accept  json
// @Produce  json
// @Param code path string true "Namespace code"
// @Param namespace body domain.Namespace true "Namespace data"
// @Success 200 {object} domain.Namespace
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /namespaces/{code} [put]
func (h *handler) Update(w http.ResponseWriter, r *http.Request) {
	code := mux.Vars(r)["namespace"]
	var namespace domain.Namespace
	if err := json.NewDecoder(r.Body).Decode(&namespace); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	if err := h.uc.Update(code, &namespace); err != nil {
		http.Error(w, "error update", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(namespace)
}

// Delete deletes a namespace
// @Summary Delete namespace
// @Description Delete namespace by code
// @Tags namespaces
// @Param code path string true "Namespace code"
// @Success 204
// @Failure 500 {string} string "Internal Server Error"
// @Router /namespaces/{code} [delete]
func (h *handler) Delete(w http.ResponseWriter, r *http.Request) {
	code := mux.Vars(r)["namespace"]
	if err := h.uc.Delete(code); err != nil {
		http.Error(w, "delete failed", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
