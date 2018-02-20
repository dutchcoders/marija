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

	go func() {
		items := []interface{}{}

		items = append(items, struct {
			ItemID string `json:"itemd_id"`
			Test   string `json:"test"`
			Test1  string `json:"test2"`
			Test2  string `json:"test3"`
			Test3  string `json:"test4"`
			Test4  string `json:"test5"`
		}{
			ItemID: "test",

			Test:  "test",
			Test1: "test1",
			Test2: "test2",
			Test3: "test3",
			Test4: "test4",
		})

		items = append(items, struct {
			ItemID string `json:"itemd_id"`
			Test   string `json:"test"`
			Test1  string `json:"test2"`
			Test2  string `json:"test3"`
			Test3  string `json:"test4"`
			Test4  string `json:"test5"`
		}{
			ItemID: "test",
			Test:   "test2",
			Test1:  "test1",
			Test2:  "test2",
			Test3:  "test3",
			Test4:  "test4",
		})

		c.Send(&ItemsResponse{
			RequestID: r.RequestID,
			Items:     items,
		})

		c.Send(&RequestCompleted{
			RequestID: r.RequestID,
		})

		/*
			response := datasource.Search(ctx, datasources.SearchOptions{
				Query: r.Query,
			})

			unique := Unique{}

			items := []interface{}{}

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
						c.Send(&ItemsResponse{
							RequestID: r.RequestID,
							Items:     items,
						})

						return
					}

					items = append(items, type Type struct {
						Test string `json:"test"`
					}{
						Test: "test",
					})

					if len(items) < 20 {
						continue
					}
				case <-time.After(time.Second * 5):
				}

				c.Send(&ItemsResponse{
					RequestID: r.RequestID,
					Items:     items,
				})

				items = []interface{}{}
			}
		*/
	}()

	return nil
}
