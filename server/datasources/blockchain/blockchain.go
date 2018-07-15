package btc

import (
	"context"
	"net/http"
	"net/url"
	"strings"

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
		defer close(itemCh)
		defer close(errorCh)

		qry := strings.Replace(so.Query, "\"", "", -1)

		sendTx := func(tx blockchain.Transaction) {
			inputs := []string{}
			sumOfInput := float64(0)

			for _, input := range tx.Inputs {
				inputs = append(inputs, input.PrevOut.Address)
				sumOfInput += float64(input.PrevOut.Value)
			}

			fields := map[string]interface{}{
				"relayed_by":   tx.RelayedBy,
				"label":        fmt.Sprintf("sum(in): %f", sumOfInput/100000000),
				"address":      inputs,
				"type":         "input",
				"hash":         tx.Hash,
				"size":         tx.Size,
				"block_height": tx.BlockHeight,
				"version":      tx.Version,
				"date":         time.Unix(tx.Time, 0),
				"input":        inputs,
				"value":        fmt.Sprintf("%f", sumOfInput/100000000),
			}

			item := datasources.Item{
				ID:     fmt.Sprintf("input.%s", tx.Hash),
				Fields: fields,
			}

			select {
			case itemCh <- item:
			case <-ctx.Done():
				return
			}

			for _, output := range tx.Outputs {
				fields := map[string]interface{}{
					"relayed_by":       tx.RelayedBy,
					"label":            fmt.Sprintf("out (%d): %f", output.Number, float64(output.Value)/100000000),
					"address":          []string{output.Address},
					"number":           output.Number,
					"type":             "output",
					"tx_index":         output.TransactionIndex,
					"size":             tx.Size,
					"block_height":     tx.BlockHeight,
					"version":          tx.Version,
					"hash":             tx.Hash,
					"date":             time.Unix(tx.Time, 0),
					"output":           []string{output.Address},
					"value":            fmt.Sprintf("%f", float64(output.Value)/100000000),
					"address_tag":      output.AddressTag,
					"address_tag_link": output.AddressTagLink,
				}

				item := datasources.Item{
					ID:     fmt.Sprintf("output.%s.%s.%d", tx.Hash, output.Address, output.TransactionIndex),
					Fields: fields,
				}

				select {
				case itemCh <- item:
				case <-ctx.Done():
					return
				}
			}
		}

		if len(qry) == 64 {
			tx := blockchain.Transaction{
				Hash: qry,
			}

			if err := b.client.Request(&tx); err != nil {
				errorCh <- err
				return
			}

			sendTx(tx)
		} else if len(qry) == 34 {
			item := blockchain.Address{Address: qry}
			if err := b.client.Request(&item); err != nil {
				errorCh <- err
				return
			}

			for {
				tx, err := item.NextTransaction()
				if err == blockchain.IterDone {
					break
				} else if err != nil {
					errorCh <- err
					return
				}

				sendTx(tx)
			}

			return
		}

	}()

	return datasources.NewSearchResponse(
		itemCh,
		errorCh,
	)
}

func (i *BTC) GetFields(context.Context) (fields []datasources.Field, err error) {
	fields = append(fields, datasources.Field{
		Path: "number",
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
		Path: "label",
		Type: "string",
	})

	fields = append(fields, datasources.Field{
		Path: "hash",
		Type: "string",
	})

	fields = append(fields, datasources.Field{
		Path: "size",
		Type: "string",
	})

	fields = append(fields, datasources.Field{
		Path: "block_height",
		Type: "string",
	})

	fields = append(fields, datasources.Field{
		Path: "version",
		Type: "string",
	})

	fields = append(fields, datasources.Field{
		Path: "tx_index",
		Type: "string",
	})

	fields = append(fields, datasources.Field{
		Path: "address",
		Type: "string",
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
		Type: "string",
	})

	return
}
