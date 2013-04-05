ws = new WebSocket("ws://localhost:8000/websocket/ws");
if (ws == null) {
    console.log("WebSocket creation failed");
    return;
} else {
    console.log("WebSocket creation succeeded");
}
