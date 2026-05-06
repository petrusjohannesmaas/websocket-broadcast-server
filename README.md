# WebSocket Broadcast Server

## 📡 Overview

A simple CLI-based WebSocket broadcast server built with Go. This tool lets you:

* **Start a WebSocket server**
* **Connect multiple interactive clients via terminal**
* **Broadcast messages between connected clients in real-time**
* **Gracefully handle disconnections with `Ctrl+C`**

> 💡 **Prerequisites:** Make sure you have Go installed (version 1.21 or later).

## 🧠 How It Works

| Command                    | Description                      |
| -------------------------- | -------------------------------- |
| `broadcast-server start`   | Starts the WebSocket server      |
| `broadcast-server connect` | Connects a terminal-based client |

All clients connect to the server and send messages that are broadcast to all other connected clients.

### Clone the Repository

```bash
git clone https://github.com/petrusjohannesmaas/websocket-broadcast-server
cd websocket-broadcast-server
```

### Build the Binary

```bash
go build -o broadcast-server .
```

This will create a `broadcast-server` binary in the current directory.

### (Optional) Install Globally

```bash
go install .
```

This will install the binary to `$GOPATH/bin` (usually `~/go/bin`), making it available from anywhere.

## 📦 Usage

**Start the WebSocket Server:**

```bash
./broadcast-server start
```

* Starts the server on `127.0.0.1:8080` (secure default, only accessible from local machine)
* Specify port with `--port` (e.g., `--port 9000`)
* Use `--remote` to allow external connections (binds to `0.0.0.0`)
* Logs client connections and broadcasts

**Connect Clients:**

Run these commands in separate terminals to simulate multiple clients.

```bash
./broadcast-server connect
```

* Prompts for a nickname
* Lets you send messages interactively
* Messages are broadcast to all other connected clients

**Connect to a Remote Server:**

```bash
./broadcast-server connect --host 192.168.1.100 --port 8080
```

* Use `--host` to specify the server's IP address
* Use `--port` to specify a custom port

## 🛑 Graceful Shutdown

* Press `Ctrl+C` in any terminal to disconnect
* Server will notify all clients and close connections

## 🧪 Example Session

**Terminal A (Server):**

```bash
./broadcast-server start
# Output: Starting server on 127.0.0.1:8080
# Output: Broadcast server running on port 8080
```

**Terminal B (Client 1):**

```bash
./broadcast-server connect
Enter your nickname: Alice
Alice > Hello!
```

**Terminal C (Client 2):**

```bash
./broadcast-server connect
Enter your nickname: Bob
Bob > Hey Alice!
# Terminal B sees: Received: Bob: Hey Alice!
```

## 📈 Future Enhancements

* Add secure WebSocket (`wss://`) support
* Add user authentication
* Add message history or logging
* Add Docker support
* Add configuration file support
* Add REST API for server management

## 📄 License

MIT License © [Petrus Johannes Maas](https://github.com/petrusjohannesmaas)
