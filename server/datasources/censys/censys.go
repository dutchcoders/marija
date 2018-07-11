package censys

import (
	"context"
	"net"
	"strings"

	"github.com/op/go-logging"

	"github.com/dutchcoders/marija/server/datasources"
	"github.com/dutchcoders/marija/server/datasources/censys/api"
)

var (
	_ = datasources.Register("censys", New)
)

func New(options ...func(datasources.Index) error) (datasources.Index, error) {
	s := Censys{}

	for _, optionFn := range options {
		optionFn(&s)
	}

	return &s, nil
}

type Config struct {
	ApiID     string
	ApiSecret string
}

var log = logging.MustGetLogger("marija/datasources/censys")

type Censys struct {
	Config
}

func (m *Config) UnmarshalTOML(p interface{}) error {
	data, _ := p.(map[string]interface{})

	if v, ok := data["api-id"]; !ok {
	} else if v, ok := v.(string); !ok {
	} else {
		m.ApiID = v
	}

	if v, ok := data["api-secret"]; !ok {
	} else if v, ok := v.(string); !ok {
	} else {
		m.ApiSecret = v
	}

	return nil
}

func (m *Censys) Type() string {
	return "censys"
}

func (b *Censys) Search(ctx context.Context, so datasources.SearchOptions) datasources.SearchResponse {
	itemCh := make(chan datasources.Item)
	errorCh := make(chan error)

	cs := api.New(b.ApiID, b.ApiSecret)

	go func() {
		defer close(itemCh)
		defer close(errorCh)

		qry := strings.Replace(so.Query, "\"", "", -1)

		if net.ParseIP(qry) == nil {
			return
		}

		vo, err := cs.IPv4.View(qry)
		if err != nil {
			errorCh <- err
			return
		}

		item := datasources.Item{
			ID:     so.Query,
			Fields: flattenFields("", *vo),
		}

		select {
		case itemCh <- item:
		case <-ctx.Done():
			return
		}
	}()

	return datasources.NewSearchResponse(
		itemCh,
		errorCh,
	)
}

func flattenFields(root string, m map[string]interface{}) map[string]interface{} {
	fields := map[string]interface{}{}

	for k, v := range m {
		key := k
		if root != "" {
			key = root + "." + key
		}

		switch s2 := v.(type) {
		case map[string]interface{}:
			for k2, v2 := range flattenFields(key, s2) {
				fields[k2] = v2
			}
		default:
			fields[key] = v
		}
	}

	return fields
}

func (i *Censys) GetFields(context.Context) (fields []datasources.Field, err error) {
	fields = append(fields, datasources.Field{
		Path: "80.http.get.title",
		Type: "string",
	})

	fields = append(fields, datasources.Field{
		Path: "80.http.get.metadata.description",
		Type: "string",
	})

	fields = append(fields, datasources.Field{
		Path: "metadata.os",
		Type: "string",
	})

	fields = append(fields, datasources.Field{
		Path: "ip",
		Type: "string",
	})

	return
}
