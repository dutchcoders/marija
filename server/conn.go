// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package server

import (
	"bytes"
	"encoding/json"
	_ "log"
	"net/http"
	"net/url"
	"path"
	"time"

	"github.com/gorilla/websocket"
	"gopkg.in/olivere/elastic.v3"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = 1 * time.Second

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

const (
	ERROR = "ERROR"

	ActionTypeItemsRequest = "ITEMS_REQUEST"
	ActionTypeItemsReceive = "ITEMS_RECEIVE"

	ActionTypeIndicesRequest = "INDICES_REQUEST"
	ActionTypeIndicesReceive = "INDICES_RECEIVE"

	ActionTypeFieldsRequest = "FIELDS_REQUEST"
	ActionTypeFieldsReceive = "FIELDS_RECEIVE"
)

type connection struct {
	ws   *websocket.Conn
	send chan json.Marshaler
	b    int
}

type ErrorMessage struct {
	Query   string `json:"query"`
	Color   string `json:"color"`
	Message string `json:"message"`
}

func (em *ErrorMessage) MarshalJSON() ([]byte, error) {
	type Alias ErrorMessage

	return json.Marshal(&struct {
		Type  string `json:"type"`
		Error Alias  `json:"error"`
	}{
		Type:  ERROR,
		Error: (Alias)(*em),
	})
}

type ResultsMessage struct {
	Server  string `json:"server"`
	Query   string `json:"query"`
	Color   string `json:"color"`
	Results []Item `json:"results"`
}

func (em *ResultsMessage) MarshalJSON() ([]byte, error) {
	type Alias ResultsMessage

	return json.Marshal(&struct {
		Type string `json:"type"`
		Hits Alias  `json:"hits"`
	}{
		Type: ActionTypeItemsReceive,
		Hits: (Alias)(*em),
	})
}

type IndicesMessage struct {
	Host    string   `json:"server"`
	Indices []string `json:"indices"`
}

func (em *IndicesMessage) MarshalJSON() ([]byte, error) {
	type Alias IndicesMessage

	return json.Marshal(&struct {
		Type string `json:"type"`
		Hits Alias  `json:"hits"`
	}{
		Type: ActionTypeIndicesReceive,
		Hits: (Alias)(*em),
	})
}

type FieldsMessage struct {
	Server string      `json:"server"`
	Index  string      `json:"index"`
	Fields interface{} `json:"fields"`
}

func (em *FieldsMessage) MarshalJSON() ([]byte, error) {
	type Alias FieldsMessage

	return json.Marshal(&struct {
		Type string `json:"type"`
		Hits Alias  `json:"hits"`
	}{
		Type: ActionTypeFieldsReceive,
		Hits: (Alias)(*em),
	})
}

func (c *connection) ConnectToEs(u *url.URL) (*elastic.Client, error) {
	return elastic.NewClient(elastic.SetURL(u.Host), elastic.SetSniff(false))
}

func (c *connection) Search(event map[string]interface{}) error {
	indexes := event["host"].([]interface{})

	for _, index := range indexes {
		u, err := url.Parse(index.(string))
		if err != nil {
			return err
		}

		i, err := NewElasticsearchIndex(u)
		if err != nil {
			return err
		}

		items, err := i.Search(SearchOptions{
			From:  event["from"].(int),
			Size:  event["size"].(int),
			Query: event["query"].(string),
		})
		if err != nil {
			return err
		}

		c.send <- &ResultsMessage{
			Query:   event["query"].(string),
			Color:   event["color"].(string),
			Server:  index.(string),
			Results: items,
		}
	}

	return nil
}

// return complete url instead of only name of index?
func (c *connection) DiscoverIndices(event map[string]interface{}) error {
	hosts := event["host"].([]interface{})
	for _, host := range hosts {
		u, err := url.Parse(host.(string))
		if err != nil {
			return err
		}

		i, err := NewElasticsearchIndex(u)
		if err != nil {
			return err
		}

		indices, err := i.Indices()
		if err != nil {
			return err
		}

		c.send <- &IndicesMessage{
			Host:    host.(string),
			Indices: indices,
		}
	}

	return nil
}

type Field struct {
	Path string `json:"path"`
	Type string `json:"type"`
}

func (c *connection) DiscoverFields(event map[string]interface{}) error {
	servers := event["host"].([]interface{})
	for _, index := range servers {
		u, err := url.Parse(index.(string))
		if err != nil {
			return err
		}

		i, err := NewElasticsearchIndex(u)
		if err != nil {
			return err
		}

		fields, err := i.Fields(path.Base(u.Path))
		if err != nil {
			return err
		}

		c.send <- &FieldsMessage{
			Server: index.(string),
			Fields: fields,
		}
	}

	return nil
}

// readPump pumps messages from the websocket connection to the hub.
func (c *connection) readPump() {
	defer func() {
		h.unregister <- c
		c.ws.Close()
	}()

	c.ws.SetReadLimit(maxMessageSize)
	c.ws.SetReadDeadline(time.Now().Add(pongWait))
	c.ws.SetPongHandler(func(string) error {
		c.ws.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, message, err := c.ws.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Errorf("error: %v", err)
			}
			break
		}

		log.Debug("Message", string(message))

		v := map[string]interface{}{}
		if err := json.NewDecoder(bytes.NewBuffer(message)).Decode(&v); err != nil {
			log.Error("Error decoding message: ", err.Error())

			c.send <- &ErrorMessage{
				Query:   v["query"].(string),
				Color:   v["color"].(string),
				Message: err.Error(),
			}
			continue
		}

		t := v["type"].(string)
		switch t {
		case ActionTypeItemsRequest:
			if err := c.Search(v); err != nil {
				log.Error("Error occured during search: %s", err.Error())
				c.send <- &ErrorMessage{
					Query:   "",
					Color:   "",
					Message: err.Error(),
				}
			}
		case ActionTypeIndicesRequest:
			if err := c.DiscoverIndices(v); err != nil {
				log.Error("Error occured during index discovery: %s", err.Error())
				c.send <- &ErrorMessage{
					Query:   "",
					Color:   "",
					Message: err.Error(),
				}
			}
		case ActionTypeFieldsRequest:
			if err := c.DiscoverFields(v); err != nil {
				log.Error("Error occured during field discovery: %s", err.Error())
				c.send <- &ErrorMessage{
					Query:   "",
					Color:   "",
					Message: err.Error(),
				}
			}
		}
	}
}

// write writes a message with the given message type and payload.
func (c *connection) write(mt int, payload []byte) error {
	c.ws.SetWriteDeadline(time.Now().Add(writeWait))
	return c.ws.WriteMessage(mt, payload)
}

// writePump pumps messages from the hub to the websocket connection.
func (c *connection) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.ws.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				c.write(websocket.CloseMessage, []byte{})
				return
			}

			buff := new(bytes.Buffer)
			if err := json.NewEncoder(buff).Encode(message); err != nil {
				log.Error(err.Error())
				return
			} else if err := c.write(websocket.TextMessage, buff.Bytes()); err != nil {
				log.Error(err.Error())
				return
			}
		case <-ticker.C:
			if err := c.write(websocket.PingMessage, []byte{}); err != nil {
				log.Error("%#v", err.Error())
				return
			}
		}
	}
}

// serveWs handles websocket requests from the peer.
func serveWs(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error(err)
		return
	}

	c := &connection{send: make(chan json.Marshaler, 256), ws: ws}

	h.register <- c

	log.Info("Connection upgraded.")

	go c.writePump()
	c.readPump()

	log.Info("Connection closed")
}
