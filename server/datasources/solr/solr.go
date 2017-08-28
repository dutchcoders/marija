package solr

import (
	"encoding/json"
	"fmt"

	"github.com/op/go-logging"

	"net/http"
	"net/url"

	"github.com/dutchcoders/marija/server/datasources"
)

type Response struct {
	FacetCounts struct {
		FacetFields struct {
		} `json:"facet_fields"`
		FacetHeatmaps struct {
		} `json:"facet_heatmaps"`
		FacetIntervals struct {
		} `json:"facet_intervals"`
		FacetQueries struct {
		} `json:"facet_queries"`
		FacetRanges struct {
		} `json:"facet_ranges"`
	} `json:"facet_counts"`
	Response struct {
		Docs     []map[string]interface{} `json:"docs"`
		NumFound int64                    `json:"numFound"`
		Start    int64                    `json:"start"`
	} `json:"response"`
	ResponseHeader struct {
		Params struct {
			Facet string `json:"facet"`
			Q     string `json:"q"`
			Rows  string `json:"rows"`
			Wt    string `json:"wt"`
		} `json:"params"`
		QTime       int64
		Status      int64 `json:"status"`
		ZkConnected bool  `json:"zkConnected"`
	} `json:"responseHeader"`
}

var log = logging.MustGetLogger("marija/datasources/elasticsearch")

type Solr struct {
	URL url.URL

	Username string
	Password string
}

func (m *Solr) UnmarshalTOML(p interface{}) error {
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

	_ = username
	_ = password

	if v, ok := data["url"]; !ok {
	} else if v, ok := v.(string); !ok {
	} else if u, err := url.Parse(v); err != nil {
		return err
	} else {
		m.URL = *u

		u.Path = ""

		/*
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
		*/
	}

	return nil
}

func (i *Solr) Search(so datasources.SearchOptions) ([]datasources.Item, int, error) {
	size := 200
	if so.Size > 0 {
		size = so.Size
	}

	rel, err := url.Parse(fmt.Sprintf("select?indent=on&q=%s&wt=json&start=%d&rows=%d", so.Query, so.From, size))
	if err != nil {
		return []datasources.Item{}, 0, err
	}

	u := i.URL.ResolveReference(rel)

	resp, err := http.DefaultClient.Get(u.String())
	if err != nil {
		return []datasources.Item{}, 0, err
	}

	response := Response{}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return []datasources.Item{}, 0, err
	}

	items := make([]datasources.Item, len(response.Response.Docs))
	for i, doc := range response.Response.Docs {

		items[i] = datasources.Item{
			ID:     doc["md5_name"].(string),
			Fields: doc,
		}
	}

	return items, int(response.Response.NumFound), nil
	/*
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
	*/
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

func (i *Solr) Fields() (fields []datasources.Field, err error) {
	/*
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
	*/
	fields = append(fields, datasources.Field{
		Path: "nlp_keywords",
		Type: "string",
	})

	fields = append(fields, datasources.Field{
		Path: "email_from_string",
		Type: "string",
	})

	fields = append(fields, datasources.Field{
		Path: "email_cc_string",
		Type: "string",
	})

	fields = append(fields, datasources.Field{
		Path: "type",
		Type: "string",
	})

	fields = append(fields, datasources.Field{
		Path: "email_to_string",
		Type: "string",
	})

	fields = append(fields, datasources.Field{
		Path: "email_bcc_string",
		Type: "string",
	})

	fields = append(fields, datasources.Field{
		Path: "file_name",
		Type: "string",
	})

	fields = append(fields, datasources.Field{
		Path: "GDPR_EMAIL",
		Type: "string",
	})

	fields = append(fields, datasources.Field{
		Path: "GDPR_NAME",
		Type: "string",
	})

	return
}
