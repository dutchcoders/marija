package datasources

type Graph struct {
	ID         string                 `json:"id"`
	Fields     map[string]interface{} `json:"fields"`
	Count      int                    `json:"count"`
	Datasource string                 `json:"datasource"`
}
