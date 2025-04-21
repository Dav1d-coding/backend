package domain

type AppData struct {
	UID  string                 `json:"uid"`  // Уникальный идентификатор
	Data map[string]interface{} `json:"data"` // Произвольные JSON данные
}
