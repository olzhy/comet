package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/olzhy/comet/server"
	"golang.org/x/net/websocket"
)

var port = flag.Int("serverPort", 8080, "server port")

func main() {
	flag.Parse()

	// server
	wsServer := server.NewWsServer()
	httpServer := server.NewHttpServer(wsServer)
	h := server.NewHandler(wsServer, httpServer)
	go wsServer.Start()

	// handler
	r := mux.NewRouter()
	r.HandleFunc("/messages", h.MessageHandler).Methods(http.MethodPost)
	r.Handle("/comet", websocket.Handler(h.CometHandler))
	r.Headers("Content-Type", "application/json; charset=UTF-8")
	http.Handle("/", r)

	// test page
	http.Handle("/index/", http.FileServer(http.Dir("web/")))
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), nil))
}
