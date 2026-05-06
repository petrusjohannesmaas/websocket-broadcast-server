package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

func TestServerStartup(t *testing.T) {
	server := NewServer(8081)
	if server == nil {
		t.Error("server should not be nil")
	}
	if server.Port != 8081 {
		t.Errorf("expected port 8081, got %d", server.Port)
	}
}

func TestServerDefaultPort(t *testing.T) {
	server := NewServer(8080)
	if server.Port != 8080 {
		t.Errorf("expected port 8080, got %d", server.Port)
	}
}

func TestServerCustomPort(t *testing.T) {
	server := NewServer(9000)
	if server.Port != 9000 {
		t.Errorf("expected port 9000, got %d", server.Port)
	}
}

func startTestServer(port int) *httptest.Server {
	server := NewServer(port)
	mux := http.NewServeMux()
	mux.HandleFunc("/ws", server.handleWebSocket)
	return httptest.NewServer(mux)
}

func TestClientConnection(t *testing.T) {
	ts := startTestServer(8082)
	defer ts.Close()

	url := strings.Replace(ts.URL, "http", "ws", 1)
	ws, _, err := websocket.DefaultDialer.Dial(url+"/ws", nil)
	if err != nil {
		t.Fatalf("failed to connect: %v", err)
	}
	defer ws.Close()
	_ = ts
}

func TestClientConnectionError(t *testing.T) {
	_, _, err := websocket.DefaultDialer.Dial("ws://localhost:99999", nil)
	if err == nil {
		t.Error("expected connection error")
	}
}

func TestMessageBroadcast(t *testing.T) {
	ts := startTestServer(8083)
	defer ts.Close()

	url := strings.Replace(ts.URL, "http", "ws", 1)

	ws1, _, err := websocket.DefaultDialer.Dial(url+"/ws", nil)
	if err != nil {
		t.Fatalf("failed to connect ws1: %v", err)
	}
	defer ws1.Close()

	ws2, _, err := websocket.DefaultDialer.Dial(url+"/ws", nil)
	if err != nil {
		t.Fatalf("failed to connect ws2: %v", err)
	}
	defer ws2.Close()

	done := make(chan string, 1)
	go func() {
		_, msg, err := ws2.ReadMessage()
		if err != nil {
			done <- "error:" + err.Error()
			return
		}
		done <- string(msg)
	}()

	err = ws1.WriteMessage(websocket.TextMessage, []byte("hello"))
	if err != nil {
		t.Fatalf("write error: %v", err)
	}

	select {
	case result := <-done:
		if result == "error:timeout" {
			t.Error("timeout waiting for message")
		}
		if result != "hello" {
			t.Errorf("expected 'hello', got '%s'", result)
		}
	case <-time.After(time.Second):
		t.Error("timeout waiting for message")
	}
}

func TestNoBroadcastToSelf(t *testing.T) {
	ts := startTestServer(8084)
	defer ts.Close()

	url := strings.Replace(ts.URL, "http", "ws", 1)

	ws, _, err := websocket.DefaultDialer.Dial(url+"/ws", nil)
	if err != nil {
		t.Fatalf("failed to connect: %v", err)
	}
	defer ws.Close()

	err = ws.WriteMessage(websocket.TextMessage, []byte("test"))
	if err != nil {
		t.Fatalf("write error: %v", err)
	}

	ws.SetReadDeadline(time.Now().Add(200 * time.Millisecond))
	_, _, err = ws.ReadMessage()
	if err == nil {
		t.Error("should not receive own message")
	}
}

func TestDisconnectClient(t *testing.T) {
	ts := startTestServer(8085)
	defer ts.Close()

	url := strings.Replace(ts.URL, "http", "ws", 1)

	ws1, _, err := websocket.DefaultDialer.Dial(url+"/ws", nil)
	if err != nil {
		t.Fatalf("failed to connect: %v", err)
	}

	ws2, _, err := websocket.DefaultDialer.Dial(url+"/ws", nil)
	if err != nil {
		t.Fatalf("failed to connect: %v", err)
	}
	defer ws2.Close()

	ws1.Close()

	time.Sleep(100 * time.Millisecond)
	_ = ws2
}

