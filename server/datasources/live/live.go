package live

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"hash/fnv"

	"github.com/op/go-logging"

	"github.com/dutchcoders/marija/server/datasources"
	"github.com/dutchcoders/marija/server/messages"
	"github.com/dutchcoders/marija/server/unique"
)

var log = logging.MustGetLogger("marija/datasources/live")

var (
	_ = datasources.Register("live", New)
)

func New(options ...func(datasources.Index) error) (datasources.Index, error) {
	s := &Live{
		liveCh: make(chan map[string]interface{}, 100),
	}

	for _, optionFn := range options {
		optionFn(s)
	}

	return s, nil
}

func (l *Live) Receive(m map[string]interface{}) {
	l.liveCh <- m
}

func (l *Live) Broadcast(ctx context.Context, datasource string) chan json.Marshaler {
	broadcastCh := make(chan json.Marshaler)

	go func() {
		defer close(broadcastCh)

		unique := unique.New()

		for {

			select {
			case fields := <-l.liveCh:
				// calculate hash of fields
				hash := fnv.New128()
				for _, field := range fields {
					switch s := field.(type) {
					case []string:
						for _, v := range s {
							hash.Write([]byte(v))
						}
					case string:
						hash.Write([]byte(s))
					default:
					}
				}

				hashHex := hex.EncodeToString(hash.Sum(nil))

				i := &datasources.Graph{
					ID:     hashHex,
					Fields: fields,
					Count:  1,
				}

				if v, ok := unique.Get(hash.Sum(nil)); ok {
					i = v

					i.Count++
				}

				unique.Add(hash.Sum(nil), i)

				select {
				case broadcastCh <- &messages.LiveResponse{
					Datasource: datasource,
					Graphs: []datasources.Graph{
						*i,
					},
				}:
				default:
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	return broadcastCh
}

type Live struct {
	liveCh chan map[string]interface{}
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
