package cmd

import (
	"fmt"

	"github.com/dutchcoders/marija/server"
	"github.com/fatih/color"
	"github.com/minio/cli"
	"github.com/op/go-logging"
)

var Version = "0.1"
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
}

type Cmd struct {
	*cli.App
}

func VersionAction(c *cli.Context) {
	fmt.Println(color.YellowString(fmt.Sprintf("Marija: Exploration and visualisation of Elasticsearch data.")))
}

func New() *Cmd {
	cli.VersionPrinter = VersionPrinter

	app := cli.NewApp()
	app.Name = "Marija"
	app.Author = ""
	app.Usage = "marija"
	app.Description = `Exploration and visualisation of Elasticsearch data`
	app.Flags = globalFlags
	app.CustomAppHelpTemplate = helpTemplate
	app.Commands = []cli.Command{
		{
			Name:   "version",
			Action: VersionAction,
		},
	}

	app.Before = func(c *cli.Context) error {
		return nil
	}

	app.Action = func(c *cli.Context) {
		options := []func(*server.Server){}

		if v := c.String("port"); v != "" {
			options = append(options, server.Address(v))

		}

		if v := c.String("path"); v != "" {
			options = append(options, server.Path(v))
		}

		if v := c.String("config"); v != "" {
			options = append(options, server.Config(v))
		}

		srvr := server.New(
			options...,
		)

		srvr.Run()
	}

	return &Cmd{
		App: app,
	}
}
