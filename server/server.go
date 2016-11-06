package server

import (
	"fmt"
	_ "log"
	"net/http"
	"os"

	web "github.com/dutchcoders/marija-web"
	"github.com/elazarl/go-bindata-assetfs"
	"github.com/fatih/color"
	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("marija/server")

type Server struct {
	path    string
	address string
	debug   bool
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
		staticHandler = http.FileServer(http.Dir(server.path))
	}

	http.Handle("/", staticHandler)

	http.HandleFunc("/ws", serveWs)

	fmt.Println(color.YellowString(fmt.Sprintf("Marija server started, listening on address %s.", server.address)))

	defer func() {
		fmt.Println(color.YellowString(fmt.Sprintf("Marija server stopped.")))
	}()

	if err := http.ListenAndServe(server.address, nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func New(options ...func(*Server)) *Server {
	server := &Server{
		debug: false,
	}

	for _, optionFn := range options {
		optionFn(server)
	}

	return server
}
