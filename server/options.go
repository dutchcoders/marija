package server

import _ "log"

func Debug() func(*Server) {
	return func(server *Server) {
		server.debug = true
	}
}

func Address(addr string) func(*Server) {
	return func(server *Server) {
		server.address = addr
	}
}
