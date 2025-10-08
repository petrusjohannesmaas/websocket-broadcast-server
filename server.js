// server.js
const WebSocket = require("ws");

function startServer(port = 8080) {
    const wss = new WebSocket.Server({ port });

    console.log(`Broadcast server running at ws://localhost:${port}`);

    wss.on("connection", (ws) => {
        console.log("Client connected.");

        ws.on("message", (message) => {
            console.log(`Broadcasting: ${message}`);
            for (const client of wss.clients) {
                if (client !== ws && client.readyState === WebSocket.OPEN) {
                    client.send(message.toString());
                }
            }
        });

        ws.on("error", (err) => console.error("Client error:", err));
        ws.on("close", () => console.log("Client disconnected."));
    });

    process.on("SIGINT", () => {
        console.log("\nGracefully shutting down server...");
        for (const client of wss.clients) {
            if (client.readyState === WebSocket.OPEN) {
                client.close(1001, "Server shutting down");
            }
        }
        wss.close(() => {
            console.log("WebSocket server closed.");
            process.exit(0);
        });
    });
}

module.exports = { startServer };
