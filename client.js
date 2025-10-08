// client.js
const WebSocket = require("ws");
const readline = require("readline");

function startClient(serverUrl = "ws://localhost:8080") {
    const ws = new WebSocket(serverUrl);
    let rl;
    let nickname = process.env.NICKNAME || "anonymous";

    ws.on("open", () => {
        console.log(`Connected to ${serverUrl}`);

        rl = readline.createInterface({
            input: process.stdin,
            output: process.stdout,
        });

        rl.question("Enter your nickname: ", (name) => {
            if (name.trim()) nickname = name.trim();

            rl.setPrompt(`${nickname}> `);
            rl.prompt();

            rl.on("line", (msg) => {
                try {
                    ws.send(`${nickname}: ${msg}`);
                } catch (err) {
                    console.error("Send failed:", err);
                }
                rl.prompt();
            });
        });
    });

    ws.on("message", (data) => console.log(`Received: ${data}`));
    ws.on("error", (err) => console.error("Connection error:", err));
    ws.on("close", () => {
        console.log("Disconnected.");
        process.exit(0);
    });

    process.on("SIGINT", () => {
        console.log("\nDisconnecting...");
        ws.close(1000, "Client shutting down");
        rl?.close();
        process.exit(0);
    });
}

module.exports = { startClient };
