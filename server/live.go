package server

import (
	"encoding/json"
	"io"
	"net/http"
	"os"

	"github.com/fatih/color"
)

func flattenFields(root string, m map[string]interface{}) map[string]interface{} {
	fields := map[string]interface{}{}

	for k, v := range m {
		key := k
		if root != "" {
			key = root + "." + key
		}

		switch s2 := v.(type) {
		case map[string]interface{}:
			for k2, v2 := range flattenFields(key, s2) {
				fields[k2] = v2
			}
		default:
			fields[key] = v
		}
	}

	return fields
}

// SubmitHandler will receive requests from other datasources
func (server *Server) SubmitHandler(w http.ResponseWriter, r *http.Request) {
	var doc map[string]interface{}
	err := json.NewDecoder(io.TeeReader(r.Body, os.Stdout)).Decode(&doc)
	if err != nil {
		log.Error(color.RedString("Submit could not parse body: %s", err.Error()))
		return
	}

	fields := flattenFields("", doc)

	if len(fields) == 0 {
		return
	}

	key := r.URL.RawQuery
	if key == "" {
		key = "wodan"
	}

	ds, ok := server.GetDatasource(key)
	if !ok {
		log.Error("Could not find datasource: %s", key)
		return
	}

	type Receiverer interface {
		Receive(m map[string]interface{})
	}

	s, ok := ds.(Receiverer)
	if !ok {
		log.Error("%s does not support receiverer", key)
		return
	}

	s.Receive(fields)
}
