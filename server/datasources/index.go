package datasources

import "context"

type Index interface {
	Search(context.Context, SearchOptions) SearchResponse
	// Items(context.Context, ItemsOptions) ItemsResponse
	GetFields(context.Context) ([]Field, error)
}
