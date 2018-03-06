package server

import (
	"context"
	_ "log"

	"github.com/dutchcoders/marija/server/messages"
)

func (c *connection) GetFields(ctx context.Context, r messages.GetFieldsRequest) error {
	for _, datasource := range r.Datasources {
		go func(datasource string) {
			log.Debug("GetFields request=%s, index=%s", r.RequestID, datasource)
			defer log.Debug("GetFields completed request=%s, index=%s", r.RequestID, datasource)

			ds, ok := c.server.GetDatasource(datasource)
			if !ok {
				log.Errorf("Could not find datasource: %s", datasource)
				return
			}

			fields, err := ds.GetFields(ctx)
			if err != nil {
				c.Send(&messages.ErrorMessage{
					RequestID: r.RequestID,
					Message:   err.Error(),
				})
				return
			}

			c.Send(&messages.GetFieldsResponse{
				RequestID:  r.RequestID,
				Datasource: datasource,
				Fields:     fields,
			})
		}(datasource)
	}

	return nil
}
