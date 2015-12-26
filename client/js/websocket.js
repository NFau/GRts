class WebSocketHandler {
    constructor() {
        this.webSocket = null;
    }

    connect() {
        this.webSocket = new WebSocket("ws://127.0.0.1:8080/socket");
        this.webSocket.onopen = function (event) {
            console.log("Connection opened");
        };
        this.webSocket.onmessage = function (event) {
            console.log("Received message", event.data);
        };
        this.webSocket.onclose = function (event) {
            console.log("Connection closed", event.data);
        }
    }
}

export default WebSocketHandler;
