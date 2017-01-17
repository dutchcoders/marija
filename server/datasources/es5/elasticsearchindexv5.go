package es5

import (
	"context"
	"encoding/json"

	"github.com/op/go-logging"
	log2 "log"

	"net/url"
	"os"
	"path"

	elastic "gopkg.in/olivere/elastic.v5"

	"fmt"
	"github.com/dutchcoders/marija/server/datasources"
)

var log = logging.MustGetLogger("marija/datasources/elasticsearch")

type Elasticsearch struct {
	client *elastic.Client
	URL    url.URL

	Username string
	Password string
}

func (m *Elasticsearch) UnmarshalTOML(p interface{}) error {
	data, _ := p.(map[string]interface{})

	username := ""
	if v, ok := data["username"]; !ok {
	} else if v, ok := v.(string); !ok {
	} else {
		username = v
	}

	password := ""
	if v, ok := data["password"]; !ok {
	} else if v, ok := v.(string); !ok {
	} else {
		password = v
	}

	if v, ok := data["url"]; !ok {
	} else if v, ok := v.(string); !ok {
	} else if u, err := url.Parse(v); err != nil {
		return err
	} else {
		m.URL = *u

		u.Path = ""

		errorlog := log2.New(os.Stdout, "APP ", log2.LstdFlags)
		params := []elastic.ClientOptionFunc{
			elastic.SetURL(u.String()),
			elastic.SetSniff(false),
			elastic.SetErrorLog(errorlog),
		}

		if username != "" {
			params = append(params, elastic.SetBasicAuth(username, password))
		}

		if client, err := elastic.NewClient(
			params...,
		); err != nil {
			return fmt.Errorf("Error connecting to: %s: %s", u.String(), err)
		} else {
			m.client = client
		}
	}

	return nil
}

func (i *Elasticsearch) Search(so datasources.SearchOptions) ([]datasources.Item, int, error) {
	hl := elastic.NewHighlight()
	hl = hl.Fields(elastic.NewHighlighterField("*").RequireFieldMatch(false).NumOfFragments(0))
	hl = hl.PreTags("<em>").PostTags("</em>")

	index := path.Base(i.URL.Path)

	q := elastic.NewQueryStringQuery(so.Query)
	results, err := i.client.Search().
		Index(index).
		Highlight(hl).
		Query(q).
		From(so.From).Size(so.Size).
		Do(context.Background())
	if err != nil {
		return nil, 0, err
	}

	log.Debugf("%#v\n", results)

	items := make([]datasources.Item, len(results.Hits.Hits))
	for i, hit := range results.Hits.Hits {
		var fields map[string]interface{}
		if err := json.Unmarshal(*hit.Source, &fields); err != nil {
			continue
		}

		items[i] = datasources.Item{
			ID:        hit.Id,
			Fields:    fields,
			Highlight: hit.Highlight,
		}
	}

	return items, int(results.Hits.TotalHits), nil
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

func (i *Elasticsearch) Fields() (fields []datasources.Field, err error) {
	index := path.Base(i.URL.Path)
	log.Debug("Using index: ", path.Base(i.URL.Path))

	mapping, err := i.client.GetMapping().
		Index(index).
		Do(context.Background())

	if err != nil {
		return nil, fmt.Errorf("Error retrieving fields for index: %s: %s", index, err.Error())
	}

	mapping = mapping[index].(map[string]interface{})
	mapping = mapping["mappings"].(map[string]interface{})
	for _, v := range mapping {
		// types
		fields = append(fields, flatten("", v.(map[string]interface{}))...)
	}

	return
}
