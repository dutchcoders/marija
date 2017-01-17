package datasources

type Index interface {
	Search(SearchOptions) ([]Item, int, error)
	Fields() ([]Field, error)
}
