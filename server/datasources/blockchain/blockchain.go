package btc

import (
	"net/http"
	"net/url"

	"fmt"
	"github.com/dutchcoders/marija/server/datasources"
	"github.com/qedus/blockchain"
	"time"
)

func (m *BTC) UnmarshalTOML(p interface{}) error {
	data, _ := p.(map[string]interface{})

	_ = data

	m.client = blockchain.New(http.DefaultClient)
	return nil
}

type BTC struct {
	client *blockchain.BlockChain
	u      *url.URL
}

func (b *BTC) Search(so datasources.SearchOptions) ([]datasources.Item, int, error) {
	address := &blockchain.Address{Address: so.Query}
	if err := b.client.Request(address); err != nil {
		return nil, 0, nil
	}

	i := 0

	items := []datasources.Item{}

main:
	for {

		tx, err := address.NextTransaction()
		if err == blockchain.IterDone {
			break
		} else if err != nil {
			return nil, 0, err
		}

		for _, input := range tx.Inputs {
			for _, output := range tx.Outputs {
				if i < so.From {
					i++
					continue
				}

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

				items = append(items, item)

				if i >= so.From+so.Size {
					break main
				}

				i++
			}
		}
	}

	return items, len(items), nil
}

func (i *BTC) Fields() (fields []datasources.Field, err error) {
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
