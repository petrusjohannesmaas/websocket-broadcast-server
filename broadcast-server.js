#!/usr/bin/env node

const { startServer } = require("./server");
const { startClient } = require("./client");

const args = process.argv.slice(2);
const command = args[0];

if (command === "start") {
    startServer();
} else if (command === "connect") {
    startClient();
} else {
    console.error("Usage: broadcast-server <start|connect>");
    process.exit(1);
}
