package server

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"hash/fnv"
	_ "log"
	"runtime"
	"time"

	"github.com/dutchcoders/marija/server/datasources"
	"github.com/dutchcoders/marija/server/messages"
	"github.com/dutchcoders/marija/server/unique"
)

func (c *connection) Search(ctx context.Context, r messages.SearchRequest) error {
	if len(r.Datasources) == 0 {
		return errors.New("No datasource set")
	}

	for _, index := range r.Datasources {
		log.Debug("Search query=%s, request=%s, index=%s", r.Query, r.RequestID, index)

		c.Send(&messages.SearchResponse{
			RequestID:  r.RequestID,
			Query:      r.Query,
			Datasource: index,
		})

		datasource, ok := c.server.GetDatasource(index)
		if !ok {
			c.Send(&messages.ErrorMessage{
				RequestID: r.RequestID,
				Message:   fmt.Sprintf("Could not find datasource: %s", index),
			})

			log.Errorf("Could not find datasource: %s", index)
			continue
		}

		go func(index string) (err error) {
			defer func() {
				if err := recover(); err != nil {
					trace := make([]byte, 1024)
					count := runtime.Stack(trace, true)
					log.Errorf("Error: %s", err)
					log.Debugf("Stack of %d bytes: %s\n", count, string(trace))
				}
			}()

			response := datasource.Search(ctx, datasources.SearchOptions{
				Query: r.Query,
			})

			unique := unique.New()

			graphs := []datasources.Graph{}

			defer func() {
				if err == context.Canceled {
					log.Debug("Search canceled query=%s, requestid=%s, index=%s", r.Query, r.RequestID, index)

					c.Send(&messages.RequestCanceled{
						RequestID: r.RequestID,
					})
				} else if err != nil {
					log.Error("Search error query=%s, requestid=%s, index=%s, error=%s", r.Query, r.RequestID, index, err.Error())

					c.Send(&messages.ErrorMessage{
						RequestID: r.RequestID,
						Message:   err.Error(),
					})
				} else {
					c.Send(&messages.SearchResponse{
						RequestID:  r.RequestID,
						Query:      r.Query,
						Graphs:     graphs,
						Datasource: index,
					})

					c.Send(&messages.RequestCompleted{
						RequestID: r.RequestID,
					})

					log.Debug("Search completed query=%s, requestid=%s, index=%s", r.Query, r.RequestID, index)
				}
			}()

			for {
				select {
				case <-ctx.Done():
					err = ctx.Err()
					return
				case err, ok := <-response.Error():
					if !ok {
						return nil
					}

					return err
				case item, ok := <-response.Item():
					if !ok {
						return nil
					}

					// filter fields
					values := map[string]interface{}{}

					for _, field := range r.Fields {
						v, ok := item.Fields[field]
						if !ok {
							continue
						}

						values[field] = v
					}

					count := 0
					for _ = range values {
						count++
					}

					if count == 0 {
						continue
					}

					// calculate hash of fields
					h := fnv.New128()
					for _, field := range values {
						switch s := field.(type) {
						case []string:
							for _, v := range s {
								h.Write([]byte(v))
							}
						case string:
							h.Write([]byte(s))
						default:
						}
					}

					hash := h.Sum(nil)
					hashHex := hex.EncodeToString(hash)

					i := &datasources.Graph{
						ID:         hashHex,
						Fields:     values,
						Datasource: index,
						Count:      1,
					}

					if v, ok := unique.Get(hash); ok {
						i = v

						i.Count++
					}

					unique.Add(hash, i)

					items, _ := c.items.LoadOrStore(i.ID, []datasources.Item{})
					items = append(items, item)

					c.items.Store(i.ID, items)

					graphs = append(graphs, *i)

					if len(graphs) < 20 {
						continue
					}
				case <-time.After(time.Second * 5):
				}

				if len(graphs) == 0 {
					continue
				}

				c.Send(&messages.SearchResponse{
					RequestID:  r.RequestID,
					Query:      r.Query,
					Graphs:     graphs,
					Datasource: index,
				})

				graphs = []datasources.Graph{}
			}
		}(index)
	}

	return nil
}
