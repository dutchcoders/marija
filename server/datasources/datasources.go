package datasources

import (
	"errors"
	"net/url"
)

var (
	datasources = map[string]DatasourceFn{}
	ErrNotFound = errors.New("Datasource not found.")
)

type DatasourceFn func(u *url.URL) (Index, error)

func Register(name string, fn DatasourceFn) DatasourceFn {
	datasources[name] = fn
	return fn
}

func Get(name string) (DatasourceFn, error) {
	if ds, ok := datasources[name]; ok {
		return ds, nil
	}

	return nil, ErrNotFound
}

type SearchResponse interface {
	Item() chan Item
	Error() chan error
}

func NewSearchResponse(itemChan chan Item, errorChan chan error) *searchResponse {
	return &searchResponse{
		itemChan,
		errorChan,
	}
}

type searchResponse struct {
	itemChan  chan Item
	errorChan chan error
}

func (sr *searchResponse) Item() chan Item {
	return sr.itemChan
}

func (sr *searchResponse) Error() chan error {
	return sr.errorChan
}
