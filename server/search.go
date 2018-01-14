package server

import (
	"context"
	"fmt"
	"hash/fnv"
	_ "log"
	"runtime"
	"time"

	"github.com/dutchcoders/marija/server/datasources"

	_ "github.com/dutchcoders/marija/server/datasources/blockchain"
	_ "github.com/dutchcoders/marija/server/datasources/es5"
	_ "github.com/dutchcoders/marija/server/datasources/twitter"
	uuid "github.com/satori/go.uuid"
)

func (c *connection) Search(ctx context.Context, r SearchRequest) error {
	c.Send(&SearchResponse{
		RequestID: r.RequestID,
		Query:     r.Query,
	})

	indexes := []string{}

	if r.Datasource != "" {
		indexes = append(indexes, r.Datasource)
	}

	if len(r.Datasources) >= 1 {
		indexes = append(indexes, r.Datasources...)
	}

	for _, index := range indexes {
		datasource, ok := c.server.Datasources[index]
		if !ok {
			c.Send(&ErrorMessage{
				RequestID: r.RequestID,
				Message:   fmt.Sprintf("Could not find datasource: %s", index),
			})

			log.Errorf("Could not find datasource: %s", index)
			continue
		}

		go func() (err error) {
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

			unique := Unique{}

			items := []datasources.Item{}

			defer func() {
				if err == context.Canceled {
					c.Send(&SearchCanceled{
						RequestID: r.RequestID,
					})
				} else if err != nil {
					log.Error("Error: ", err.Error())

					c.Send(&ErrorMessage{
						RequestID: r.RequestID,
						Message:   err.Error(),
					})
				} else {
					c.Send(&SearchResponse{
						RequestID: r.RequestID,
						Query:     r.Query,
						Nodes:     items,
					})

					c.Send(&SearchCompleted{
						RequestID: r.RequestID,
					})
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

					// calculate hash of fields
					h := fnv.New128()
					for _, field := range values {
						switch s := field.(type) {
						case string:
							h.Write([]byte(s))
						default:
						}
					}

					hash := h.Sum(nil)

					i := &datasources.Item{
						ID:     uuid.NewV4().String(),
						Fields: values,
						Count:  1,
					}

					if v, ok := unique.Get(hash); ok {
						i = v

						i.Count++
					}

					unique.Add(hash, i)

					items = append(items, *i)

					if len(items) < 20 {
						continue
					}
				case <-time.After(time.Second * 5):
				}

				c.Send(&SearchResponse{
					RequestID: r.RequestID,
					Query:     r.Query,
					Nodes:     items,
				})

				items = []datasources.Item{}
			}
		}()
	}

	return nil
}
