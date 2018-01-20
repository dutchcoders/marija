package server

import (
	"bytes"
	"context"
	"encoding/json"
	_ "log"
	"net/http"
	"time"

	_ "github.com/dutchcoders/marija/server/datasources/blockchain"
	_ "github.com/dutchcoders/marija/server/datasources/es5"
	_ "github.com/dutchcoders/marija/server/datasources/twitter"
	"github.com/gorilla/websocket"
	"github.com/op/go-logging"
)

var format = logging.MustStringFormatter(
	"%{color}%{time:15:04:05.000} %{shortfunc} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}",
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

	InitialStateReceive = "INITIAL_STATE_RECEIVE"

	ActionTypeCancel = "CANCEL_REQUEST"

	ActionTypeSearchRequest   = "SEARCH_REQUEST"
	ActionTypeSearchReceive   = "SEARCH_RECEIVE"
	ActionTypeSearchCanceled  = "SEARCH_CANCELED"
	ActionTypeSearchCompleted = "SEARCH_COMPLETED"

	ActionTypeItemsRequest = "ITEMS_REQUEST"
	ActionTypeItemsReceive = "ITEMS_RECEIVE"

	ActionTypeGetFieldsRequest = "FIELDS_REQUEST"
	ActionTypeGetFieldsReceive = "FIELDS_RECEIVE"
)

type connection struct {
	ws     *websocket.Conn
	send   chan json.Marshaler
	b      int
	server *Server
	closed bool
}

func (c *connection) Send(v json.Marshaler) {
	if c.closed {
		return
	}

	c.send <- v
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

	datasources := make([]Datasource, 0, len(c.server.Datasources))
	for k := range c.server.Datasources {
		datasources = append(datasources, Datasource{ID: k, Name: k})
	}

	c.Send(&InitialStateMessage{
		Datasources: datasources,
		Version:     Version,
		CommitID:    CommitID,
	})

	cancelFuncs := map[string]context.CancelFunc{}

	defer func() {
		for _, cancel := range cancelFuncs {
			cancel()
		}
	}()

	for {
		_, data, err := c.ws.ReadMessage()
		if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
			return
		} else if err != nil {
			log.Errorf("Error reading message: %+v", err.Error())
			return
		}

		r := Request{}
		if err := json.Unmarshal(data, &r); err != nil {
			log.Error("Error decoding message: ", err.Error())
			continue
		}

		if r.Type == ActionTypeCancel {
			cancel, ok := cancelFuncs[r.RequestID]
			if !ok {
				log.Error("Could not find cancel func for requestid: %s", r.RequestID)

			}

			cancel()
			continue
		}

		ctx, cancel := context.WithCancel(context.Background())
		cancelFuncs[r.RequestID] = cancel

		switch r.Type {
		case ActionTypeSearchRequest:
			r := SearchRequest{}
			if err := json.Unmarshal(data, &r); err != nil {
				log.Error("Error occured during search: %s", err.Error())

				c.Send(&ErrorMessage{
					RequestID: r.RequestID,
					Message:   err.Error(),
				})
			} else if err := c.Search(ctx, r); err != nil {
				log.Error("Error occured during search: %s", err.Error())

				c.Send(&ErrorMessage{
					RequestID: r.RequestID,
					Message:   err.Error(),
				})
			}
		case ActionTypeItemsRequest:
			r := ItemsRequest{}
			if err := json.Unmarshal(data, &r); err != nil {
				log.Error("Error occured retrieving items: %s", err.Error())
				c.Send(&ErrorMessage{
					RequestID: r.RequestID,
					Message:   err.Error(),
				})
			} else if err := c.Items(ctx, r); err != nil {
				log.Error("Error occured retrieving items: %s", err.Error())

				c.Send(&ErrorMessage{
					RequestID: r.RequestID,
					Message:   err.Error(),
				})
			}
		case ActionTypeGetFieldsRequest:
			r := GetFieldsRequest{}
			if err := json.Unmarshal(data, &r); err != nil {
				log.Error("Error occured during search: %s", err.Error())
				c.Send(&ErrorMessage{
					RequestID: r.RequestID,
					Message:   err.Error(),
				})
			} else if err := c.GetFields(ctx, r); err != nil {
				log.Error("Error occured during field discovery: %s", err.Error())
				c.Send(&ErrorMessage{
					RequestID: r.RequestID,
					Message:   err.Error(),
				})
			}
		default:
			log.Error("Unknown request: %s", r.Type)
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
				c.closed = true
				c.write(websocket.CloseMessage, []byte{})
				return
			}

			buff := new(bytes.Buffer)
			if err := json.NewEncoder(buff).Encode(message); err != nil {
				log.Error("%s", err.Error())
				return
			} else if err := c.write(websocket.TextMessage, buff.Bytes()); err != nil {
				log.Error("%s", err.Error())
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
func (s *Server) serveWs(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error(err.Error())
		return
	}

	c := &connection{
		send:   make(chan json.Marshaler, 256),
		ws:     ws,
		server: s,
	}

	h.register <- c

	log.Info("Connection upgraded.")
	defer log.Info("Connection closed")

	go c.writePump()
	c.readPump()
}
