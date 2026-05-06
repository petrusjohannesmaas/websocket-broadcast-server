package main

import (
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

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type Server struct {
	Port       int
	clients   map[*websocket.Conn]bool
	mu        sync.RWMutex
	upgrader  websocket.Upgrader
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
	httpAddr := fmt.Sprintf("0.0.0.0:%d", s.Port)
	log.Printf("Broadcast server running at ws://localhost:%d", s.Port)

	http.HandleFunc("/ws", s.handleWebSocket)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGINT)

	go func() {
		<-sigChan
		log.Println("\nServer is shutting down...")
		s.shutdown()
	}()

	return http.ListenAndServe(httpAddr, nil)
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
		conn.Close()
		s.mu.Lock()
		delete(s.clients, conn)
		s.mu.Unlock()
		log.Println("Client disconnected.")
	}()

	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			if _, ok := err.(*websocket.CloseError); ok {
				break
			}
			log.Printf("Client error: %v", err)
			break
		}

		log.Printf("Broadcasting: %s", message)
		s.broadcast(message, conn, messageType)
	}
}

func (s *Server) broadcast(message []byte, sender *websocket.Conn, messageType int) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for client := range s.clients {
		if client != sender && client.WriteMessage(messageType, message) != nil {
			log.Printf("Error broadcasting to client")
		}
	}
}

func (s *Server) shutdown() {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for client := range s.clients {
		client.WriteControl(websocket.CloseMessage, []byte("Server shutting down"), time.Now().Add(5*time.Second))
		client.Close()
	}
}

var globalServer *Server

func GetTestServer() *Server {
	return globalServer
}