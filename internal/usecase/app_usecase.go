package usecase

import "app/backendv1/internal/domain"

type AppUsecase interface {
	Create(app *domain.App) error
	GetAll() ([]*domain.App, error)
	GetAllByCodeNamespace(code string) ([]*domain.App, error)
	Update(app *domain.App) error
	Delete(code, namespaceCode string) error
}

type appUsecase struct {
	repo AppUsecase
}

func NewAppUsecase(repo AppUsecase) AppUsecase {
	return &appUsecase{repo: repo}
}

func (u *appUsecase) Create(app *domain.App) error {
	return u.repo.Create(app)
}

func (u *appUsecase) GetAll() ([]*domain.App, error) {
	return u.repo.GetAll()
}

func (u *appUsecase) GetAllByCodeNamespace(code string) ([]*domain.App, error) {
	return u.repo.GetAllByCodeNamespace(code)
}

func (u *appUsecase) Update(app *domain.App) error {
	return u.repo.Update(app)
}

func (u *appUsecase) Delete(code, namespaceCode string) error {
	return u.repo.Delete(code, namespaceCode)
}
