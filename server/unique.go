package server

import (
	"bytes"
	_ "log"

	_ "github.com/dutchcoders/marija/server/datasources/blockchain"
	_ "github.com/dutchcoders/marija/server/datasources/es5"
	_ "github.com/dutchcoders/marija/server/datasources/twitter"
)

type Unique [][]byte

func (u *Unique) Contains(hash []byte) bool {
	for i := range *u {
		if bytes.Compare((*u)[i], hash) == 0 {
			return true
		}
	}

	return false
}

func (u *Unique) Add(hash []byte) {
	*u = append(*u, hash)
}
