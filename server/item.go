package server

type Item struct {
	ID        string                 `json:"id"`
	Fields    map[string]interface{} `json:"fields"`
	Highlight map[string][]string    `json:"highlight"`
}
