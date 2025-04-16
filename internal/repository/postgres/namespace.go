package postgres

import (
	"app/backendv1/internal/domain"
	"database/sql"
	"errors"
)

type namespaceRepo struct {
	db *sql.DB
}

func NewNamespaceRepo(db *sql.DB) *namespaceRepo {
	return &namespaceRepo{db: db}
}

func (r *namespaceRepo) Create(namespace *domain.Namespace) error {
	_, err := r.db.Exec("INSERT INTO namespaces (code, name) VALUES ($1, $2)", namespace.Code, namespace.Name)
	return err
}

func (r *namespaceRepo) GetAll() ([]domain.Namespace, error) {
	rows, err := r.db.Query("SELECT code, name FROM namespaces")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var namespaces []domain.Namespace
	for rows.Next() {
		var namespace domain.Namespace
		if err := rows.Scan(&namespace.Code, &namespace.Name); err != nil {
			return nil, err
		}
		namespaces = append(namespaces, namespace)
	}
	return namespaces, nil
}

func (r *namespaceRepo) GetByCode(code string) (*domain.Namespace, error) {
	var namespace domain.Namespace
	err := r.db.QueryRow("SELECT code, name FROM namespaces WHERE code = $1", code).Scan(&namespace.Code, &namespace.Name)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return &namespace, err
}

func (r *namespaceRepo) Update(code string, namespace *domain.Namespace) error {
	_, err := r.db.Exec("UPDATE namespaces SET name = $1 WHERE code = $2", namespace.Name, code)
	return err
}

func (r *namespaceRepo) Delete(code string) error {
	_, err := r.db.Exec("DELETE FROM namespaces WHERE code = $1", code)
	return err
}
