package datasources

type Index interface {
	Search(SearchOptions) ([]Item, error)
	Indices() ([]string, error)
	Fields(string) ([]Field, error)
}
