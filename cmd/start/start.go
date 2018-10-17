package main

import (
	"flag"
	"net/http"
	"ws-cluster/server"

	"github.com/gorilla/websocket"
)

func main() {
	port := flag.String("port", "4000", "port service")
	flag.Parse()

	u := websocket.Upgrader{
		EnableCompression: true,
		ReadBufferSize:    1024,
		WriteBufferSize:   1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	s := server.NewServer(u)
	s.Start(*port)
}
