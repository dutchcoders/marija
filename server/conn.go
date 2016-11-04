// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	_ "log"
	"net/http"
	"net/url"
	"path"
	"time"

	"gopkg.in/olivere/elastic.v3"

	"github.com/gorilla/websocket"
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
	CheckOrigin:     func(r *http.Request) bool { return true },
}

// connection is an middleman between the websocket connection and the hub.
type connection struct {
	// The websocket connection.
	ws *websocket.Conn

	// Buffered channel of outbound messages.
	send chan json.Marshaler

	b int
}

type ErrorMessage struct {
	Query   string `json:"query"`
	Color   string `json:"color"`
	Message string `json:"message"`
}

func (em *ErrorMessage) MarshalJSON() ([]byte, error) {
	type Alias ErrorMessage

	return json.Marshal(&struct {
		Error Alias `json:"error"`
	}{
		Error: (Alias)(*em),
	})
}

type ResultsMessage struct {
	Query   string      `json:"query"`
	Color   string      `json:"color"`
	Results interface{} `json:"results"`
}

func (em *ResultsMessage) MarshalJSON() ([]byte, error) {
	type Alias ResultsMessage

	return json.Marshal(&struct {
		Hits Alias `json:"hits"`
	}{
		Hits: (Alias)(*em),
	})
}

// readPump pumps messages from the websocket connection to the hub.
func (c *connection) readPump() {
	defer func() {
		h.unregister <- c
		c.ws.Close()
	}()

	c.ws.SetReadLimit(maxMessageSize)
	c.ws.SetReadDeadline(time.Now().Add(pongWait))
	c.ws.SetPongHandler(func(string) error { c.ws.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	for {
		_, message, err := c.ws.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Errorf("error: %v", err)
			}
			break
		}

		// check type of message
		// mapping
		// query

		log.Debug(string(message))

		v := map[string]interface{}{}
		if err := json.NewDecoder(bytes.NewBuffer(message)).Decode(&v); err != nil {
			c.send <- &ErrorMessage{
				Query:   v["query"].(string),
				Color:   v["color"].(string),
				Message: err.Error(),
			}

			return
		}

		func() {
			index := v["index"].(string)

			u, err := url.Parse(index)
			if err != nil {
				return
			}

			es, err := elastic.NewClient(elastic.SetURL(u.Host), elastic.SetSniff(false))
			if err != nil {
				fmt.Println(err.Error())

				c.send <- &ErrorMessage{
					Query:   v["query"].(string),
					Color:   v["color"].(string),
					Message: err.Error(),
				}
				return
			}

			hl := elastic.NewHighlight()
			hl = hl.Fields(elastic.NewHighlighterField("_all").NumOfFragments(0))
			hl = hl.PreTags("<em>").PostTags("</em>")

			results, err := es.Search().
				Index(path.Base(u.Path)). // search in index "twitter"
				Highlight(hl).
				Query(elastic.NewQueryStringQuery(v["query"].(string))). // specify the query
				From(c.b).Size(500).                                     // take documents 0-9
				Do()                                                     // execute
			if err != nil {
				fmt.Println(err.Error())

				c.send <- &ErrorMessage{
					Query:   v["query"].(string),
					Color:   v["color"].(string),
					Message: err.Error(),
				}
				return
			}

			c.send <- &ResultsMessage{
				Query:   v["query"].(string),
				Color:   v["color"].(string),
				Results: results,
			}

			c.b += 10
		}()

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
