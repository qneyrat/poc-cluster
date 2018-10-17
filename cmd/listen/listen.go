package main

import (
	"flag"
	"fmt"
	"net/url"
	"os"
	"os/signal"

	"github.com/gorilla/websocket"
)

func main() {
	n := flag.Int("n", 0, "number of clients")
	addr := flag.String("addr", "localhost:4000", "http service address")
	flag.Parse()
	u := url.URL{Scheme: "ws", Host: *addr, Path: "/"}
	for i := 0; i < *n; i++ {
		go listen(u)
	}
	listen(u)
}

func listen(u url.URL) {
	fmt.Printf("connecting to %s\n", u.String())
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		fmt.Println(err)
	}
	defer c.Close()

	done := make(chan struct{})
	go func() {
		defer c.Close()
		defer close(done)

		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println(string(message))
		}
	}()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	for {
		select {
		case <-interrupt:
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				fmt.Println(err)
				return
			}
			c.Close()
			return
		}
	}
}
