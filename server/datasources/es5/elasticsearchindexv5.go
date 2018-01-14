package es5

import (
	"context"
	"encoding/json"
	"io"

	log2 "log"

	logging "github.com/op/go-logging"
	cache "github.com/patrickmn/go-cache"

	"net/url"
	"path"

	elastic "gopkg.in/olivere/elastic.v5"

	"fmt"

	"os"

	"github.com/dutchcoders/marija/server/datasources"
)

// implement via
// return partial results to websocket
// scripted fields from config

var log = logging.MustGetLogger("marija/datasources/elasticsearch")

type Elasticsearch struct {
	client *elastic.Client
	URL    url.URL

	cache *cache.Cache

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

// return chan datasources.Item instead
type SearchResponse struct {
	itemChan  chan datasources.Item
	errorChan chan error
}

func (sr *SearchResponse) Item() chan datasources.Item {
	return sr.itemChan
}

func (sr *SearchResponse) Error() chan error {
	return sr.errorChan
}

func (i *Elasticsearch) Search(ctx context.Context, so datasources.SearchOptions) datasources.SearchResponse {
	itemCh := make(chan datasources.Item)
	errorCh := make(chan error)

	go func() {
		defer close(itemCh)
		defer close(errorCh)

		hl := elastic.NewHighlight()
		hl = hl.Fields(elastic.NewHighlighterField("*").RequireFieldMatch(false).NumOfFragments(0))
		hl = hl.PreTags("<em>").PostTags("</em>")

		index := path.Base(i.URL.Path)

		scriptFields := []*elastic.ScriptField{
		/*
			//		elastic.NewScriptField("src-ip_dst-ip_port", elastic.NewScript("params['_source']['source-ip'] + '_' + params['_source']['destination-ip'] + '_' + params['_source']['destination-port']")),
			elastic.NewScriptField("src-ip_dst-net_port", elastic.NewScript("params['_source']['source-ip'] + '_' + params['_source']['destination-net'] + '_' + params['_source']['destination-port']")),
		*/
		}

		q := elastic.NewQueryStringQuery(so.Query)

		src := elastic.NewSearchSource().
			Query(q).
			FetchSource(false).
			Highlight(hl).
			FetchSource(true).
			From(so.From).
			Size(100)

		if len(scriptFields) > 0 {
			src = src.ScriptFields(scriptFields...)
		}

		hits := make(chan *elastic.SearchHit)

		go func() {
			defer close(hits)

			scroll := i.client.Scroll().Index(index).SearchSource(src)
			for {
				results, err := scroll.Do(ctx)
				if err == io.EOF {
					return
				} else if err != nil {
					errorCh <- err
					return
				}

				for _, hit := range results.Hits.Hits {
					select {
					case hits <- hit:
					case <-ctx.Done():
						log.Debug("Search canceled query=%s, index=%s", q, index)
						return
					}
				}
			}

		}()

		for {
			select {
			case hit, ok := <-hits:
				if !ok {
					return
				}

				var fields map[string]interface{}
				if err := json.Unmarshal(*hit.Source, &fields); err != nil {
					errorCh <- err
					continue
				}

				for key, val := range hit.Fields {
					fields[key] = val
				}

				itemCh <- datasources.Item{
					ID:        hit.Id,
					Fields:    fields,
					Highlight: hit.Highlight,
				}
			}
		}

		return
	}()

	return &SearchResponse{
		itemCh,
		errorCh,
	}
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

func (i *Elasticsearch) GetFields(ctx context.Context) (fields []datasources.Field, err error) {
	index := path.Base(i.URL.Path)

	exists, err := i.client.IndexExists(index).Do(ctx)
	if err != nil {
		return nil, err
	} else if !exists {
		return nil, fmt.Errorf("Index %s doesn't exist", index)
	}

	mapping, err := i.client.GetMapping().
		Index(index).
		Do(ctx)

	if err != nil {
		return nil, fmt.Errorf("Error retrieving fields for index: %s: %s", index, err.Error())
	}

	mapping = mapping[index].(map[string]interface{})
	mapping = mapping["mappings"].(map[string]interface{})
	for _, v := range mapping {
		fields = append(fields, flatten("", v.(map[string]interface{}))...)
	}

	/*
		fields = append(fields, datasources.Field{
			Path: "src-ip_dst-net_port",
			Type: "string",
		})

		fields = append(fields, datasources.Field{
			Path: "src-ip_dst-ip_port",
			Type: "string",
		})
	*/

	return
}
