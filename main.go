package main

import (
	"flag"
	"log"
	"net/http"
)

var addr = flag.String("addr", "0.0.0.0:8089", "http service address")

/*
func serveHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "Not found", 404)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	homeTempl.Execute(w, r.Host)
}
*/

type Packet struct {
}

func main() {
	flag.Parse()

	go h.run()

	http.Handle("/", http.FileServer(http.Dir(".")))
	http.HandleFunc("/ws", serveWs)

	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
