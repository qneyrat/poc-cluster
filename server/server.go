package server

import (
	"net/http"
	"ws-cluster/internal/client"
	"ws-cluster/internal/storage"

	"github.com/garyburd/redigo/redis"
	"github.com/gorilla/websocket"
	uuid "github.com/satori/go.uuid"
)

type Server struct {
	Channel chan []byte
	Storage *storage.Storage

	conn     redis.Conn
	upgrader websocket.Upgrader
}

func NewServer(u websocket.Upgrader) *Server {
	return &Server{
		Channel:  make(chan []byte),
		Storage:  storage.NewStorage(),
		upgrader: u,
	}
}

func (s *Server) Start(port string) error {
	conn, err := redis.Dial("tcp", ":6379")
	if err != nil {
		return err
	}
	defer conn.Close()
	s.conn = conn

	go sync(s.Channel)

	go s.handleMessages()
	http.HandleFunc("/", s.handleConnections)
	return http.ListenAndServe(":"+port, nil)
}

func (s *Server) handleConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer ws.Close()

	id := uuid.NewV4().String()
	s.Storage.Add(&client.Client{
		ID:   id,
		Conn: ws,
	})
	defer s.Storage.Delete(id)

	for {
		_, p, err := ws.ReadMessage()
		if err != nil {
			break
		}
		s.conn.Do("PUBLISH", "*", p)
	}
}

func (s *Server) handleMessages() {
	for {
		m := <-s.Channel
		for _, c := range s.Storage.GetAll() {
			ws := c.Conn
			err := ws.WriteMessage(websocket.TextMessage, m)
			if err != nil {
				ws.Close()
				s.Storage.Delete(c.ID)
			}
		}
	}
}

func sync(c chan []byte) {
	conn, err := redis.Dial("tcp", ":6379")
	if err != nil {
		return
	}
	defer conn.Close()
	ps := &redis.PubSubConn{Conn: conn}
	defer ps.Close()
	ps.Subscribe("*")
	for {
		switch v := ps.Receive().(type) {
		case redis.Message:
			c <- v.Data
		case error:
			return
		}
	}
}
