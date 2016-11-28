package datasources

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/dutchcoders/marija/server/datasources"
	"github.com/dutchcoders/marija/server/datasources/es2"
	"github.com/dutchcoders/marija/server/datasources/es5"
)

var (
	_ = datasources.Register("es", NewElasticsearchIndex)
)

func NewElasticsearchIndex(u *url.URL) (datasources.Index, error) {
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
		return es5.NewElasticsearchIndexV5(u)
	} else {
		return es2.NewElasticsearchIndexV3(u)
	}
}
