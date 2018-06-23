package solr

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/op/go-logging"

	"net/http"
	"net/url"

	"github.com/dutchcoders/marija/server/datasources"
)

var (
	_ = datasources.Register("tronscan", New)
)

func New(options ...func(datasources.Index) error) (datasources.Index, error) {
	s := Tronscan{}

	for _, optionFn := range options {
		optionFn(&s)
	}

	return &s, nil
}

type Config struct {
	URL url.URL
}

type Response struct {
	Data []struct {
		Amount              interface{} `json:"amount"`
		Block               int64       `json:"block"`
		Confirmed           bool        `json:"confirmed"`
		ID                  string      `json:"id"`
		Timestamp           int64       `json:"timestamp"`
		TokenName           string      `json:"tokenName"`
		TransactionHash     string      `json:"transactionHash"`
		TransferFromAddress string      `json:"transferFromAddress"`
		TransferToAddress   string      `json:"transferToAddress"`
	} `json:"data"`
	Total int64 `json:"total"`
}

var log = logging.MustGetLogger("marija/datasources/tronscan")

type Tronscan struct {
	Config
}

func (m *Config) UnmarshalTOML(p interface{}) error {
	data, _ := p.(map[string]interface{})

	_ = data

	return nil
}

func (m *Tronscan) Type() string {
	return "tronscan"
}

func (b *Tronscan) Search(ctx context.Context, so datasources.SearchOptions) datasources.SearchResponse {
	itemCh := make(chan datasources.Item)
	errorCh := make(chan error)

	client := http.DefaultClient

	go func() {
		defer close(itemCh)
		defer close(errorCh)

		size := 50
		if so.Size > 0 {
			size = so.Size
		}

		q := url.Values{}
		q.Add("address", "7YxAaK71utTpYJ8u4Zna7muWxd1pQwimpGxy8")

		u, err := url.Parse(fmt.Sprintf("https://api.tronscan.org/api/transfer"))
		if err != nil {
			errorCh <- err
			return
		}

		u.RawQuery = q.Encode()

		count := 0

		for {
			req, err := http.NewRequest("GET", u.String(), nil)
			if err != nil {
				errorCh <- err
				return
			}

			resp, err := client.Do(req)
			if err != nil {
				errorCh <- err
				return
			}

			response := Response{}
			if err := json.NewDecoder(io.TeeReader(resp.Body, os.Stdout)).Decode(&response); err != nil {
				errorCh <- err
				return
			}

			for _, doc := range response.Data {
				if count > size {
					return
				}

				fields := map[string]interface{}{
					"id":                  doc.ID,
					"confirmed":           doc.Confirmed,
					"block":               doc.Block,
					"transactionHash":     doc.TransactionHash,
					"timestamp":           doc.Timestamp,
					"transferFromAddress": doc.TransferFromAddress,
					"transferToAddress":   doc.TransferToAddress,
					"amount":              fmt.Sprintf("%f TRX", doc.Amount),
					"tokenName":           doc.TokenName,
				}

				item := datasources.Item{
					ID:     doc.TransactionHash,
					Fields: fields,
				}

				select {
				case itemCh <- item:
				case <-ctx.Done():
					return
				}

				count++
			}
		}
	}()

	return datasources.NewSearchResponse(
		itemCh,
		errorCh,
	)
}

func flatten(root string, m map[string]interface{}) (fields []datasources.Field) {
	for k, v := range m {
		if k == "mappings" {
			fields = append(fields, flatten(root, v.(map[string]interface{}))...)
		} else if k == "properties" {
			fields = append(fields, flatten(root, v.(map[string]interface{}))...)
		} else if k == "fields" {
			fields = append(fields, flatten(root, v.(map[string]interface{}))...)
		} else {
			if v2, ok := v.(map[string]interface{}); ok {
				key := k
				if root != "" {
					key = root + "." + key
				}
				fields = append(fields, flatten(key, v2)...)
			} else if k == "type" {
				fields = append(fields, datasources.Field{
					Path: root,
					Type: v.(string),
				})
			}
		}
	}

	return
}

func (i *Tronscan) GetFields(context.Context) (fields []datasources.Field, err error) {
	fields = append(fields, datasources.Field{
		Path: "id",
		Type: "string",
	})

	fields = append(fields, datasources.Field{
		Path: "confirmed",
		Type: "bool",
	})

	fields = append(fields, datasources.Field{
		Path: "block",
		Type: "string",
	})

	fields = append(fields, datasources.Field{
		Path: "transactionHash",
		Type: "string",
	})

	fields = append(fields, datasources.Field{
		Path: "timestamp",
		Type: "date",
	})

	fields = append(fields, datasources.Field{
		Path: "transferFromAddress",
		Type: "string",
	})

	fields = append(fields, datasources.Field{
		Path: "transferToAddress",
		Type: "string",
	})

	fields = append(fields, datasources.Field{
		Path: "amount",
		Type: "string",
	})

	fields = append(fields, datasources.Field{
		Path: "tokenName",
		Type: "string",
	})

	return
}
