package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
)

type index struct {
	name string
	c    *client
}

type client struct {
	username string
	password string
	*http.Client
	baseURL *url.URL

	IPv4         *index
	Websites     *index
	Certificates *index
}

/*

type SearchIPv4Input struct {
	Query  string   `json:"query"`
	Page   int      `json:"page"`
	Fields []string `json:"fields"`
}

type SearchWebsiteInput struct {
	Query  string   `json:"query"`
	Page   int      `json:"page"`
	Fields []string `json:"fields"`
}

type SearchWebsiteOutput struct {
}

type SearchCertificateInput struct {
	Query  string   `json:"query"`
	Page   int      `json:"page"`
	Fields []string `json:"fields"`
}

type SearchCertificateOutput struct {
	Parsed struct {
		Extensions struct {
			AuthorityKeyId   string `json:"authority_key_id"`
			BasicConstraints struct {
				IsCa bool `json:"is_ca"`
			} `json:"basic_constraints"`
			CertificatePolicies []interface{} `json:"certificate_policies"`
			SubjectKeyId        string        `json:"subject_key_id"`
		} `json:"extensions"`
		FingerprintMd5    string `json:"fingerprint_md5"`
		FingerprintSha1   string `json:"fingerprint_sha1"`
		FingerprintSha256 string `json:"fingerprint_sha256"`
		Issuer            struct {
			CommonName []string `json:"common_name"`
		} `json:"issuer"`
		IssuerDn     string `json:"issuer_dn"`
		SerialNumber string `json:"serial_number"`
		Signature    struct {
			SelfSigned         bool `json:"self_signed"`
			SignatureAlgorithm struct {
				Name string `json:"name"`
				Oid  string `json:"oid"`
			} `json:"signature_algorithm"`
			Valid bool   `json:"valid"`
			Value string `json:"value"`
		} `json:"signature"`
		SignatureAlgorithm struct {
			Name string `json:"name"`
			Oid  string `json:"oid"`
		} `json:"signature_algorithm"`
		Subject struct {
			CommonName []string `json:"common_name"`
		} `json:"subject"`
		SubjectDn      string `json:"subject_dn"`
		SubjectKeyInfo struct {
			KeyAlgorithm struct {
				Name string `json:"name"`
				Oid  string `json:"oid"`
			} `json:"key_algorithm"`
			RsaPublicKey struct {
				Exponent int64  `json:"exponent"`
				Length   int64  `json:"length"`
				Modulus  string `json:"modulus"`
			} `json:"rsa_public_key"`
		} `json:"subject_key_info"`
		Validity struct {
			End   string `json:"end"`
			Start string `json:"start"`
		} `json:"validity"`
		Version int64 `json:"version"`
	} `json:"parsed"`
	Raw                 string `json:"raw"`
	UpdatedAt           string `json:"updated_at"`
	ValidNss            bool   `json:"valid_nss"`
	ValidationTimestamp string `json:"validation_timestamp"`
}

// GET /api/v1/view/:index/:id
type ViewRequest struct {
	// The search index the document is in. Must be one of ipv4, websites, or certificates.
	Index string `json:"index"`
	// The ID of the document you are requesting. In the IPv4 index, this is IP address (e.g., 192.168.1.1), domain in the websites index (e.g., google.com) and SHA-256 fingerprint in the certificates index (e.g., 9d3b51a6b80daf76e074730f19dc01e643ca0c3127d8f48be64cf3302f6622cc).
	ID string `json:"index"`
}

type Result struct {
	IP       string `json:"ip"`
	Location struct {
		Country []string `json:"location.country"`
	}
	AutonomousSystemASN []int `json:"autonomous_system.asn"`
}

func (r *Result) UnmarshalJSON(b []byte) error {
	var f interface{}
	json.Unmarshal(b, &f)

	v := f.(map[string]interface{})

	r.Location.Country = v["location.country"].([]string)
	r.IP = v["ip"].(string)
	r.AutonomousSystemASN = v["autonomous_system.asn"].([]int)
	return nil
}

type SearchIPv4Output struct {
	MetaData struct {
		Count       int    `json:"count"`
		Query       string `json:"query"`
		Page        int    `json:"page"`
		Pages       int    `json:"pages"`
		BackendTime int    `json:"backend_time"`
	} `json:"metadata"`
	Results []Result `json:"results"`
}
*/

type searchInput struct {
	Query  string   `json:"query"`
	Page   int      `json:"page"`
	Fields []string `json:"fields"`
}

type SearchOutput struct {
	MetaData struct {
		Count       int    `json:"count"`
		Query       string `json:"query"`
		Page        int    `json:"page"`
		Pages       int    `json:"pages"`
		BackendTime int    `json:"backend_time"`
	} `json:"metadata"`
	Results []map[string]interface{} `json:"results"`
	// todo(nl5887): Results json.RawMessage
}

type ViewOutput map[string]interface{}

// POST /api/v1/report/:index
// POST /api/v1/export
// GET /api/v1/export/:job_id
// GET /api/v1/data
// GET /api/v1/data/:series
// GET /api/v1/data/:series/:result

type SearchOption func(input searchInput) searchInput

func Page(i int) SearchOption {
	return func(input searchInput) searchInput {
		input.Page = i
		return input
	}
}

func Field(s string) SearchOption {
	return func(input searchInput) searchInput {
		input.Fields = append(input.Fields, s)
		return input
	}
}

