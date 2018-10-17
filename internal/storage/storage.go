package storage

import (
	"errors"
	"sync"
	"ws-cluster/internal/client"
)

type Storage struct {
	cc client.Clients
	sync.RWMutex
}

func (s *Storage) Add(c *client.Client) {
	s.Lock()
	defer s.Unlock()
	s.cc[c.ID] = c
}

func (s *Storage) GetAll() client.Clients {
	s.RLock()
	defer s.RUnlock()
	return s.cc
}

func (s *Storage) Get(id string) (*client.Client, error) {
	s.RLock()
	defer s.RUnlock()
	c, ok := s.cc[id]
	if !ok {
		return nil, errors.New("Client not found")
	}
	return c, nil
}

func (s *Storage) Delete(id string) {
	s.Lock()
	defer s.Unlock()
	delete(s.cc, id)
}

func NewStorage() *Storage {
	return &Storage{
		cc: make(client.Clients),
	}
}
