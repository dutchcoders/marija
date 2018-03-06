package datasources

type Node struct {
	ID     string                 `json:"id"`
	Fields map[string]interface{} `json:"fields"`
	Count  int                    `json:"count"`
}
