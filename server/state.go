package main

import (
	"fmt"
	"sync"
	"time"
)

type TicketState struct {
	Name       string
	State      string
	TrackCount int
	CreatedAt  time.Time
}

type LabState struct {
	mu      sync.Mutex
	counter int64
	tickets map[string]*TicketState
}

func NewLabState() *LabState {
	return &LabState{
		tickets: make(map[string]*TicketState),
	}
}

func (s *LabState) nextID(prefix string) string {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.counter++

	return fmt.Sprintf(
		"%s-%d-%d",
		prefix,
		time.Now().UnixNano(),
		s.counter,
	)
}