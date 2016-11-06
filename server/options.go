package server

import _ "log"

func Debug() func(*Server) {
	return func(server *Server) {
		server.debug = true
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
