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
	_ = datasources.Register("openkvk", New)
)

func New(options ...func(datasources.Index) error) (datasources.Index, error) {
	s := OpenKvK{}

	for _, optionFn := range options {
		optionFn(&s)
	}

	return &s, nil
}

type Config struct {
	URL url.URL

	APIKey string
}

type Response struct {
	Embedded struct {
		Bedrijf []struct {
			BTW                   string   `json:"BTW"`
			LEI                   string   `json:"LEI"`
			RSIN                  string   `json:"RSIN"`
			VboID                 string   `json:"vbo_id"`
			BestaandeHandelsnaam  []string `json:"bestaandehandelsnaam"`
			StatutaireHandelsnaam []string `json:"statutairehandelsnaam"`
			PandID                string   `json:"pand_id"`
			Dossiernummer         string   `json:"dossiernummer"`
			Handelsnaam           string   `json:"handelsnaam"`
			Links                 struct {
				Self struct {
					Href string `json:"href"`
				} `json:"self"`
			} `json:"_links"`
			Postcode string `json:"postcode"`
			Locatie  struct {
				Lat string `json:"lat"`
				Lon string `json:"lon"`
			} `json:"locatie"`
			Straat           string `json:"straat"`
			Plaats           string `json:"plaats"`
			Subdossiernummer string `json:"subdossiernummer"`
			Vestigingsnummer string `json:"vestigingsnummer"`
		} `json:"bedrijf"`
	} `json:"_embedded"`
	Links struct {
		First struct {
			Href string `json:"href"`
		} `json:"first"`
		Last struct {
			Href string `json:"href"`
		} `json:"last"`
		Next struct {
			Href string `json:"href"`
		} `json:"next"`
		Self struct {
			Href string `json:"href"`
		} `json:"self"`
	} `json:"_links"`
	PageCount      int64 `json:"pageCount"`
	Size           int64 `json:"size"`
	TotalItemCount int64 `json:"totalItemCount"`
}

var log = logging.MustGetLogger("marija/datasources/openkvk")

type OpenKvK struct {
	Config
}

func (m *Config) UnmarshalTOML(p interface{}) error {
	data, _ := p.(map[string]interface{})

	if v, ok := data["api_key"]; !ok {
	} else if v, ok := v.(string); !ok {
	} else {
		m.APIKey = v
	}

	return nil
}

func (m *OpenKvK) Type() string {
	return "openkvk"
}

func (b *OpenKvK) Search(ctx context.Context, so datasources.SearchOptions) datasources.SearchResponse {
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
		q.Add("query", so.Query)
		for _, name := range []string{
			"statutairehandelsnaam",
			"bestaandehandelsnaam",
			"postcode",
			"straat",
			"handelsnaam",
			"locatie",
			"vestigingsnummer",
			"dossiernummer",
			"btw",
			"rsin",
			"lei",
			"pand_id",
			"vbo_id",
		} {
			q.Add("fields[]", name)
		}

		u, err := url.Parse(fmt.Sprintf("https://api.overheid.io/openkvk"))
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

			req.Header.Add("ovio-api-key", b.APIKey)

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

			for _, doc := range response.Embedded.Bedrijf {
				if count > size {
					return
				}

				fields := map[string]interface{}{
					"dossiernummer":         doc.Dossiernummer,
					"handelsnaam":           doc.Handelsnaam,
					"subdossiernummer":      doc.Subdossiernummer,
					"vestigingsnummer":      doc.Vestigingsnummer,
					"btw":                   doc.BTW,
					"lei":                   doc.LEI,
					"rsin":                  doc.RSIN,
					"bestaandehandelsnaam":  doc.BestaandeHandelsnaam,
					"statutairehandelsnaam": doc.BestaandeHandelsnaam,
					"pand_id":               doc.PandID,
					"vbo_id":                doc.VboID,
					"postcode":              doc.Postcode,
					"straat":                doc.Straat,
					"plaats":                doc.Plaats,
				}

				if doc.Locatie.Lon == "" {
				} else if doc.Locatie.Lat == "" {
				} else {
					fields["locatie"] = fmt.Sprintf("%s,%s", doc.Locatie.Lat, doc.Locatie.Lon)
				}

				item := datasources.Item{
					ID:     doc.Dossiernummer,
					Fields: fields,
				}

				select {
				case itemCh <- item:
				case <-ctx.Done():
					return
				}

				count++
			}

			if response.Links.Next.Href == "" {
				return
			}

			u, err = url.Parse(response.Links.Next.Href)
			if err != nil {
				errorCh <- err
				return
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

func (i *OpenKvK) GetFields(context.Context) (fields []datasources.Field, err error) {
	fields = append(fields, datasources.Field{
		Path: "btw",
		Type: "string",
	})
	fields = append(fields, datasources.Field{
		Path: "lei",
		Type: "string",
	})

	fields = append(fields, datasources.Field{
		Path: "rsin",
		Type: "string",
	})

	fields = append(fields, datasources.Field{
		Path: "bestaandehandelsnaam",
		Type: "string",
	})
	fields = append(fields, datasources.Field{
		Path: "statutairehandelsnaam",
		Type: "string",
	})

	fields = append(fields, datasources.Field{
		Path: "vbo_id",
		Type: "string",
	})

	fields = append(fields, datasources.Field{
		Path: "pand_id",
		Type: "string",
	})

	fields = append(fields, datasources.Field{
		Path: "dossiernummer",
		Type: "string",
	})

	fields = append(fields, datasources.Field{
		Path: "locatie",
		Type: "location",
	})

	fields = append(fields, datasources.Field{
		Path: "handelsnaam",
		Type: "string",
	})

	fields = append(fields, datasources.Field{
		Path: "postcode",
		Type: "string",
	})

	fields = append(fields, datasources.Field{
		Path: "subdossiernummer",
		Type: "string",
	})

	fields = append(fields, datasources.Field{
		Path: "vestigingsnummer",
		Type: "string",
	})

	fields = append(fields, datasources.Field{
		Path: "straat",
		Type: "string",
	})

	fields = append(fields, datasources.Field{
		Path: "plaats",
		Type: "string",
	})

	return
}
