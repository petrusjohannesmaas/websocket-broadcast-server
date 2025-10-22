### 🗼 Part 1: The Server (`server.js`)
> "Let’s start with the heart of the system — the server."

```js
const WebSocket = require("ws");
```
> "We’re using the `ws` library to handle WebSocket connections. It’s lightweight and perfect for real-time communication."

```js
function startServer(port = 8080) {
    const wss = new WebSocket.Server({ port });
```
> "Here we create a WebSocket server listening on port 8080 by default. You can change this if needed."

```js
    wss.on("connection", (ws) => {
        console.log("Client connected.");
```
> "Every time a client connects, we log it. Then we listen for messages from that client."

```js
        ws.on("message", (message) => {
            console.log(`Broadcasting: ${message}`);
```
> "When a message comes in, we broadcast it to all other connected clients."

```js
            for (const client of wss.clients) {
                if (client !== ws && client.readyState === WebSocket.OPEN) {
                    client.send(message.toString());
                }
            }
```
> "This loop ensures the message is sent only to other clients that are still connected."

```js
        ws.on("error", (err) => console.error("Client error:", err));
        ws.on("close", () => console.log("Client disconnected."));
```
> "We also handle errors and disconnections gracefully."

```js
    process.on("SIGINT", () => {
        console.log("\nServer is shutting down ...");
```
> "And finally, we listen for Ctrl+C to shut down the server cleanly and notify all clients."

---

### 💬 Part 2: The Client (`client.js`)
> "Now let’s look at the client — this is what users run to connect and chat."

```js
const WebSocket = require("ws");
const readline = require("readline");
```
> "We use `ws` again for the connection, and `readline` to capture user input from the terminal."

```js
    let nickname = process.env.NICKNAME || "Anonymous";
```
> "Each client gets a nickname, either from an environment variable or defaulting to 'Anonymous'."

```js
    rl.question("Enter your nickname: ", (name) => {
        if (name.trim()) nickname = name.trim();
```
> "We prompt the user to enter a nickname interactively."

```js
        rl.on("line", (msg) => {
            ws.send(`${nickname}: ${msg}`);
```
> "Every line the user types gets sent to the server, prefixed with their nickname."

```js
    ws.on("message", (data) => console.log(`Received: ${data}`));
```
> "Incoming messages from other users are printed to the console."

```js
    process.on("SIGINT", () => {
        console.log("\nDisconnecting...");
        ws.close(1000, "Client shutting down");
        rl?.close();
        process.exit(0);
    });
```
> "Just like the server, the client shuts down cleanly when you hit Ctrl+C."

---

### 🚀 Part 3: The Entry Point (`broadcast-server.js`)
> "Finally, we have the entry script that ties everything together."

```js
const args = process.argv.slice(2);
const command = args[0];
```
> "We grab the command-line argument to decide what to do."

```js
if (command === "start") {
    startServer();
} else if (command === "connect") {
    startClient();
```
> "If you run `broadcast-server start`, it launches the server. If you run `broadcast-server connect`, it launches a client."

```js
} else {
    console.error("Usage: broadcast-server <start|connect>");
    process.exit(1);
}
```
> "And if you mess up the command, it shows you how to use it properly."

---