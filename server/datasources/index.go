package datasources

import "context"

type Index interface {
	Search(context.Context, SearchOptions) SearchResponse
	Fields() ([]Field, error)
}
