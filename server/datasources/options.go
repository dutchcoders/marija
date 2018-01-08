package datasources

type SearchOptions struct {
	Size  int
	From  int
	Query string

	Fields []string
}
