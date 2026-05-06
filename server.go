package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/gorilla/websocket"
)

type Server struct {
	Port     int
	clients  map[*websocket.Conn]bool
	mu       sync.Mutex
	upgrader websocket.Upgrader
	http     *http.Server
}

func NewServer(port int) *Server {
	return &Server{
		Port:    port,
		clients: make(map[*websocket.Conn]bool),
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool { return true },
		},
	}
}

func (s *Server) Start() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/ws", s.handleWebSocket)

	s.http = &http.Server{
		Addr:    fmt.Sprintf("0.0.0.0:%d", s.Port),
		Handler: mux,
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigChan
		log.Println("Server shutting down...")
		s.shutdown()
	}()

	log.Printf("Broadcast server running at ws://localhost:%d", s.Port)
	if err := s.http.ListenAndServe(); err != http.ErrServerClosed {
		return err
	}
	return nil
}

func (s *Server) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("upgrade error: %v", err)
		return
	}

	s.mu.Lock()
	s.clients[conn] = true
	s.mu.Unlock()

	log.Println("Client connected.")
	go s.handleMessages(conn)
}

func (s *Server) handleMessages(conn *websocket.Conn) {
	defer func() {
		s.mu.Lock()
		delete(s.clients, conn)
		s.mu.Unlock()
		conn.Close()
		log.Println("Client disconnected.")
	}()

	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			if _, ok := err.(*websocket.CloseError); !ok {
				log.Printf("Client error: %v", err)
			}
			return
		}
		log.Printf("Broadcasting: %s", message)
		s.broadcast(message, conn, messageType)
	}
}

func (s *Server) broadcast(message []byte, sender *websocket.Conn, messageType int) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for client := range s.clients {
		if client != sender {
			if err := client.WriteMessage(messageType, message); err != nil {
				log.Printf("broadcast error: %v", err)
			}
		}
	}
}

func (s *Server) shutdown() {
	s.mu.Lock()
	for client := range s.clients {
		client.WriteControl(websocket.CloseMessage, []byte("Server shutting down"), time.Now().Add(time.Second))
		client.Close()
	}
	s.mu.Unlock()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	s.http.Shutdown(ctx)
}
