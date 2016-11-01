// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"path"
	"time"

	"gopkg.in/olivere/elastic.v3"
	"github.com/gorilla/websocket"
)

const (
	writeWait = 10 * time.Second
	pongWait = 60 * time.Second
	pingPeriod = 1 * time.Second
	maxMessageSize = 512
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool {
		return true
	},
}

type connection struct {
	ws   *websocket.Conn
	send chan interface{}
}

func (c *connection) reportError(event map[string]interface{}, error error) {
	log.Println(error.Error())

	c.send <- map[string]interface{}{
		"error": map[string]interface{}{
			"query":   event["query"].(string),
			"color":   event["color"].(string),
			"message": error.Error(),
		},
	}
}

func (c *connection) connectToEs(event map[string]interface{}) (esClient *elastic.Client, url *url.URL) {
	index := event["index"].(string)

	url, parseError := url.Parse(index)
	if parseError != nil {
		return
	}

	esClient, connectionError := elastic.NewClient(elastic.SetURL(url.Host), elastic.SetSniff(false))
	if connectionError != nil {
		c.reportError(event, connectionError)
		return
	}

	return esClient, url
}

func (c *connection) search(event map[string]interface{}) {

	es, url := c.connectToEs(event)

	hl := elastic.NewHighlight()
	hl = hl.Fields(elastic.NewHighlighterField("_all").NumOfFragments(0))
	hl = hl.PreTags("<em>").PostTags("</em>")

	results, err := es.Search().
		Index(path.Base(url.Path)).
		Highlight(hl).
		Query(elastic.NewQueryStringQuery(event["query"].(string))).
		From(0).Size(200).
		Do()

	if err != nil {
		c.reportError(event, err)
		return
	}

	c.send <- map[string]interface{}{
		"type": "search",
		"hits": map[string]interface{}{
			"query":   event["query"].(string),
			"color":   event["color"].(string),
			"results": results,
		},
	}
}

func (c *connection) discoverIndices(event map[string]interface{}) {
	es, _ := c.connectToEs(event)

	results, err := es.IndexStats().
		Metric("index").
		Do()

	if err != nil {
		c.reportError(event, err)
		return
	}

	c.send <- map[string]interface{}{
		"type": "index_discovery",
		"hits": map[string]interface{}{
			"server":   event["server"].(string),
			"indices": results.Indices,
		},
	}
}

func (c *connection) discoverFields(event map[string]interface{}) {
	es, url := c.connectToEs(event)

	results, err := es.FieldStats().
		Index(path.Base(url.Path)).
		Do()

	if err != nil {
		c.reportError(event, err)
		return
	}

	c.send <- map[string]interface{}{
		"type": "field_discovery",
		"hits": map[string]interface{}{
			"server":   event["server"].(string),
			"index":   event["index"].(string),
			"fields": results.Indices,
		},
	}
}

func (c *connection) readPump() {
	defer func() {
		h.unregister <- c
		c.ws.Close()
	}()

	c.ws.SetReadLimit(maxMessageSize)
	c.ws.SetReadDeadline(time.Now().Add(pongWait))
	c.ws.SetPongHandler(func(string) error {
		c.ws.SetReadDeadline(time.Now().Add(pongWait)); return nil
	})

	for {
		_, message, err := c.ws.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Printf("error: %v", err)
			}
			break
		}

		event := map[string]interface{}{}

		if err := json.NewDecoder(bytes.NewBuffer(message)).Decode(&event); err != nil {
			c.reportError(event, err)
			return
		}

		func() {
			event_type := event["event_type"].(float64)

			switch(event_type){
			case 1:
				c.search(event)
			case 2:
				c.discoverIndices(event)
			case 3:
				c.discoverFields(event)
			}
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
				log.Println(err.Error())
				return
			} else if err := c.write(websocket.TextMessage, buff.Bytes()); err != nil {
				log.Println(err.Error())
				return
			}
		case <-ticker.C:
			if err := c.write(websocket.PingMessage, []byte{}); err != nil {
				log.Println("%#v", err.Error())
				return
			}
		}
	}
}

func serveWs(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	c := &connection{send: make(chan interface{}, 256), ws: ws}

	h.register <- c

	log.Println("Connection upgraded.")

	go c.writePump()
	c.readPump()

	log.Println("Connection closed")
}
