package voertuiggegevens

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
	_ = datasources.Register("voertuiggegevens", New)
)

func New(options ...func(datasources.Index) error) (datasources.Index, error) {
	s := VoertuigGegevens{}

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
		Kenteken []struct {
			DatumEersteAfgifteNederland string `json:"datumeersteafgiftenederland"`
			VervalDatumAPK              string `json:"vervaldatumapk"`
			EersteKleur                 string `json:"eerstekleur"`
			HandelsBenaming             string `json:"handelsbenaming"`
			CatalogusPrijs              string `json:"catalogusprijs"`
			Kenteken                    string `json:"kenteken"`
			Merk                        string `json:"merk"`
			VoertuigSoort               string `json:"voertuigsoort"`
		} `json:"kenteken"`
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

type VoertuigGegevens struct {
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

func (m *VoertuigGegevens) Type() string {
	return "openkvk"
}

func (b *VoertuigGegevens) Search(ctx context.Context, so datasources.SearchOptions) datasources.SearchResponse {
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
			"merk",
			"datumeersteafgiftenederland",
			"vervaldatumapk",
			"catalogusprijs",
			"eerstekleur",
			"kenteken",
			"handelsbenaming",
		} {
			q.Add("fields[]", name)
		}

		u, err := url.Parse(fmt.Sprintf("https://overheid.io/api/voertuiggegevens"))
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

			if resp.StatusCode != http.StatusOK {
				errorCh <- fmt.Errorf("Invalid status: %d", resp.StatusCode)
				return
			}

			response := Response{}
			if err := json.NewDecoder(io.TeeReader(resp.Body, os.Stdout)).Decode(&response); err != nil {
				errorCh <- err
				return
			}

			for _, doc := range response.Embedded.Kenteken {
				if count > size {
					return
				}

				fields := map[string]interface{}{
					"datumeersteafgiftenederland": doc.DatumEersteAfgifteNederland,
					"vervaldatumapk":              doc.VervalDatumAPK,
					"merk":                        doc.Merk,
					"kenteken":                    doc.Kenteken,
					"handelsbenaming":             doc.HandelsBenaming,
					"eerstekleur":                 doc.EersteKleur,
					"catalogusprijs":              doc.CatalogusPrijs,
				}

				item := datasources.Item{
					ID:     doc.Kenteken,
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

func (i *VoertuigGegevens) GetFields(context.Context) (fields []datasources.Field, err error) {
	fields = append(fields, datasources.Field{
		Path: "datumeersteafgiftenederland",
		Type: "string",
	})

	fields = append(fields, datasources.Field{
		Path: "vervaldatumapk",
		Type: "string",
	})

	fields = append(fields, datasources.Field{
		Path: "eerstekleur",
		Type: "string",
	})

	fields = append(fields, datasources.Field{
		Path: "handelsbenaming",
		Type: "string",
	})

	fields = append(fields, datasources.Field{
		Path: "kenteken",
		Type: "string",
	})

	fields = append(fields, datasources.Field{
		Path: "merk",
		Type: "string",
	})

	fields = append(fields, datasources.Field{
		Path: "voertuigsoort",
		Type: "string",
	})

	return
}
