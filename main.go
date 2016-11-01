package main

import (
	"flag"
	"log"
	"net/http"
)

var addr = flag.String("addr", "0.0.0.0:8089", "http service address")

type Packet struct{}

func main() {
	flag.Parse()

	go h.run()

	http.Handle("/", http.FileServer(http.Dir(".")))
	http.HandleFunc("/ws", serveWs)

	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
