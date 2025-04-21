package postgres

import (
	"app/backendv1/internal/domain"
	"database/sql"
	"encoding/json"
	"fmt"
)

type appRepo struct {
	db *sql.DB
}

func NewAppRepo(db *sql.DB) *appRepo {
	return &appRepo{db: db}
}

func (r *appRepo) Create(app *domain.App) error {
	// Создаем базовую структуру для fields
	// Создаем структуру fields согласно требованиям
	fields := map[string]interface{}{
		"item": []interface{}{
			map[string]interface{}{
				"name": map[string]interface{}{
					"type": "string",
				},
			},
		},
	}

	// Преобразуем в JSON
	fieldsJSON, err := json.Marshal(fields)
	if err != nil {
		return fmt.Errorf("failed to marshal fields: %w", err)
	}
	// fmt.Print(app.Code, app.Name, app.NamespaceCode, app.Icon)
	_, err = r.db.Exec("INSERT INTO apps (code, name, namespace_code, icon, fields) VALUES ($1, $2, $3, $4, $5)", app.Code, app.Name, app.NamespaceCode, app.Icon, fieldsJSON)
	if err != nil {
		return fmt.Errorf("failed to marshal fields: %w", err)
	}
	query := "CREATE TABLE IF NOT EXISTS " + app.NamespaceCode + "." + app.Code + " (uid uuid PRIMARY KEY DEFAULT gen_random_uuid(), data jsonb not null default '{}'::jsonb)"
	_, err = r.db.Exec(query)
	// fmt.Printf(query)
	// fmt.Print(err.Error())
	return err
}

func (r *appRepo) GetAllByCodeNamespace(code string) ([]*domain.App, error) {
	rows, err := r.db.Query("SELECT code, name, namespace_code, icon FROM apps WHERE namespace_code = $1", code)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var apps []*domain.App
	for rows.Next() {
		var app domain.App
		if err := rows.Scan(&app.Code, &app.Name, &app.NamespaceCode, &app.Icon); err != nil {
			return nil, err
		}
		apps = append(apps, &app)
	}
	return apps, nil
}

func (r *appRepo) GetAll() ([]*domain.App, error) {
	query := `SELECT code, name, namespace_code, icon FROM apps`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var apps []*domain.App
	for rows.Next() {
		var app domain.App
		if err := rows.Scan(&app.Code, &app.Name, &app.NamespaceCode, &app.Icon); err != nil {
			return nil, err
		}
		apps = append(apps, &app)
	}
	return apps, nil
}

func (r *appRepo) Update(app *domain.App) error {
	result, err := r.db.Exec("UPDATE apps SET name = $1, icon = $2  WHERE code = $3 AND namespace_code = $4", app.Name, app.Icon, app.Code, app.NamespaceCode)
	if err != nil {
		return err
	}

	count, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf("no rows updated")
	}
	return nil
}

func (r *appRepo) Delete(code, namespace_code string) error {
	_, err := r.db.Exec("DELETE FROM apps WHERE code = $1 AND namespace_code = $2", code, namespace_code)
	return err
}
