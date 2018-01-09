package server

import (
	"context"
	_ "log"

	_ "github.com/dutchcoders/marija/server/datasources/blockchain"
	_ "github.com/dutchcoders/marija/server/datasources/es5"
	_ "github.com/dutchcoders/marija/server/datasources/twitter"
)

func (c *connection) Items(ctx context.Context, r ItemsRequest) error {
	c.Send(&ItemsResponse{
		RequestID: r.RequestID,
	})

	/*
		indexes := r.Datasources

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

			go func() {
				response := datasource.Search(ctx, datasources.SearchOptions{
					Query: r.Query,
				})

				unique := Unique{}

				items := []datasources.Item{}

				for {
					select {
					case <-ctx.Done():
						return
					case err, ok := <-response.Error():
						if !ok {
							return
						}

						log.Error("Error: ", err.Error())

						c.Send(&ErrorMessage{
							RequestID: r.RequestID,
							Message:   err.Error(),
						})

					case item, ok := <-response.Item():
						if !ok {
							c.Send(&SearchResponse{
								RequestID: r.RequestID,
								Query:     r.Query,
								Nodes:     items,
							})

							return
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

						if unique.Contains(hash) {
							continue
						}

						unique.Add(hash)

						items = append(items, datasources.Item{
							Fields: values,
						})

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
	*/

	return nil
}
