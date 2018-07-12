package splunk

import (
	"context"
	"fmt"
	"strings"
	"time"

	"net/url"

	"github.com/dutchcoders/marija/server/datasources"
	logging "github.com/op/go-logging"
)

var log = logging.MustGetLogger("marija/datasources/elasticsearch")

var (
	_ = datasources.Register("splunk", New)
)

func New(options ...func(datasources.Index) error) (datasources.Index, error) {
	s := Splunk{}

	for _, optionFn := range options {
		optionFn(&s)
	}

	client := NewSplunkClient(s.URL)

	client.Username = s.Username
	client.Password = s.Password

	s.client = client

	return &s, nil
}

type Config struct {
	URL url.URL

	Username string
	Password string
}

type Splunk struct {
	Config

	client *Client
}

func (m *Splunk) Type() string {
	return "splunk"
}

func (m *Splunk) UnmarshalTOML(p interface{}) error {
	data, _ := p.(map[string]interface{})

	if v, ok := data["username"]; !ok {
	} else if v, ok := v.(string); !ok {
	} else {
		m.Username = v
	}

	if v, ok := data["password"]; !ok {
	} else if v, ok := v.(string); !ok {
	} else {
		m.Password = v
	}

	if v, ok := data["url"]; !ok {
	} else if v, ok := v.(string); !ok {
	} else if u, err := url.Parse(v); err != nil {
		return err
	} else {
		m.URL = *u

		u.Path = ""
	}

	return nil
}

func (i *Splunk) Search(ctx context.Context, so datasources.SearchOptions) datasources.SearchResponse {
	itemCh := make(chan datasources.Item)
	errorCh := make(chan error)

	go func() {
		defer close(itemCh)
		defer close(errorCh)

		size := 200
		if so.Size > 0 {
			size = so.Size
		}

		_ = size

		data := url.Values{}
		data.Add("output_mode", "json")
		data.Add("rf", "*")
		data.Add("search", fmt.Sprintf("search %s", so.Query))

		req, err := i.client.NewRequest("POST", "/services/search/jobs", strings.NewReader(data.Encode()))
		if err != nil {
			errorCh <- err
			return
		}

		response := JobResponse{}

		if err := i.client.Do(req, &response); err != nil {
			errorCh <- err
			return
		}

		sid := response.SID

		hits := make(chan map[string]interface{})

		go func() {
			defer close(hits)

			offset := 0

			for {
				data = url.Values{}
				data.Add("output_mode", "json")
				data.Add("count", "1000")
				data.Add("offset", fmt.Sprintf("%d", offset))

				req, err = i.client.NewRequest("GET", fmt.Sprintf("/services/search/jobs/%s/results/?%s", sid, data.Encode()), nil)
				if err != nil {
					errorCh <- err
					return
				}

				rr := ResultsResponse{}
				if err := i.client.Do(req, &rr); err == ErrNoContent {
					time.Sleep(time.Second * 1)
					continue
				} else if err != nil {
					errorCh <- err
					return
				}

				if len(rr.Results) == 0 {
					break
				}

				for _, hit := range rr.Results {
					select {
					case hits <- hit:
					case <-ctx.Done():
						return
					}
				}

				offset += len(rr.Results)
			}
		}()

		uniqueFields := map[string]bool{}

		for {
			select {
			case hit, ok := <-hits:
				if !ok {
					return
				}

				fields := flattenFields("", hit)

				for key, val := range fields {
					fields[key] = val

					if _, ok := uniqueFields[key]; ok {
						continue
					}

					uniqueFields[key] = true
				}

				itemCh <- datasources.Item{
					ID:        "", // fields["_bkt"].(string),
					Fields:    fields,
					Highlight: nil,
				}
			}
		}

		return
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

func (i *Splunk) GetFields(context.Context) (fields []datasources.Field, err error) {
	data := url.Values{}
	data.Add("output_mode", "json")
	data.Add("preview", "true")
	data.Add("auto_cancel", "30")
	data.Add("rf", "*")
	data.Add("status_buckets", "300")
	data.Add("sample_ratio", "1")
	data.Add("search", "search *")

	req, err := i.client.NewRequest("POST", "/services/search/jobs", strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}

	response := JobResponse{}
	if err := i.client.Do(req, &response); err != nil {
		return nil, err
	}

	sid := response.SID

	data = url.Values{}
	data.Add("output_mode", "json")
	data.Add("min_freq", "0")

	for {

		req, err = i.client.NewRequest("GET", fmt.Sprintf("/services/search/jobs/%s/summary/?%s", sid, data.Encode()), nil)
		if err != nil {
			return nil, err
		}

		sr := SummaryResponse{}

		if err := i.client.Do(req, &sr); err == ErrNoContent {
			time.Sleep(time.Second * 1)
			continue
		} else if err != nil {
			return nil, err
		}

		for name, _ := range sr.Fields {
			fields = append(fields, datasources.Field{
				Path: name,
				Type: "string",
			})
		}

		break
	}

	return
}
