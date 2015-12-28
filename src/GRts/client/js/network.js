import pb from "protobufjs"

// Protobuf encoder/decoder
// Basic version, will probably need to handle complex (de)serialization in the future
class NetworkSerializer {
    constructor() {
        this._builder = pb.loadProtoFile("/resources/protocol/lobby.proto");
        this._grtsproto = this._builder.build('grtsproto');

        this._protomap = {
            "CONNECTED": this._grtsproto.Connected
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
        let protoer = this._getProtoer('Decode', message.type)
        if (!protoer) { return null; }
        message.data = protoer.decode(message.data);
        return message;
    }

    encode(message) {
        let protoer = this._getProtoer('Encode', message.type)
        if (!protoer) { return null; }
        message.data = new protoer(message.data).encode;
        return message;
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
        const message = JSON.parse(event.data);
        const deserializedMessage = this._serializer.decode(message);
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

    send(message) {
        const serializedMessage = this._serializer.encode(message);

    }
}

export default NetworkManager;
