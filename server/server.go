package server

import (
	"fmt"
	"hash/fnv"
	_ "log"
	"net/http"
	"os"
	"os/signal"

	"github.com/BurntSushi/toml"
	"github.com/dutchcoders/marija/server/datasources"
	isatty "github.com/mattn/go-isatty"

	_ "github.com/dutchcoders/marija/server/datasources/blockchain"
	_ "github.com/dutchcoders/marija/server/datasources/es5"
	_ "github.com/dutchcoders/marija/server/datasources/live"
	_ "github.com/dutchcoders/marija/server/datasources/splunk"
	_ "github.com/dutchcoders/marija/server/datasources/twitter"

	assetfs "github.com/elazarl/go-bindata-assetfs"

	logging "github.com/op/go-logging"

	"encoding/hex"
	"encoding/json"

	web "github.com/dutchcoders/marija-web"
	"github.com/fatih/color"
)

var log = logging.MustGetLogger("marija/server")

type Server struct {
	*config

	Datasources map[string]datasources.Index

	liveCh chan map[string]interface{}
}

func New(options ...func(*Server)) *Server {
	server := &Server{
		config: &config{
			debug: false,
		},

		liveCh: make(chan map[string]interface{}, 100),
	}

	for _, optionFn := range options {
		optionFn(server)
	}

	return server
}

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

func IsTerminal(f *os.File) bool {
	if isatty.IsTerminal(f.Fd()) {
		return true
	} else if isatty.IsCygwinTerminal(f.Fd()) {
		return true
	}

	return false
}

func (server *Server) SubmitHandler(w http.ResponseWriter, r *http.Request) {
	var fields map[string]interface{}

	err := json.NewDecoder(r.Body).Decode(&fields)
	if err != nil {
		log.Error(color.RedString(err.Error()))
		return
	}

	fields = flattenFields("", fields)

	if len(fields) == 0 {
		return
	}

	server.liveCh <- fields
}

func (server *Server) Run() {
	go h.run()

	staticHandler := http.FileServer(
		&assetfs.AssetFS{
			Asset:    web.Asset,
			AssetDir: web.AssetDir,
			AssetInfo: func(path string) (os.FileInfo, error) {
				return os.Stat(path)
			},
			Prefix: web.Prefix,
		})

	if server.path != "" {
		log.Debug("Using static file path: ", server.path)

		// check local css first
		staticHandler = http.FileServer(http.Dir(server.path))
	}

	http.Handle("/", staticHandler)

	http.HandleFunc("/submit", server.SubmitHandler)
	http.HandleFunc("/ws", server.serveWs)

	if IsTerminal(os.Stdout) {
		fmt.Println(color.YellowString(`
 __  __            _  _
|  \/  | __ _ _ __(_)(_) __ _
| |\/| |/ _' | '__| || |/ _' |
| |  | | (_| | |  | || | (_| |
|_|  |_|\__,_|_|  |_|/ |\__,_|
                   |__/
`))
	}

	fmt.Println(color.YellowString("Marija server started %s (%s)", Version, ShortCommitID))
	fmt.Println(color.YellowString("Listening on address %s.", server.address))

	defer func() {
		fmt.Println(color.YellowString("Marija server stopped"))
	}()

	server.Datasources = map[string]datasources.Index{}

	for key, s := range server.config.Datasources {
		x := struct {
			Type string `toml:"type"`
		}{}

		err := toml.PrimitiveDecode(s, &x)
		if err != nil {
			log.Error("Error parsing configuration of datasource: %s: %s", key, err.Error())
			continue
		}

		fn, err := datasources.Get(x.Type)
		if err != nil {
			log.Error("Error parsing configuration of datasource: %s: %s", key, err.Error())
			continue
		}

		ds, err := fn(datasources.WithConfig(s))
		if err != nil {
			log.Error("Error parsing configuration of datasource: %s: %s", key, err.Error())
			continue
		}

		server.Datasources[key] = ds
	}

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	go func() {
		if err := http.ListenAndServe(server.address, nil); err != nil {
			log.Fatal("ListenAndServe: ", err)
		}
	}()

	go func() {
		for fields := range server.liveCh {
			// calculate hash of fields
			hash := fnv.New128()
			for _, field := range fields {
				switch s := field.(type) {
				case []string:
					for _, v := range s {
						hash.Write([]byte(v))
					}
				case string:
					hash.Write([]byte(s))
				default:
				}
			}

			hashHex := hex.EncodeToString(hash.Sum(nil))

			graphs := []datasources.Graph{
				datasources.Graph{
					ID:     hashHex,
					Fields: fields,
					Count:  1,
				},
			}

			for c := range h.connections {
				//
				c.Send(&LiveResponse{
					Datasource: "wodan",
					Graphs:     graphs,
				})
			}
		}
	}()

	for {
		select {
		case <-signals:
			fmt.Println(color.YellowString("Marija server stopping..."))
			return
		}
	}

}

func (s *Server) GetDatasource(key string) (datasources.Index, bool) {
	if _, ok := s.Datasources[key]; !ok {
		return nil, false
	}

	return s.Datasources[key], true
}
