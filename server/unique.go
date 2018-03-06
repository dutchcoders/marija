package server

import (
	_ "log"
	"sync"

	"github.com/dutchcoders/marija/server/datasources"
)

type Unique struct {
	sync.Map // map[string]*datasources.Node
}

func (u *Unique) Get(hash []byte) (*datasources.Node, bool) {
	item, ok := u.Load(string(hash))
	if !ok {
		return nil, false
	}

	return item.(*datasources.Node), ok
}

func (u *Unique) Add(hash []byte, value *datasources.Node) {
	u.Store(string(hash), value)
}
