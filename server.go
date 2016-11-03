package main

import (
	"fmt"
	_ "log"
	"net/http"

	"github.com/fatih/color"
)

// Debug - enables debugging.
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

type Server struct {
	address string
	debug   bool
}

func (server *Server) Run() {
	go h.run()

	http.Handle("/", http.FileServer(http.Dir(".")))
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
