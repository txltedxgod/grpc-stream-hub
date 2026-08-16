package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type Hub struct {
	mu      sync.RWMutex
	clients map[string]map[chan *MessageResponse]struct{}
}

func NewHub() *Hub {
	return &Hub{
		clients: make(map[string]map[chan *MessageResponse]struct{}),
	}
}

func (h *Hub) Subscribe(roomID string, ch chan *MessageResponse) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if _, ok := h.clients[roomID]; !ok {
		h.clients[roomID] = make(map[chan *MessageResponse]struct{})
	}
	h.clients[roomID][ch] = struct{}{}
}

func (h *Hub) Unsubscribe(roomID string, ch chan *MessageResponse) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if subs, ok := h.clients[roomID]; ok {
		delete(subs, ch)
		close(ch)
		if len(subs) == 0 {
			delete(h.clients, roomID)
		}
	}
}

func (h *Hub) Broadcast(msg *MessageResponse) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	if subs, ok := h.clients[msg.RoomId]; ok {
		for ch := range subs {
			select {
			case ch <- msg:
			default:
				// Skip slow consumer to prevent blocking hub
			}
		}
	}
}
