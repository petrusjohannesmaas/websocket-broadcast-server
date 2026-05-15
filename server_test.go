package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

func TestServer_Broadcast(t *testing.T) {
	// Initialize server with simplified signature
	s := NewServer()

	// Setup httptest server using the handler directly
	ts := httptest.NewServer(http.HandlerFunc(s.handleWebSocket))
	defer ts.Close()

	// Convert URL to websocket scheme
	url := "ws" + strings.TrimPrefix(ts.URL, "http")

	// Table-driven tests for broadcast logic
	tests := []struct {
		name    string
		message string
	}{
		{
			name:    "Send simple greeting",
			message: "Hello, World!",
		},
		{
			name:    "Send JSON payload",
			message: `{"type": "chat", "user": "alice", "text": "hi"}`,
		},
		{
			name:    "Send empty message",
			message: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c1, _, err := websocket.DefaultDialer.Dial(url, nil)
			if err != nil {
				t.Fatalf("Dial c1 failed: %v", err)
			}
			defer c1.Close()

			c2, _, err := websocket.DefaultDialer.Dial(url, nil)
			if err != nil {
				t.Fatalf("Dial c2 failed: %v", err)
			}
			defer c2.Close()

			// Set a read deadline for c2 to prevent hanging if message is not received
			c2.SetReadDeadline(time.Now().Add(2 * time.Second))

			// Send message from client 1
			err = c1.WriteMessage(websocket.TextMessage, []byte(tt.message))
			if err != nil {
				t.Fatalf("WriteMessage failed: %v", err)
			}

			// Read message from client 2 (broadcast)
			_, p, err := c2.ReadMessage()
			if err != nil {
				t.Fatalf("ReadMessage failed: %v", err)
			}

			if string(p) != tt.message {
				t.Errorf("got %q, want %q", string(p), tt.message)
			}
		})
	}

	t.Run("Sender isolation", func(t *testing.T) {
		c1, _, err := websocket.DefaultDialer.Dial(url, nil)
		if err != nil {
			t.Fatalf("Dial c1 failed: %v", err)
		}
		defer c1.Close()

		c2, _, err := websocket.DefaultDialer.Dial(url, nil)
		if err != nil {
			t.Fatalf("Dial c2 failed: %v", err)
		}
		defer c2.Close()

		msg := "private message"
		if err := c1.WriteMessage(websocket.TextMessage, []byte(msg)); err != nil {
			t.Fatalf("WriteMessage failed: %v", err)
		}

		// Client 1 should NOT receive its own message. 
		// We set a short deadline and expect a timeout.
		c1.SetReadDeadline(time.Now().Add(100 * time.Millisecond))
		_, _, err = c1.ReadMessage()
		if err == nil {
			t.Error("Sender received its own message, but it should be isolated")
		}

		// Confirm that Client 2 DID receive it to ensure the broadcast worked
		c2.SetReadDeadline(time.Now().Add(100 * time.Millisecond))
		_, p, err := c2.ReadMessage()
		if err != nil {
			t.Fatalf("Receiver failed to get broadcast: %v", err)
		}
		if string(p) != msg {
			t.Errorf("Receiver got %q, want %q", string(p), msg)
		}
	})
}
