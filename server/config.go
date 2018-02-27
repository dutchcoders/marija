package server

import (
	"io"
	_ "log"
	"os"

	"github.com/BurntSushi/toml"
	_ "github.com/dutchcoders/marija/server/datasources/blockchain"
	_ "github.com/dutchcoders/marija/server/datasources/es5"
	_ "github.com/dutchcoders/marija/server/datasources/twitter"
	"github.com/op/go-logging"
)

func Debug2() func(*Server) {
	return func(server *Server) {
		server.debug = true
	}
}

func Config(val string) func(*Server) {
	return func(server *Server) {
		if _, err := toml.DecodeFile(val, &server); err != nil {
			panic(err)
		}

		logBackends := []logging.Backend{}
		for _, log := range server.Logging {
			var err error

			var output io.Writer = os.Stdout

			switch log.Output {
			case "stdout":
				output = os.Stdout
			case "stderr":
				output = os.Stderr
			default:
				output, err = os.OpenFile(os.ExpandEnv(log.Output), os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0660)
			}

			if err != nil {
				panic(err)
			}

			backend := logging.NewLogBackend(output, "", 0)
			backendFormatter := logging.NewBackendFormatter(backend, format)
			backendLeveled := logging.AddModuleLevel(backendFormatter)

			level, err := logging.LogLevel(log.Level)
			if err != nil {
				panic(err)
			}

			backendLeveled.SetLevel(level, "")

			logBackends = append(logBackends, backendLeveled)
		}

		logging.SetBackend(logBackends...)

	}
}

func Path(val string) func(*Server) {
	return func(server *Server) {
		server.path = val
	}
}

func Address(addr string) func(*Server) {
	return func(server *Server) {
		server.address = addr
	}
}
