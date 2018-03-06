package splunk

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"net/http"
	"net/http/httputil"
	"net/url"
)

var (
	ErrNoContent = errors.New("No content")
)

var debug = false

type JobResponse struct {
	SID string `json:"sid`
}

type SummaryResponse struct {
	EventCount int `json:"event_count"`

	Fields map[string]struct {
		Name string `json:"name"`
	} `json:"fields"`
}

type ResultsResponse struct {
	Preview    bool            `json:"preview"`
	InitOffset int             `json:"init_offset"`
	Messages   json.RawMessage `json:"messages"`

	Fields []struct {
		Name string `json:"name"`
	} `json:"fields"`

	Results     []map[string]interface{} `json:"results"`
	Highlighted json.RawMessage          `json:"highlighted"`
}

type Client struct {
	*http.Client

	Username string
	Password string

	BaseURL *url.URL
}

func NewSplunkClient(baseURL url.URL) *Client {
	return &Client{
		Client: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			},
		},
		BaseURL: &baseURL,
	}
}

func (c *Client) NewRequest(method, urlStr string, body io.Reader) (*http.Request, error) {
	rel, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	u := c.BaseURL.ResolveReference(rel)

	log.Debug("%s", u.String())

	req, err := http.NewRequest(method, u.String(), body)
	if err != nil {
		return nil, err
	}

	return req, nil
}

func (c *Client) Do(req *http.Request, v interface{}) error {
	req.SetBasicAuth(c.Username, c.Password)

	req.Header.Set("Accept", "*/*")
	req.Header.Set("User-Agent", "Marija Splunk Connector")

	if debug {
		data, _ := httputil.DumpRequest(req, false)
		fmt.Println(string(data))
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		return err
	}

	if debug {
		data, _ := httputil.DumpResponse(resp, false)
		fmt.Println(string(data))
	}

	if resp.StatusCode == http.StatusNoContent {
		return ErrNoContent
	} else if resp.StatusCode >= http.StatusMultipleChoices {
		return fmt.Errorf("%s", http.StatusText(resp.StatusCode))
	}

	r := resp.Body
	return json.NewDecoder(r).Decode(v)
}
