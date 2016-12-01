package es5

import (
	"context"
	"encoding/json"
	log2 "log"
	"net/url"
	"os"
	"path"
	"strings"

	elastic "gopkg.in/olivere/elastic.v5"

	"github.com/dutchcoders/marija/server/datasources"
)

var (
	_ = datasources.Register("es", NewElasticsearchIndexV5)
)

func NewElasticsearchIndexV5(u *url.URL) (datasources.Index, error) {
	errorlog := log2.New(os.Stdout, "APP ", log2.LstdFlags)

	u2 := *u
	u2.Path = ""

	client, err := elastic.NewClient(elastic.SetURL(u2.String()), elastic.SetSniff(false), elastic.SetErrorLog(errorlog))
	if err != nil {
		return nil, err
	}

	// check version here, and return appropriate version
	return &ElasticsearchIndexV5{
		client: client,
		u:      u,
	}, nil
}

type ElasticsearchIndexV5 struct {
	client *elastic.Client
	u      *url.URL
}

func (i *ElasticsearchIndexV5) Search(so datasources.SearchOptions) ([]datasources.Item, error) {
	hl := elastic.NewHighlight()
	hl = hl.Fields(elastic.NewHighlighterField("*").RequireFieldMatch(false).NumOfFragments(0))
	hl = hl.PreTags("<em>").PostTags("</em>")

	q := elastic.NewQueryStringQuery(so.Query)
	results, err := i.client.Search().
		Index(path.Base(i.u.Path)).
		Highlight(hl).
		Query(q).
		From(so.From).Size(so.Size).
		Do(context.Background())
	if err != nil {
		return nil, err
	}

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

	return items, nil
}

func (i *ElasticsearchIndexV5) Indices() ([]string, error) {
	stats, err := i.client.IndexStats().
		Metric("index").
		Do(context.Background())
	if err != nil {
		return nil, err
	}

	indices := []string{}
	for k, _ := range stats.Indices {
		if strings.HasPrefix(k, ".") {
			continue
		}

		// todo(nl5887):
		// should create index url here, and we don't want to return strings,
		// but index object with name and url
		indices = append(indices, k)
	}

	return indices, nil
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

func (i *ElasticsearchIndexV5) Fields(index string) (fields []datasources.Field, err error) {
	mapping, err := i.client.GetMapping().
		Index(index).
		Do(context.Background())

	if err != nil {
		return nil, err
	}

	mapping = mapping[index].(map[string]interface{})
	mapping = mapping["mappings"].(map[string]interface{})
	for _, v := range mapping {
		// types
		fields = append(fields, flatten("", v.(map[string]interface{}))...)
	}

	return
}