func Query(s string) SearchOption {
	return func(input searchInput) searchInput {
		input.Query = s
		return input
	}
}

// The search endpoint allows searches against the IPv4, Alexa Top Million, and Certificates indexes using the same search syntax as the main site.
// The endpoint returns a paginated result of the most recent information we know for the set of user selected fields. More information
// about the returned hosts, websites, and certificates can be fetched using the /view endpoint.
func (i *index) Search(options ...SearchOption) (*SearchOutput, error) {
	input := searchInput{
		Query:  "",
		Page:   1,
		Fields: []string{},
	}

	for _, fn := range options {
		input = fn(input)
	}

	request, err := i.c.NewRequest("POST", fmt.Sprintf("/api/v1/search/%s", i.name), input)
	if err != nil {
		return nil, err
	}

	output := SearchOutput{}
	if err := i.c.Do(request, &output); err != nil {
		return nil, err
	}

	return &output, nil
}

// The ID of the document you are requesting. In the IPv4 index, this is IP address (e.g., 192.168.1.1), domain in the
// websites index (e.g., google.com) and SHA-256 fingerprint in the certificates index (e.g.,
// 9d3b51a6b80daf76e074730f19dc01e643ca0c3127d8f48be64cf3302f6622cc).
func (i *index) View(id string) (*ViewOutput, error) {
	request, err := i.c.NewRequest("GET", fmt.Sprintf("/api/v1/view/%s/%s", i.name, id), nil)
	if err != nil {
		return nil, err
	}

	output := ViewOutput{}
	if err := i.c.Do(request, &output); err != nil {
		return nil, err
	}

	return &output, nil
}

type reportInput struct {
	Query   string `json:"query"`
	Field   string `json:"field"`
	Buckets int    `json:"buckets"`
}

type reportOutput struct {
	MetaData struct {
		Count            int    `json:"count"`
		Query            string `json:"query"`
		NonNullCount     int    `json:"nonnull_count"`
		OtherResultCount int    `json:"other_result_count"`
		ErrorBound       int    `json:"error_bound"`
		Buckets          int    `json:"buckets"`
		BackendTime      int    `json:"backend_time"`
	} `json:"metadata"`
	Results []struct {
		Key   string `json:"key"`
		Count int    `json:"doc_count"`
	} `json:"results"`
}

// The build report endpoint lets you run aggregate reports on the breakdown of a field in a result set analogous to
// the "Build Report" functionality in the front end. For example, if you wanted to determine the breakdown of cipher
// suites selected by Top Million Websites. Data should be posted as a JSON request document.
func (i *index) Report() (*reportOutput, error) {
	input := reportInput{
		Query:   "80.http.get.headers.server: Apache",
		Field:   "location.country",
		Buckets: 100,
	}

	request, err := i.c.NewRequest("GET", fmt.Sprintf("/api/v1/report/%s", i.name), input)
	if err != nil {
		return nil, err
	}

	output := reportOutput{}
	if err := i.c.Do(request, &output); err != nil {
		return nil, err
	}

	return &output, nil
}

type queryInput struct {
	Query string `json:"query"`
}

type queryOutput struct {
	Configuration struct {
		Query string `json:"query"`
	} `json:"configuration"`
	JobId  string `json:"job_id"`
	Status string `json:"status"`
}

// POST /api/v1/query
// GET /api/v1/query/:job_id/:page
// GET /api/v1/query_definitions
func (c *client) Query(q string) (*queryOutput, error) {
	input := queryInput{
		Query: q,
	}

	request, err := c.NewRequest("POST", "/api/v1/query", input)
	if err != nil {
		return nil, err
	}

	output := queryOutput{}
	if err := c.Do(request, &output); err != nil {
		return nil, err
	}

	return &output, nil
}

func (c *client) NewRequest(method, urlStr string, body interface{}) (*http.Request, error) {
	rel, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	u := c.baseURL.ResolveReference(rel)

	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "text/json; charset=UTF-8")
	req.Header.Add("Accept", "text/json")

	req.SetBasicAuth(c.username, c.password)
	return req, nil
}

func New(username, password string) *client {
	if baseURL, err := url.Parse("https://censys.io"); err != nil {
		panic(err)
	} else {
		c := &client{
			username: username,
			password: password,
			baseURL:  baseURL,
			Client:   http.DefaultClient,
		}

		c.IPv4 = &index{
			name: "ipv4",
			c:    c,
		}

		c.Websites = &index{
			name: "websites",
			c:    c,
		}

		c.Certificates = &index{
			name: "certificates",
			c:    c,
		}

		return c
	}
}

func (wd *client) Do(req *http.Request, v interface{}) error {
	if dump, err := httputil.DumpRequestOut(req, true); err == nil {
		os.Stdout.Write(dump)
	}

	resp, err := wd.Client.Do(req)
	if err != nil {
		return err
	}

	r := resp.Body
	defer r.Close()

	if true {

	} else if dump, err := httputil.DumpResponse(resp, true); err == nil {
		os.Stdout.Write(dump)
	}

	r2 := r

	if resp.StatusCode != http.StatusOK {
		err := Error{}
		json.NewDecoder(r2).Decode(&err)
		return &err
	}

	err = json.NewDecoder(r2).Decode(&v)
	if err != nil {
		return err
	}

	return nil
}