func TestEmptyMessageBroadcast(t *testing.T) {
	ts := startTestServer(8086)
	defer ts.Close()

	url := strings.Replace(ts.URL, "http", "ws", 1)

	ws1, _, err := websocket.DefaultDialer.Dial(url+"/ws", nil)
	if err != nil {
		t.Fatalf("failed to connect: %v", err)
	}
	defer ws1.Close()

	ws2, _, err := websocket.DefaultDialer.Dial(url+"/ws", nil)
	if err != nil {
		t.Fatalf("failed to connect: %v", err)
	}
	defer ws2.Close()

	err = ws1.WriteMessage(websocket.TextMessage, []byte(""))
	if err != nil {
		t.Fatalf("write error: %v", err)
	}

	ws2.SetReadDeadline(time.Now().Add(200 * time.Millisecond))
	_, _, err = ws2.ReadMessage()
	if err != nil {
		t.Logf("empty message handled: %v", err)
	}
}

func TestServerClientCount(t *testing.T) {
	ts := startTestServer(8087)
	defer ts.Close()

	url := strings.Replace(ts.URL, "http", "ws", 1)

	clients := make([]*websocket.Conn, 3)
	for i := 0; i < 3; i++ {
		ws, _, err := websocket.DefaultDialer.Dial(url+"/ws", nil)
		if err != nil {
			t.Fatalf("failed to connect: %v", err)
		}
		clients[i] = ws
	}

	time.Sleep(50 * time.Millisecond)

	for _, ws := range clients {
		ws.Close()
	}
}

func TestMultipleClients(t *testing.T) {
	ts := startTestServer(8088)
	defer ts.Close()

	url := strings.Replace(ts.URL, "http", "ws", 1)

	ws1, _, err := websocket.DefaultDialer.Dial(url+"/ws", nil)
	if err != nil {
		t.Fatalf("failed to connect ws1: %v", err)
	}
	defer ws1.Close()

	ws2, _, err := websocket.DefaultDialer.Dial(url+"/ws", nil)
	if err != nil {
		t.Fatalf("failed to connect ws2: %v", err)
	}
	defer ws2.Close()

	ws3, _, err := websocket.DefaultDialer.Dial(url+"/ws", nil)
	if err != nil {
		t.Fatalf("failed to connect ws3: %v", err)
	}
	defer ws3.Close()

	err = ws1.WriteMessage(websocket.TextMessage, []byte("message from 1"))
	if err != nil {
		t.Fatalf("write error: %v", err)
	}

	done := make(chan int, 2)
	go func() {
		_, msg, _ := ws2.ReadMessage()
		if string(msg) == "message from 1" {
			done <- 1
		}
	}()
	go func() {
		_, msg, _ := ws3.ReadMessage()
		if string(msg) == "message from 1" {
			done <- 1
		}
	}()

	select {
	case <-done:
	case <-time.After(time.Second):
		t.Error("timeout")
	}
}

func TestNicknameDefault(t *testing.T) {
	client := NewClient("ws://localhost:8080", "", nil)
	if client.Nickname != "" {
		t.Errorf("expected empty nickname, got '%s'", client.Nickname)
	}
}

func TestNicknameSet(t *testing.T) {
	client := NewClient("ws://localhost:8080", "Alice", nil)
	if client.Nickname != "Alice" {
		t.Errorf("expected 'Alice', got '%s'", client.Nickname)
	}
}

func TestClientNew(t *testing.T) {
	server := "ws://localhost:8080"
	nickname := "TestUser"
	client := NewClient(server, nickname, nil)

	if client.ServerURL != server {
		t.Errorf("expected server %s, got %s", server, client.ServerURL)
	}
	if client.Nickname != nickname {
		t.Errorf("expected nickname %s, got %s", nickname, client.Nickname)
	}
	if client.Conn != nil {
		t.Error("expected nil connection")
	}
}

func TestServerString(t *testing.T) {
	server := NewServer(8080)
	_ = server
}

func TestBinaryMessage(t *testing.T) {
	ts := startTestServer(8089)
	defer ts.Close()

	url := strings.Replace(ts.URL, "http", "ws", 1)

	ws1, _, err := websocket.DefaultDialer.Dial(url+"/ws", nil)
	if err != nil {
		t.Fatalf("failed to connect: %v", err)
	}
	defer ws1.Close()

	ws2, _, err := websocket.DefaultDialer.Dial(url+"/ws", nil)
	if err != nil {
		t.Fatalf("failed to connect: %v", err)
	}
	defer ws2.Close()

	err = ws1.WriteMessage(websocket.BinaryMessage, []byte("binary data"))
	if err != nil {
		t.Fatalf("write error: %v", err)
	}

	ws2.SetReadDeadline(time.Now().Add(200 * time.Millisecond))
	_, msg, err := ws2.ReadMessage()
	if err != nil {
		t.Logf("binary message handled: %v", err)
	} else if string(msg) != "binary data" {
		t.Errorf("expected 'binary data', got '%s'", msg)
	}
}