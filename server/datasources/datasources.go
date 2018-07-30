package datasources

import (
	"errors"

	"github.com/BurntSushi/toml"
)

var (
	datasources = map[string]DatasourceFn{}

	ErrNotFound = errors.New("Datasource not found.")
)

type DatasourceFn func(options ...func(Index) error) (Index, error)

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

func WithConfig(c toml.Primitive) func(Index) error {
	return func(d Index) error {
		err := toml.PrimitiveDecode(c, d)
		return err
	}
}

type AdvancedQuery struct {
	Field    string `json:"field"`
	Operator string `json:"operator"`
	Value    string `json:"value"`
}
