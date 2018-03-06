package server

import (
	_ "log"

	"github.com/BurntSushi/toml"
)

type config struct {
	path    string
	address string
	debug   bool

	ListenerString string `toml:"listen"`

	Username string `toml:"username"`
	Password string `toml:"password"`
	Service  string `toml:"service"`

	Datasources map[string]toml.Primitive `toml:"datasource"`

	Logging []struct {
		Output string `toml:"output"`
		Level  string `toml:"level"`
	} `toml:"logging"`
}
