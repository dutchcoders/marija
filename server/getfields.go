package server

import (
	"context"
	_ "log"

	"github.com/dutchcoders/marija/server/datasources"

	_ "github.com/dutchcoders/marija/server/datasources/blockchain"
	_ "github.com/dutchcoders/marija/server/datasources/es5"
	_ "github.com/dutchcoders/marija/server/datasources/twitter"
)

func (c *connection) GetFields(ctx context.Context, r GetFieldsRequest) error {
	for _, server := range r.Datasources {
		var datasource datasources.Index
		if d, ok := c.server.Datasources[server]; !ok {
			log.Errorf("Could not find datasource: %s", server)
			continue
		} else {
			datasource = d
		}

		fields, err := datasource.Fields()
		if err != nil {
			return err
		}

		c.Send(&GetFieldsResponse{
			RequestID: r.RequestID,
			Server:    server,
			Fields:    fields,
		})
	}

	return nil
}
