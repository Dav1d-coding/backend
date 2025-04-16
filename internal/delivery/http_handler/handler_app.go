package http_handler

import (
	"app/backendv1/internal/domain"
	"app/backendv1/internal/usecase"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

type appHandler struct {
	uc usecase.AppUsecase
}

func NewAppHandler(uc usecase.AppUsecase) *appHandler {
	return &appHandler{uc: uc}
}

func (h *appHandler) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/namespace/{namespace}/app", h.Create).Methods("POST")
	r.HandleFunc("/apps", h.GetAll).Methods("GET")
	r.HandleFunc("/namespace/{namespace}/apps", h.GetAllByCodeNamespace).Methods("GET")
	r.HandleFunc("/namespace/{namespace}/app/{app}", h.Update).Methods("PUT")
	r.HandleFunc("/namespace/{namespace}/app/{app}", h.Delete).Methods("DELETE")
}

func (h *appHandler) Create(w http.ResponseWriter, r *http.Request) {
	namespaceCode := mux.Vars(r)["namespace"]
	fmt.Print(namespaceCode)
	var app domain.App
	if err := json.NewDecoder(r.Body).Decode(&app); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	// Важный момент: привязка к namespace из пути
	app.NamespaceCode = namespaceCode

	if err := h.uc.Create(&app); err != nil {
		http.Error(w, "create failed", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(app)
}

func (h *appHandler) GetAllByCodeNamespace(w http.ResponseWriter, r *http.Request) {
	code := mux.Vars(r)["namespace"]
	apps, err := h.uc.GetAllByCodeNamespace(code)
	if err != nil {
		http.Error(w, "error", http.StatusInternalServerError)
	}
	json.NewEncoder(w).Encode(apps)
}

func (h *appHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	apps, err := h.uc.GetAll()
	if err != nil {
		http.Error(w, "get all failed", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(apps)
}

func (h *appHandler) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	code := vars["app"]
	namespaceCode := vars["namespace"]

	var app domain.App
	if err := json.NewDecoder(r.Body).Decode(&app); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	// Обязательно подставляем ключи из URL
	app.Code = code
	app.NamespaceCode = namespaceCode

	if err := h.uc.Update(&app); err != nil {
		http.Error(w, "update failed", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(app)
}

func (h *appHandler) Delete(w http.ResponseWriter, r *http.Request) {
	appCode := mux.Vars(r)["app"]
	namespaceCode := mux.Vars(r)["namespace"]

	if err := h.uc.Delete(appCode, namespaceCode); err != nil {
		http.Error(w, "delete failed", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
