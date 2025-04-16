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

func (h *handler) GetAll(w http.ResponseWriter, r *http.Request) {
	namespaces, err := h.uc.GetAll()
	if err != nil {
		http.Error(w, "error loading all", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(namespaces)
}

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

func (h *handler) Delete(w http.ResponseWriter, r *http.Request) {
	code := mux.Vars(r)["namespace"]
	if err := h.uc.Delete(code); err != nil {
		http.Error(w, "delete failed", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
