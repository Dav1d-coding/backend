package http_handler

import (
	"app/backendv1/internal/domain"
	"app/backendv1/internal/usecase"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

type appDataHandler struct {
	uc usecase.AppDataUsecase
}

func NewAppDataHandler(uc usecase.AppDataUsecase) *appDataHandler {
	return &appDataHandler{uc: uc}
}

func (h *appDataHandler) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/namespace/{namespace}/app/{app}/data", h.Create).Methods("POST")
	r.HandleFunc("/namespace/{namespace}/app/{app}/data/{uid}", h.GetDataByUID).Methods("GET")
	r.HandleFunc("/namespace/{namespace}/app/{app}/data", h.GetAll).Methods("GET")
	r.HandleFunc("/namespace/{namespace}/app/{app}/data/{uid}", h.Update).Methods("PUT")
	r.HandleFunc("/namespace/{namespace}/app/{app}/data/{uid}", h.UpdateDataPartial).Methods("PATCH")
	r.HandleFunc("/namespace/{namespace}/app/{app}/data/{uid}", h.Delete).Methods("DELETE")
}

// CreateDataHandler godoc
// @Summary Создать новые данные приложения
// @Description Создаёт новые данные для указанного приложения в namespace
// @Tags app-data
// @Accept json
// @Produce json
// @Param namespace path string true "Namespace Code"
// @Param app path string true "App Code"
// @Param data body domain.AppData true "Данные приложения"
// @Success 201 {object} domain.AppData
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /namespace/{namespace}/app/{app}/data [post]
func (h *appDataHandler) Create(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	namespace := vars["namespace"]
	appName := vars["app"]

	var data domain.AppData
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.uc.Create(namespace, appName, &data); err != nil {
		http.Error(w, "failed to create data", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(data)
}

// GetDataByUIDHandler godoc
// @Summary Получить данные по UID
// @Description Возвращает данные приложения по уникальному идентификатору
// @Tags app-data
// @Produce json
// @Param namespace path string true "Namespace Code"
// @Param app path string true "App Code"
// @Param uid path string true "Data UID"
// @Success 200 {object} domain.AppData
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /namespace/{namespace}/app/{app}/data/{uid} [get]
func (h *appDataHandler) GetDataByUID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	namespace := vars["namespace"]
	appName := vars["app"]
	uid := vars["uid"]

	data, err := h.uc.GetDataByUID(namespace, appName, uid)
	if err != nil {
		http.Error(w, "data not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(data)
}

// GetAllDataHandler godoc
// @Summary Получить все данные приложения
// @Description Возвращает все данные для указанного приложения в namespace
// @Tags app-data
// @Produce json
// @Param namespace path string true "Namespace Code"
// @Param app path string true "App Code"
// @Success 200 {array} domain.AppData
// @Failure 500 {object} map[string]string
// @Router /namespace/{namespace}/app/{app}/data [get]
func (h *appDataHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	namespace := vars["namespace"]
	appName := vars["app"]

	data, err := h.uc.GetAll(namespace, appName)
	if err != nil {
		http.Error(w, "failed to get data", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(data)
}

// UpdateDataHandler godoc
// @Summary Полностью обновить данные
// @Description Заменяет все данные для указанного UID
// @Tags app-data
// @Accept json
// @Produce json
// @Param namespace path string true "Namespace Code"
// @Param app path string true "App Code"
// @Param uid path string true "Data UID"
// @Param data body domain.AppData true "Новые данные"
// @Success 200 {object} domain.AppData
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /namespace/{namespace}/app/{app}/data/{uid} [put]
func (h *appDataHandler) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	namespace := vars["namespace"]
	appName := vars["app"]
	uid := vars["uid"]

	var data domain.AppData
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	// Устанавливаем UID из пути
	data.UID = uid

	if err := h.uc.Update(namespace, appName, &data); err != nil {
		http.Error(w, "failed to update data", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(data)
}

// UpdateDataPartialHandler godoc
// @Summary Частично обновить данные
// @Description Обновляет только указанные поля данных
// @Tags app-data
// @Accept json
// @Produce json
// @Param namespace path string true "Namespace Code"
// @Param app path string true "App Code"
// @Param uid path string true "Data UID"
// @Param data body map[string]interface{} true "Поля для обновления"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /namespace/{namespace}/app/{app}/data/{uid} [patch]
func (h *appDataHandler) UpdateDataPartial(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	namespace := vars["namespace"]
	appName := vars["app"]
	uid := vars["uid"]

	var partialData map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&partialData); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.uc.UpdateDataPartial(namespace, appName, uid, partialData); err != nil {
		http.Error(w, "failed to update data", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "updated"})
}

// DeleteDataHandler godoc
// @Summary Удалить данные
// @Description Удаляет данные по UID
// @Tags app-data
// @Param namespace path string true "Namespace Code"
// @Param app path string true "App Code"
// @Param uid path string true "Data UID"
// @Success 204
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /namespace/{namespace}/app/{app}/data/{uid} [delete]
func (h *appDataHandler) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	namespace := vars["namespace"]
	appName := vars["app"]
	uid := vars["uid"]

	if err := h.uc.Delete(namespace, appName, uid); err != nil {
		http.Error(w, "failed to delete data", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
