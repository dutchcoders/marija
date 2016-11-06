package cmd

import (
	"github.com/dutchcoders/marija/server"
	"github.com/minio/cli"
	"github.com/op/go-logging"
)

var Version = "0.1"

var format = logging.MustStringFormatter(
	"%{color}%{time:15:04:05.000} %{shortfunc} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}",
)

var helpTemplate = `NAME:
{{.Name}} - {{.Usage}}

DESCRIPTION:
{{.Description}}

USAGE:
{{.Name}} {{if .Flags}}[flags] {{end}}command{{if .Flags}}{{end}} [arguments...]

COMMANDS:
{{range .Commands}}{{join .Names ", "}}{{ "\t" }}{{.Usage}}
{{end}}{{if .Flags}}
FLAGS:
{{range .Flags}}{{.}}
{{end}}{{end}}
VERSION:
` + Version +
	`{{ "\n"}}`

var log = logging.MustGetLogger("marija/cmd")

var globalFlags = []cli.Flag{
	cli.StringFlag{
		Name:  "p,port",
		Usage: "port",
		Value: "127.0.0.1:8080",
	},
	cli.StringFlag{
		Name:  "path",
		Usage: "path to static files",
		Value: "",
	},
	cli.StringFlag{
		Name:  "c,config",
		Usage: "config file",
		Value: "config.toml",
	},
	cli.BoolFlag{
		Name:  "help, h",
		Usage: "Show help.",
	},
}

type Cmd struct {
	*cli.App
}

func New() *Cmd {
	app := cli.NewApp()
	app.Name = "Marija"
	app.Author = ""
	app.Usage = "marija"
	app.Description = `Marija, graphing for Elasticsearch`
	app.Flags = globalFlags
	app.CustomAppHelpTemplate = helpTemplate
	app.Commands = []cli.Command{}

	app.Before = func(c *cli.Context) error {
		return nil
	}

	app.Action = func(c *cli.Context) {
		srvr := server.New(
			server.Address(c.String("port")),
			server.Path(c.String("path")),
		)

		srvr.Run()
	}

	return &Cmd{
		App: app,
	}
}
