package main

import (
	"flag"
	"fmt"
	"net/url"

	"github.com/gorilla/websocket"
)

func main() {
	addr := flag.String("addr", "localhost:4000", "http service address")
	message := flag.String("message", "hello", "message sended")
	flag.Parse()
	u := url.URL{Scheme: "ws", Host: *addr, Path: "/"}
	fmt.Printf("connecting to %s\n", u.String())
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		fmt.Println(err)
	}
	defer c.Close()
	c.WriteMessage(websocket.TextMessage, []byte(*message))
}
