package live

import (
	"context"

	"github.com/op/go-logging"

	"github.com/dutchcoders/marija/server/datasources"
)

var log = logging.MustGetLogger("marija/datasources/live")

var (
	_ = datasources.Register("live", New)
)

func New(options ...func(datasources.Index) error) (datasources.Index, error) {
	s := Live{}

	return &s, nil
}

type Live struct {
}

func (m *Live) Type() string {
	return "live"
}

func (m *Live) UnmarshalTOML(p interface{}) error {
	return nil
}

func (i *Live) Search(ctx context.Context, so datasources.SearchOptions) datasources.SearchResponse {
	itemCh := make(chan datasources.Item)
	errorCh := make(chan error)

	return datasources.NewSearchResponse(
		itemCh,
		errorCh,
	)
}

func (i *Live) GetFields(context.Context) (fields []datasources.Field, err error) {
	return
}
