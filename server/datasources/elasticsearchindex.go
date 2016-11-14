package datasources

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"strings"
)

func NewElasticsearchIndex(u *url.URL) (Index, error) {
	rel, err := url.Parse("/")
	if err != nil {
		log.Fatal(err)
	}

	u2 := u.ResolveReference(rel)

	resp, err := http.Get(u2.String())
	if err != nil {
		return nil, err
	}

	v := struct {
		Name    string `json:"name"`
		Version struct {
			Number string `json:"number"`
		} `json:"version"`
	}{}

	if err := json.NewDecoder(resp.Body).Decode(&v); err != nil {
		return nil, err
	}

	if strings.HasPrefix(v.Version.Number, "5.") {
		return NewElasticsearchIndexV5(u)
	} else {
		return NewElasticsearchIndexV3(u)
	}
}
