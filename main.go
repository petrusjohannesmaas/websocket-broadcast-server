package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <start|connect>\n", os.Args[0])
		os.Exit(1)
	}

	command := os.Args[1]
	os.Args = append(os.Args[:1], os.Args[2:]...)

	port := flag.Int("port", 8080, "port to listen on")
	host := flag.String("host", "localhost", "host to connect to")
	flag.Parse()

	if command == "start" {
		server := NewServer(*port)
		log.Printf("Starting server on port %d", *port)
		log.Fatal(server.Start())
	} else if command == "connect" {
		reader := bufio.NewReader(os.Stdin)
		url := fmt.Sprintf("ws://%s:%d/ws", *host, *port)
		client := NewClient(url, "", reader)
		if err := client.Connect(); err != nil {
			log.Fatalf("Connection failed: %v", err)
		}
	} else {
		fmt.Fprintf(os.Stderr, "Usage: %s <start|connect>\n", os.Args[0])
		os.Exit(1)
	}
}
