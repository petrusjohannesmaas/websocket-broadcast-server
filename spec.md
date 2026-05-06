# WebSocket Broadcast Server - Specification (Go)

## Overview
A CLI-based WebSocket broadcast server built with Go enabling real-time message broadcasting between multiple connected clients.

## Current Features (to be ported to Go)

### 1. WebSocket Server (`server.go`)
- Creates WebSocket server on configurable port (default: 8080)
- Accepts client connections
- Broadcasts messages from any client to all other connected clients
- Handles client disconnections gracefully
- Handles server shutdown (SIGINT) and notifies all clients

### 2. Interactive Client (`client.go`)
- Connects to WebSocket server (default: ws://localhost:8080)
- Prompts for nickname interactively (or uses NICKNAME env var, defaults to "Anonymous")
- Sends messages prefixed with nickname (format: `nickname: message`)
- Displays received messages from other clients
- Handles disconnection and shutdown gracefully

### 3. CLI Entry Point (`main.go`)
- `broadcast-server start` - Starts the WebSocket server
- `broadcast-server connect` - Connects an interactive client
- Error message on invalid command usage

## Required Tests

### Unit Tests

#### Server Tests (`server_test.go`)
1. **Server Startup**
   - Should start server on specified port
   - Should start server on default port (8080) when no port specified
   - Should reject starting on already-in-use port

2. **Client Connection Handling**
   - Should accept new client connections
   - Should track connected clients count
   - Should handle client disconnection (remove from clients set)
   - Should emit events on connection/disconnection

3. **Message Broadcasting**
   - Should broadcast message to all connected clients except sender
   - Should not broadcast to disconnected clients
   - Should not broadcast to client with non-OPEN readyState
   - Should handle empty message broadcasting

4. **Server Shutdown**
   - Should close all client connections on SIGINT
   - Should send close code 1001 on server shutdown
   - Should exit cleanly after closing

#### Client Tests (`client_test.go`)
1. **Client Connection**
   - Should connect to server with default URL (ws://localhost:8080)
   - Should connect to server with custom URL
   - Should handle connection errors

2. **Nickname Handling**
   - Should use NICKNAME env var when available
   - Should use "Anonymous" when NICKNAME not set and no input
   - Should use provided nickname from prompt input

3. **Message Sending**
   - Should send messages prefixed with nickname
   - Should handle send errors gracefully
   - Should not crash on send failure

4. **Message Receiving**
   - Should log received messages
   - Should handle malformed messages

5. **Client Shutdown**
   - Should close connection on SIGINT
   - Should send close code 1000 on client shutdown
   - Should close readline interface on disconnect

### Integration Tests (`integration_test.go`)

1. **Multi-Client Broadcasting**
   - Should allow multiple clients to connect simultaneously
   - Should broadcast message from Client A to Client B and C (but not A)
   - Should handle message ordering correctly

2. **Client Disconnect Scenario**
   - Should continue broadcasting when one client disconnects
   - Should remove disconnected client from broadcast list

3. **Nickname Prefix Verification**
   - Should verify messages include correct nickname prefix
   - Should handle special characters in nicknames

4. **Graceful Server Shutdown with Connected Clients**
   - Should notify all clients when server shuts down
   - Should close all connections with appropriate close codes

### Edge Cases (`edge_cases_test.go`)

1. **Rapid Connect/Disconnect**
   - Should handle rapid client connections
   - Should handle rapid client disconnections

2. **Invalid Messages**
   - Should handle non-string messages
   - Should handle very large messages

3. **Server Not Running**
   - Client should handle connection to non-existent server
   - Client should emit appropriate error events

## Implementation Requirements

### Project Structure
```
websocket-broadcast-server/
├── main.go          # CLI entry point
├── server.go       # WebSocket server implementation
├── client.go      # Interactive client implementation
├── go.mod        # Go module file
└── *_test.go     # Test files

```

### Dependencies
- `github.com/gorilla/websocket` - WebSocket implementation
- Standard library: `bufio`, `flag`, `log`, `os`, `os/signal`, `fmt`

### Refactoring for Testability

1. **Server Module (`server.go`)**
   - Return the `http.Server` instance for testing purposes
   - Use interfaces for dependency injection
   - Extract message broadcasting logic into testable function
   - Use channels for client management

2. **Client Module (`client.go`)**
   - Accept `bufio.Reader` as parameter for testing
   - Accept websocket connection interface for testing
   - Return reference for test control

### Test Infrastructure

1. **Testing Framework**: Go's built-in `testing` package
   - Run with: `go test -v`
   - Coverage: `go test -cover`

2. **WebSocket Testing**
   - Use `gorilla/websocket` test helper utilities
   - Create helper functions for spawning test clients
   - Use goroutines and channels for async verification

3. **Test File Structure**
   ```
   server_test.go
   client_test.go
   integration_test.go
   edge_cases_test.go
   ```

### Example Test Pattern

```go
// server_test.go
package main

import (
    "net/http"
    "testing"
    "time"
    
    "github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

func TestServerBroadcast(t *testing.T) {
    server := startServer(8081)
    defer server.Close()
    
    // Connect two clients
    ws1, _, err := websocket.DefaultDialer.Dial("ws://localhost:8081", nil)
    if err != nil {
        t.Fatalf("failed to dial: %v", err)
    }
    defer ws1.Close()
    
    ws2, _, err := websocket.DefaultDialer.Dial("ws://localhost:8081", nil)
    if err != nil {
        t.Fatalf("failed to dial: %v", err)
    }
    defer ws2.Close()
    
    // Send message from client 1
    err = ws1.WriteMessage(websocket.TextMessage, []byte("test message"))
    if err != nil {
        t.Fatalf("failed to write message: %v", err)
    }
    
    // Receive message on client 2
    _, msg, err := ws2.ReadMessage()
    if err != nil {
        t.Fatalf("failed to read message: %v", err)
    }
    
    if string(msg) != "test message" {
        t.Errorf("expected 'test message', got '%s'", msg)
    }
}
```

## Acceptance Criteria

- [ ] All unit tests pass (server, client)
- [ ] All integration tests pass (multi-client scenarios)
- [ ] All edge case tests pass
- [ ] Test coverage > 80%
- [ ] Code refactored for testability (dependency injection)
- [ ] `go test` runs all tests successfully
- [ ] No regressions in existing functionality
- [ ] Build succeeds with `go build`
- [ ] Binary can be run with `./broadcast-server start` and `./broadcast-server connect`

## Developer Notes

The senior developer should:
1. First create Go module and dependencies
2. Write all tests BEFORE any implementation
3. Run tests to confirm they fail (TDD red phase)
4. Implement Go code for server
5. Implement Go code for client
6. Run tests to make them pass (green phase)
7. Refactor as needed (refactor phase)
8. Ensure all tests pass before considering the task complete

**Important**: No implementation changes should be made until all tests are written and the initial failure state is confirmed.