package domain

type App struct {
	Code          string `json:"code"`
	Name          string `json:"name"`
	NamespaceCode string `json:"namespaceCode"`
	Icon          string `json:"icon"`
}
