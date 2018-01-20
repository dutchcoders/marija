package server

import (
	_ "log"

	"github.com/dutchcoders/marija/server/datasources"
	_ "github.com/dutchcoders/marija/server/datasources/blockchain"
	_ "github.com/dutchcoders/marija/server/datasources/es5"
	_ "github.com/dutchcoders/marija/server/datasources/twitter"
)

type Unique map[string]*datasources.Item

func (u Unique) Get(hash []byte) (*datasources.Item, bool) {
	item, ok := u[string(hash)]
	return item, ok
}

func (u Unique) Add(hash []byte, value *datasources.Item) {
	u[string(hash)] = value
}
