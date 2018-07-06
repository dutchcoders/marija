package server

import (
	"context"
	_ "log"
	"runtime"

	"github.com/dutchcoders/marija/server/datasources"
	"github.com/dutchcoders/marija/server/messages"
)

func (c *connection) Items(ctx context.Context, r messages.ItemsRequest) error {
	c.Send(&messages.ItemsResponse{
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
				c.Send(&messages.RequestCanceled{
					RequestID: r.RequestID,
				})
			} else if err != nil {
				log.Error("Error: ", err.Error())

				c.Send(&messages.ErrorMessage{
					RequestID: r.RequestID,
					Message:   err.Error(),
				})
			} else {
				c.Send(&messages.ItemsResponse{
					RequestID: r.RequestID,
					Items:     items,
				})

				c.Send(&messages.RequestCompleted{
					RequestID: r.RequestID,
				})
			}
		}()

		for _, itemid := range r.Items {
			items, ok := c.items.Load(itemid)
			if !ok {
				continue
			}

			if len(items) == 0 {
				continue
			}

			c.Send(&messages.ItemsResponse{
				RequestID: r.RequestID,
				ItemID:    itemid,
				Items:     items,
			})
		}

		return nil
	}()

	return nil
}
