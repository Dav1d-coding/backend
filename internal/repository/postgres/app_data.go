package postgres

import (
	"app/backendv1/internal/domain"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
)

type appDataRepo struct {
	db *sql.DB
}

func NewAppDataRepo(db *sql.DB) *appDataRepo {
	return &appDataRepo{db: db}
}

// Create создает новую запись
func (r *appDataRepo) Create(namespace, table string, data *domain.AppData) error {
	jsonData, err := json.Marshal(data.Data)
	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}

	query := fmt.Sprintf(`
		INSERT INTO %s.%s (data) 
		VALUES ($1)
	`, namespace, table)

	_, err = r.db.Exec(query, jsonData)
	if err != nil {
		return fmt.Errorf("failed to insert data: %w", err)
	}

	return nil
}

// GetByUID возвращает запись по UID
func (r *appDataRepo) GetDataByUID(namespace, table, uid string) (*domain.AppData, error) {
	query := fmt.Sprintf(`
		SELECT uid, data 
		FROM %s.%s 
		WHERE uid = $1
	`, namespace, table)

	var (
		dbUID    string
		jsonData []byte
	)

	err := r.db.QueryRow(query, uid).Scan(&dbUID, &jsonData)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("record with uid %s not found", uid)
		}
		return nil, fmt.Errorf("failed to query data: %w", err)
	}

	var data map[string]interface{}
	if err := json.Unmarshal(jsonData, &data); err != nil {
		return nil, fmt.Errorf("failed to unmarshal data: %w", err)
	}

	return &domain.AppData{
		UID:  dbUID,
		Data: data,
	}, nil
}

// GetAll возвращает все записи
func (r *appDataRepo) GetAll(namespace, table string) ([]*domain.AppData, error) {
	query := fmt.Sprintf(`
		SELECT uid, data 
		FROM %s.%s
	`, namespace, table)

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query data: %w", err)
	}
	defer rows.Close()

	var results []*domain.AppData
	for rows.Next() {
		var (
			uid      string
			jsonData []byte
		)

		if err := rows.Scan(&uid, &jsonData); err != nil {
			return nil, fmt.Errorf("failed to scan data: %w", err)
		}

		var data map[string]interface{}
		if err := json.Unmarshal(jsonData, &data); err != nil {
			return nil, fmt.Errorf("failed to unmarshal data: %w", err)
		}

		results = append(results, &domain.AppData{
			UID:  uid,
			Data: data,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return results, nil
}

// Update полностью обновляет запись
func (r *appDataRepo) Update(namespace, table string, data *domain.AppData) error {
	jsonData, err := json.Marshal(data.Data)
	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}

	query := fmt.Sprintf(`
		UPDATE %s.%s 
		SET data = $1
		WHERE uid = $2
	`, namespace, table)

	result, err := r.db.Exec(query, jsonData, data.UID)
	if err != nil {
		return fmt.Errorf("failed to update data: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("record with uid %s not found", data.UID)
	}

	return nil
}

// UpdatePartial частично обновляет JSON данные
func (r *appDataRepo) UpdateDataPartial(namespace, table, uid string, partialData map[string]interface{}) error {
	setParts := make([]string, 0, len(partialData))
	args := make([]interface{}, 0, len(partialData)+1)
	argPos := 1

	for field, value := range partialData {
		setParts = append(setParts, fmt.Sprintf("data = jsonb_set(data, '{%s}', $%d)", field, argPos))
		jsonValue, err := json.Marshal(value)
		if err != nil {
			return fmt.Errorf("failed to marshal field %s: %w", field, err)
		}
		args = append(args, jsonValue)
		argPos++
	}

	args = append(args, uid)

	query := fmt.Sprintf(`
		UPDATE %s.%s 
		SET %s
		WHERE uid = $%d
	`, namespace, table, strings.Join(setParts, ", "), argPos)

	result, err := r.db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("failed to update data: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("record with uid %s not found", uid)
	}

	return nil
}

// Delete удаляет запись
func (r *appDataRepo) Delete(namespace, table, uid string) error {
	query := fmt.Sprintf(`
		DELETE FROM %s.%s 
		WHERE uid = $1
	`, namespace, table)

	result, err := r.db.Exec(query, uid)
	if err != nil {
		return fmt.Errorf("failed to delete data: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("record with uid %s not found", uid)
	}

	return nil
}
