package server

import (
	"encoding/json"
	_ "log"

	"github.com/dutchcoders/marija/server/datasources"

	_ "github.com/dutchcoders/marija/server/datasources/blockchain"
	_ "github.com/dutchcoders/marija/server/datasources/es5"
	_ "github.com/dutchcoders/marija/server/datasources/twitter"
)

type Datasource struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type InitialStateMessage struct {
	Datasources []Datasource
	Version     string
	CommitID    string
}

func (em *InitialStateMessage) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Type        string       `json:"type"`
		Datasources []Datasource `json:"datasources"`
		Version     string       `json:"version"`
		CommitID    string       `json:"commit-id"`
	}{
		Type:        InitialStateReceive,
		Datasources: em.Datasources,
		Version:     em.Version,
		CommitID:    em.CommitID,
	})
}

type Request struct {
	RequestID string `json:"request-id"`
	Type      string `json:"type"`
}

type SearchRequest struct {
	Request

	Datasources []string `json:"datasources"`
	Fields      []string `json:"fields"`
	Query       string   `json:"query"`
}

type SearchResponse struct {
	RequestID string
	Server    string
	Query     string
	Nodes     []datasources.Item
}

func (em *SearchResponse) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Type      string             `json:"type"`
		RequestID string             `json:"request-id"`
		Server    string             `json:"server,omitempty"`
		Query     string             `json:"query"`
		Nodes     []datasources.Item `json:"results"`
	}{
		Type:      ActionTypeItemsReceive,
		RequestID: em.RequestID,
		Server:    em.Server,
		Query:     em.Query,
		Nodes:     em.Nodes,
	})
}

type GetFieldsRequest struct {
	Request

	Datasources []string `json:"datasources"`
}

type GetFieldsResponse struct {
	RequestID string
	Server    string
	Index     string
	Fields    interface{}
}

func (em *GetFieldsResponse) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Type      string      `json:"type"`
		RequestID string      `json:"request-id"`
		Server    string      `json:"server,omitempty"`
		Index     string      `json:"index,omitempty"`
		Fields    interface{} `json:"fields"`
	}{
		Type:      ActionTypeFieldsReceive,
		RequestID: em.RequestID,
		Server:    em.Server,
		Index:     em.Index,
		Fields:    em.Fields,
	})
}

type ErrorMessage struct {
	RequestID string
	Message   string
}

func (em *ErrorMessage) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Type      string `json:"type"`
		RequestID string `json:"request-id"`
		Message   string `json:"message"`
	}{
		Type:      ERROR,
		RequestID: em.RequestID,
		Message:   em.Message,
	})
}
