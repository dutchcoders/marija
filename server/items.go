package server

import (
	"context"
	_ "log"
	"runtime"

	"github.com/dutchcoders/marija/server/datasources"
)

func (c *connection) Items(ctx context.Context, r ItemsRequest) error {
	c.Send(&ItemsResponse{
		RequestID: r.RequestID,
	})

	go func() (err error) {
		defer func() {
			if err := recover(); err != nil {
				trace := make([]byte, 1024)
				count := runtime.Stack(trace, true)
				log.Errorf("Error: %s", err)
				log.Debugf("Stack of %d bytes: %s\n", count, string(trace))
			}
		}()

		items := []datasources.Item{}

		defer func() {
			if err == context.Canceled {
				c.Send(&RequestCanceled{
					RequestID: r.RequestID,
				})
			} else if err != nil {
				log.Error("Error: ", err.Error())

				c.Send(&ErrorMessage{
					RequestID: r.RequestID,
					Message:   err.Error(),
				})
			} else {
				c.Send(&ItemsResponse{
					RequestID: r.RequestID,
					Items:     items,
				})

				c.Send(&RequestCompleted{
					RequestID: r.RequestID,
				})
			}
		}()

		for _, itemid := range r.Items {
			items := c.items[itemid]

			if len(items) == 0 {
				continue
			}

			c.Send(&ItemsResponse{
				Items: items,
			})
		}

		return nil
	}()

	return nil
}
