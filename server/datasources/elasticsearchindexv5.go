package datasources

import (
	"context"
	"encoding/json"
	"fmt"
	log2 "log"
	"net/url"
	"os"
	"path"
	"strings"

	elastic "gopkg.in/olivere/elastic.v5"
)

func NewElasticsearchIndexV5(u *url.URL) (Index, error) {
	fmt.Println("Elasticsearch v5", u.String())

	errorlog := log2.New(os.Stdout, "APP ", log2.LstdFlags)

	u2 := *u
	u2.Path = ""

	client, err := elastic.NewClient(elastic.SetURL(u2.String()), elastic.SetSniff(true), elastic.SetErrorLog(errorlog))
	if err != nil {
		return nil, err
	}

	// Ping the Elasticsearch server to get e.g. the version number
	info, code, err := client.Ping(u.String()).Do(context.Background())
	if err != nil {
		// Handle error
		panic(err)
	}

	fmt.Printf("Elasticsearch returned with code %d and version %s", code, info.Version.Number)

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

func (i *ElasticsearchIndexV5) Search(so SearchOptions) ([]Item, error) {
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

	items := make([]Item, len(results.Hits.Hits))
	for i, hit := range results.Hits.Hits {
		var fields map[string]interface{}
		if err := json.Unmarshal(*hit.Source, &fields); err != nil {
			continue
		}

		items[i] = Item{
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

/*
func flatten(root string, m map[string]interface{}) (fields []Field) {
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
				fields = append(fields, Field{
					Path: root,
					Type: v.(string),
				})
			}
		}
	}

	return
}
*/

func (i *ElasticsearchIndexV5) Fields(index string) (fields []Field, err error) {
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
