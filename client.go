package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

type Client struct {
	ServerURL string
	Nickname string
	Conn    *websocket.Conn
	Reader  *bufio.Reader
}

func NewClient(serverURL, nickname string, reader *bufio.Reader) *Client {
	return &Client{
		ServerURL: serverURL,
		Nickname: nickname,
		Reader:  reader,
	}
}

func (c *Client) Connect() error {
	if c.Nickname == "" {
		c.Nickname = os.Getenv("NICKNAME")
		if c.Nickname == "" {
			c.Nickname = "Anonymous"
		}
	}

	log.Printf("Connecting to %s", c.ServerURL)

	conn, _, err := websocket.DefaultDialer.Dial(c.ServerURL, nil)
	if err != nil {
		return fmt.Errorf("connection failed: %v", err)
	}
	c.Conn = conn

	log.Printf("Connected to %s", c.ServerURL)

	go c.readLoop()
	c.writeLoop()

	return nil
}

func (c *Client) readLoop() {
	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			log.Printf("Disconnected: %v", err)
			c.Conn.Close()
			os.Exit(0)
			return
		}
		log.Printf("Received: %s", message)
	}
}

func (c *Client) writeLoop() {
	fmt.Print("Enter your nickname: ")
	nickname, err := c.Reader.ReadString('\n')
	if err != nil {
		return
	}
	nickname = strings.TrimSpace(nickname)
	if nickname != "" {
		c.Nickname = nickname
	}

	for {
		fmt.Printf("%s > ", c.Nickname)
		message, err := c.Reader.ReadString('\n')
		if err != nil {
			break
		}
		message = strings.TrimSpace(message)
		if message == "" {
			continue
		}

		fullMessage := fmt.Sprintf("%s: %s", c.Nickname, message)
		err = c.Conn.WriteMessage(websocket.TextMessage, []byte(fullMessage))
		if err != nil {
			log.Printf("Send failed: %v", err)
			break
		}
	}
}

func (c *Client) Close() error {
	if c.Conn != nil {
		return c.Conn.WriteControl(websocket.CloseMessage, []byte("Client shutting down"), time.Now().Add(5*time.Second))
	}
	return nil
}