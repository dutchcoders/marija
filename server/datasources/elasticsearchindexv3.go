package datasources

import (
	"encoding/json"
	"net/url"
	"path"
	"strings"

	elastic "gopkg.in/olivere/elastic.v3"
)

func NewElasticsearchIndexV3(u *url.URL) (Index, error) {
	client, err := elastic.NewClient(elastic.SetURL(u.Host), elastic.SetSniff(false))
	if err != nil {
		return nil, err
	}

	// check version here, and return appropriate version
	return &ElasticsearchIndexV3{
		client: client,
		u:      u,
	}, nil
}

type ElasticsearchIndexV3 struct {
	client *elastic.Client
	u      *url.URL
}

func (i *ElasticsearchIndexV3) Search(so SearchOptions) ([]Item, error) {
	hl := elastic.NewHighlight()
	hl = hl.Fields(elastic.NewHighlighterField("*").RequireFieldMatch(false).NumOfFragments(0))
	hl = hl.PreTags("<em>").PostTags("</em>")

	q := elastic.NewQueryStringQuery(so.Query)
	results, err := i.client.Search().
		Index(path.Base(i.u.Path)).
		Highlight(hl).
		Query(q).
		From(so.From).Size(so.Size).
		Do()
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

func (i *ElasticsearchIndexV3) Indices() ([]string, error) {
	stats, err := i.client.IndexStats().
		Metric("index").
		Do()
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

func (i *ElasticsearchIndexV3) Fields(index string) (fields []Field, err error) {
	mapping, err := i.client.GetMapping().
		Index(index).
		Do()

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
