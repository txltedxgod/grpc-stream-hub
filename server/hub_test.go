package main

import (
	"testing"
)

func TestHub_SubscribeBroadcastAndUnsubscribe(t *testing.T) {
	hub := NewHub()
	ch1 := make(chan *MessageResponse, 1)
	chOther := make(chan *MessageResponse, 1)
	slowCh := make(chan *MessageResponse) // unbuffered, should not block broadcast

	hub.Subscribe("general", ch1)
	hub.Subscribe("general", slowCh)
	hub.Subscribe("random", chOther)

	msg := &MessageResponse{
		MessageId: "m1",
		RoomId:    "general",
		Sender:    "alice",
		Content:   "hello",
	}
	hub.Broadcast(msg)

	select {
	case got := <-ch1:
		if got.MessageId != "m1" || got.Content != "hello" {
			t.Fatalf("unexpected message received: %+v", got)
		}
	default:
		t.Fatal("expected subscriber in room 'general' to receive message")
	}

	select {
	case extra := <-chOther:
		t.Fatalf("subscriber in room 'random' should not receive 'general' message: %+v", extra)
	default:
	}

	// Unsubscribe should be idempotent and safe against double-close
	hub.Unsubscribe("general", ch1)
	hub.Unsubscribe("general", ch1)
	hub.Unsubscribe("general", slowCh)
	hub.Unsubscribe("random", chOther)

	hub.mu.RLock()
	remaining := len(hub.clients)
	hub.mu.RUnlock()
	if remaining != 0 {
		t.Fatalf("expected all rooms cleaned up after unsubscribe, got %d", remaining)
	}
}
