package datasources

import "context"

type Index interface {
	Search(context.Context, SearchOptions) SearchResponse
	// Items(context.Context, ItemsOptions) ItemsResponse
	GetFields(context.Context) ([]Field, error)

	Type() string
}

type SetNamerer interface {
	SetName(name string) string
}

type SetHuberer interface {
	SetHub(name string) string
}
