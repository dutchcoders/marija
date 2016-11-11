package server

import (
	"encoding/json"
	"net/url"
	"path"

	"gopkg.in/olivere/elastic.v3"
)

type Elasticsearch3Index struct {
	client *elastic.Client
	u      *url.URL
}

func NewElasticsearchIndex(u *url.URL) (Index, error) {
	client, err := elastic.NewClient(elastic.SetURL(u.Host), elastic.SetSniff(false))
	if err != nil {
		return nil, err
	}

	return &Elasticsearch3Index{
		client: client,
		u:      u,
	}, nil
}

func (i *Elasticsearch3Index) Search(query string) ([]Item, error) {
	hl := elastic.NewHighlight()
	hl = hl.Fields(elastic.NewHighlighterField("*").RequireFieldMatch(false).NumOfFragments(0))
	hl = hl.PreTags("<em>").PostTags("</em>")

	q := elastic.NewQueryStringQuery(query)
	results, err := i.client.Search().
		Index(path.Base(i.u.Path)).
		Highlight(hl).
		Query(q).
		From(0).Size(500).
		Do()
	if err != nil {
		return nil, err
	}

	items := make([]Item, len(results.Hits.Hits))
	for i, hit := range results.Hits.Hits {
		var fields map[string]interface{}
		if err := json.Unmarshal(*hit.Source, &fields); err != nil {
			log.Error("Error unmarshalling source: %s", err.Error())
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
