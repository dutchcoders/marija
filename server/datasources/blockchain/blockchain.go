package btc

import (
	"context"
	"net/http"
	"net/url"

	"fmt"
	"time"

	"github.com/dutchcoders/marija/server/datasources"
	logging "github.com/op/go-logging"
	"github.com/qedus/blockchain"
)

var (
	_ = datasources.Register("blockchain", New)
)

var log = logging.MustGetLogger("marija/datasources/blockchain")

func New(options ...func(datasources.Index) error) (datasources.Index, error) {
	s := BTC{}

	for _, optionFn := range options {
		optionFn(&s)
	}

	s.client = blockchain.New(http.DefaultClient)

	return &s, nil
}

func (m *BTC) Type() string {
	return "blockchain"
}

type Config struct {
}

func (m *BTC) UnmarshalTOML(p interface{}) error {
	data, _ := p.(map[string]interface{})

	_ = data
	return nil
}

type BTC struct {
	Config

	client *blockchain.BlockChain

	u *url.URL
}

func (b *BTC) Search(ctx context.Context, so datasources.SearchOptions) datasources.SearchResponse {
	itemCh := make(chan datasources.Item)
	errorCh := make(chan error)

	go func() {
		address := &blockchain.Address{Address: so.Query}
		if err := b.client.Request(address); err != nil {
			errorCh <- err
			return
		}

		for {

			tx, err := address.NextTransaction()
			if err == blockchain.IterDone {
				break
			} else if err != nil {
				errorCh <- err
				return
			}

			for _, input := range tx.Inputs {
				for _, output := range tx.Outputs {
					fields := map[string]interface{}{
						"relayed_by":       tx.RelayedBy,
						"date":             time.Unix(tx.Time, 0),
						"input":            input.PrevOut.Address,
						"output":           output.Address,
						"value":            input.PrevOut.Value,
						"address_tag":      output.AddressTag,
						"address_tag_link": output.AddressTagLink,
					}

					item := datasources.Item{
						ID:     fmt.Sprintf("%s%s%s", tx.Hash, input.PrevOut.Address, output.Address),
						Fields: fields,
					}

					select {
					case itemCh <- item:
					case <-ctx.Done():
						return
					}
				}
			}
		}
	}()

	return datasources.NewSearchResponse(
		itemCh,
		errorCh,
	)
}

func (i *BTC) GetFields(context.Context) (fields []datasources.Field, err error) {
	fields = append(fields, datasources.Field{
		Path: "id",
		Type: "string",
	})

	fields = append(fields, datasources.Field{
		Path: "relayed_by",
		Type: "string",
	})

	fields = append(fields, datasources.Field{
		Path: "date",
		Type: "date",
	})

	fields = append(fields, datasources.Field{
		Path: "address_tag",
		Type: "string",
	})

	fields = append(fields, datasources.Field{
		Path: "address_tag_link",
		Type: "string",
	})

	fields = append(fields, datasources.Field{
		Path: "input",
		Type: "string",
	})
	fields = append(fields, datasources.Field{
		Path: "output",
		Type: "string",
	})

	fields = append(fields, datasources.Field{
		Path: "value",
		Type: "int",
	})

	return
}
