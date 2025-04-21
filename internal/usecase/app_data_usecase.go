package usecase

import "app/backendv1/internal/domain"

type AppDataUsecase interface {
	Create(namespace, appName string, data *domain.AppData) error
	GetDataByUID(namespace, appName, uid string) (*domain.AppData, error)
	GetAll(namespace, appName string) ([]*domain.AppData, error)
	Update(namespace, appName string, data *domain.AppData) error
	UpdateDataPartial(namespace, appName, uid string, partialData map[string]interface{}) error
	Delete(namespace, appName, uid string) error
}

type appDataUsecase struct {
	repo AppDataUsecase
}

func NewAppDataUsecase(repo AppDataUsecase) AppDataUsecase {
	return &appDataUsecase{repo: repo}
}

func (u *appDataUsecase) Create(namespace, appName string, data *domain.AppData) error {
	return u.repo.Create(namespace, appName, data)
}

func (u *appDataUsecase) GetDataByUID(namespace, appName, uid string) (*domain.AppData, error) {
	return u.repo.GetDataByUID(namespace, appName, uid)
}

func (u *appDataUsecase) GetAll(namespace, appName string) ([]*domain.AppData, error) {
	return u.repo.GetAll(namespace, appName)
}

func (u *appDataUsecase) Update(namespace, appName string, data *domain.AppData) error {
	return u.repo.Update(namespace, appName, data)
}

func (u *appDataUsecase) UpdateDataPartial(namespace, appName, uid string, partialData map[string]interface{}) error {
	return u.repo.UpdateDataPartial(namespace, appName, uid, partialData)
}

func (u *appDataUsecase) Delete(namespace, appName, uid string) error {
	return u.repo.Delete(namespace, appName, uid)
}
