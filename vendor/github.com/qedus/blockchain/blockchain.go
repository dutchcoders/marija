package blockchain

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"reflect"
)

const (
	rootURL = "https://blockchain.info"
)

var (
	// Used when an iterator has exhausted its supply.
	IterDone = errors.New("iterator done")
)

type BlockChain struct {
	client         *http.Client
	GUID           string
	Password       string
	SecondPassword string
	APICode        string
}

type Item interface {
	load(bc *BlockChain) error
}

func New(c *http.Client) *BlockChain {
	return &BlockChain{client: c}
}

func checkHTTPResponse(r *http.Response) error {
	if r.StatusCode == 200 {
		return nil
	}

	bodyErr, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}
	return fmt.Errorf("%s: %s: %.30q...",
		r.Request.URL, r.Status, bodyErr)
}

func (bc *BlockChain) Request(item Item) error {
	return item.load(bc)
}

func (bc *BlockChain) httpGetJSON(url string, v interface{}) error {
	resp, err := bc.client.Get(url)
	if err != nil {
		return err
	}

	if err := checkHTTPResponse(resp); err != nil {
		return err
	}

	defer resp.Body.Close()
	return decodeJSON(resp.Body, v)
}

func decodeJSON(r io.Reader, v interface{}) error {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(data, v); err != nil {
		return fmt.Errorf("%s with data %.30q...", err.Error(), data)
	}

	// Check for errors.
	errVal := reflect.ValueOf(v).Elem().FieldByName("Error")
	if errVal.IsValid() && errVal.String() != "" {
		return errors.New(errVal.String())
	}

	return nil
}
