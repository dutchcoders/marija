package server

import (
	"context"
	_ "log"
)

func (c *connection) GetFields(ctx context.Context, r GetFieldsRequest) error {
	for _, server := range r.Datasources {
		datasource, ok := c.server.GetDatasource(server)
		if !ok {
			log.Errorf("Could not find datasource: %s", server)
			continue
		}

		fields, err := datasource.GetFields(ctx)
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
