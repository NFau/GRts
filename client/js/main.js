import WebSocketHandler from './websocket';

class GameManager {
    constructor() {
        this.webSocketHandler = new WebSocketHandler();
    }

    start() {
        this.webSocketHandler.connect();
    }
}

let gm = new GameManager;
gm.start()
