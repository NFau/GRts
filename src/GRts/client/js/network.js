import pb from "protobufjs"

// Protobuf encoder/decoder
// Basic version, will probably need to handle complex (de)serialization in the future
class NetworkSerializer {
    constructor() {
        this._builder = pb.loadProtoFile("/resources/protocol/lobby.proto");
        this._grtsproto = this._builder.build('grtsproto');

        this._protomap = {
            /* Receive */
            "CONNECTED": this._grtsproto.Connected,
            "JOINED_GAME": this._grtsproto.JoinedGame,
            "PLAYER_JOINED_GAME": this._grtsproto.PlayerJoinedGame,
            "PLAYER_LEFT_GAME": this._grtsproto.PlayerLeftGame,
            "PLAYER_CHANGED_LOGIN": this._grtsproto.PlayerChangedLogin,

            /* Send */
            "CHANGE_LOGIN": this._grtsproto.ChangeLogin,
            "LOBBY_READY": this._grtsproto.LobbyReady,
        }
    }

    _getProtoer(context, type) {
        let protoer = this._protomap[type];
        if (!protoer) {
            // TODO: Add a proper log system (CAULIFLOWER? :D)
            console.error(`[Serializer:${context}] Cannot find protoer for event type ${message.type}`);
            return null;
        }
        return protoer
    }

    decode(message) {
        message = this._grtsproto.Message.decode(message);
        let protoer = this._getProtoer('Decode', message.type)
        if (!protoer) { return null; }
        message.data = protoer.decode(message.data);
        return message;
    }

    encode(type, data) {
        let protoer = this._getProtoer('Decode', type);
        if (!protoer) { return null; }
        data = new protoer(data);
        let message = new this._grtsproto.Message({
            'type': type,
            'data': data.encode()
        })
        return message.toBuffer();
    }
}


class NetworkManager {

    constructor() {
        this._ws = null;
        this._serializer = new NetworkSerializer();

        this._onMessageHandlers = [];
        this._onWSOpenHandlers = [];
        this._onWSCloseHandlers = [];
    }

    /* Private methods */
    _onMessage(event) {
        const deserializedMessage = this._serializer.decode(event.data);
        this._onMessageHandlers.forEach((handler) => {
            handler(deserializedMessage);
        });
    }

    _onWSOpen() {
        this._onWSOpenHandlers.forEach((handler) => {
            handler();
        });
    }

    _onWSClose() {
        this._onWSCloseHandlers.forEach((handler) => {
            handler();
        });
    }

    _registerHandler(list, handler) {
        if (typeof(handler) == "function") {
            list.push(handler);
        }
    }

    /* Public methods */
    connect() {
        this._ws = new WebSocket("ws://127.0.0.1:8080/socket");
        this._ws.binaryType = "arraybuffer";
        this._ws.onopen = this._onWSOpen.bind(this);
        this._ws.onclose = this._onWSClose.bind(this);
        this._ws.onmessage = this._onMessage.bind(this);
    }

    registerOnMessageHandler(handler) {
        this._registerHandler(this._onMessageHandlers, handler);
    }

    unregisterOnMessageHandler(handler) {
        this._onMessageHandlers.remove(handler);
    }

    registerOnWSOpenHandler(handler) {
        this._registerHandler(this._onWSOpenHandlers, handler);
    }

    unregisterOnWSOpenHandler(handler) {
        this._onWSOpenHandlers.remove(handler);
    }

    registerOnWSCloseHandler(handler) {
        this._registerHandler(this._onWSCloseHandlers, handler);
    }

    unregisterOnWSCloseHandler(handler) {
        this._onWSCloseHandlers.remove(handler);
    }

    send(type, data) {
        console.log("Send", type, data);
        const serializedMessage = this._serializer.encode(type, data);
        this._ws.send(serializedMessage);
    }
}

export default NetworkManager;
