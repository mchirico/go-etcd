package socket

import (
	"flag"
	"log"
	"net/http"
)

var addr = flag.String("addr", ":8000", "http service address")

func ServerHome(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	if r.URL.Path != "/t" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Write(HomeHTML())

}

func Example() {

	hub := NewHub()
	go hub.Run()

	h := HANDLE{}
	// Broadcast channel
	h.broadcast = hub.broadcast

	http.HandleFunc("/base", h.BaseRoot)
	http.HandleFunc("/t", ServerHome)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		ServWs(hub, w, r)
	})
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
