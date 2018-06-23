package server

import (
	"context"
	"encoding/json"
	"fmt"
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
	_ "github.com/dutchcoders/marija/server/datasources/openkvk"
	_ "github.com/dutchcoders/marija/server/datasources/splunk"
	_ "github.com/dutchcoders/marija/server/datasources/tronscan"
	_ "github.com/dutchcoders/marija/server/datasources/twitter"
	_ "github.com/dutchcoders/marija/server/datasources/voertuiggegevens"

	assetfs "github.com/elazarl/go-bindata-assetfs"

	logging "github.com/op/go-logging"

	web "github.com/dutchcoders/marija-web"
	"github.com/fatih/color"
)

var log = logging.MustGetLogger("marija/server")

type Server struct {
	*config

	Datasources map[string]datasources.Index
}

func New(options ...func(*Server)) *Server {
	server := &Server{
		config: &config{
			debug: false,
		},
	}

	for _, optionFn := range options {
		optionFn(server)
	}

	return server
}

func IsTerminal(f *os.File) bool {
	if isatty.IsTerminal(f.Fd()) {
		return true
	} else if isatty.IsCygwinTerminal(f.Fd()) {
		return true
	}

	return false
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

		ds, err := fn(
			datasources.WithConfig(s),
		)
		if err != nil {
			log.Error("Error parsing configuration of datasource: %s: %s", key, err.Error())
			continue
		}

		server.Datasources[key] = ds

		type Broadcasterer interface {
			Broadcast(context.Context, string) chan json.Marshaler
		}

		if bc, ok := ds.(Broadcasterer); ok {
			go func(key string) {
				for m := range bc.Broadcast(context.Background(), key) {
					h.Send(m)
				}
			}(key)
		}
	}

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	go func() {
		if err := http.ListenAndServe(server.address, nil); err != nil {
			log.Fatal("ListenAndServe: ", err)
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
