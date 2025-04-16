package usecase

import "app/backendv1/internal/domain"

type NamespaceUsecase interface {
	Create(*domain.Namespace) error
	GetAll() ([]domain.Namespace, error)
	GetByCode(code string) (*domain.Namespace, error)
	Update(code string, record *domain.Namespace) error
	Delete(code string) error
}

// RecordService — конкретная реализация бизнес-логики
type namespaceService struct {
	repo NamespaceUsecase
}

func NewNamespaceService(repo NamespaceUsecase) *namespaceService {
	return &namespaceService{repo: repo}
}

// Реализация интерфейса RecordUsecase

func (s *namespaceService) Create(record *domain.Namespace) error {
	// можно добавить валидацию или другую бизнес-логику
	return s.repo.Create(record)
}

func (s *namespaceService) GetAll() ([]domain.Namespace, error) {
	return s.repo.GetAll()
}

func (s *namespaceService) GetByCode(code string) (*domain.Namespace, error) {
	return s.repo.GetByCode(code)
}

func (s *namespaceService) Update(code string, record *domain.Namespace) error {
	return s.repo.Update(code, record)
}

func (s *namespaceService) Delete(code string) error {
	return s.repo.Delete(code)
}
