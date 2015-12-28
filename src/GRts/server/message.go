package main

import (
    "encoding/json"
    "GRts/server/protocol"
    "github.com/golang/protobuf/proto"
)

/* NetworkMessage */
type NetworkMessage struct {
    Type    string `json:"type"`
    Data    []byte `json:"data"`
}

/* Message */
type MessageFactory struct {}

func (this MessageFactory) serialize(eventType string, data []byte) []byte {
    message := NetworkMessage{
        Type: eventType,
        Data: data,
    }
    serializedMessage, err := json.Marshal(message)
    if err != nil { Error.Println("Problem marshalling NetworkMessage", err) }
    return serializedMessage
}

func (this MessageFactory) Connected(player *Player) ([]byte) {
    message := &grtsproto.Connected {
        PlayerId: player.id,
        PlayerLogin: player.login,
    }
    data, err := proto.Marshal(message)
    if err != nil { Error.Println("Problem marshalling proto message", err) }

    return this.serialize("CONNECTED", data)
}

var GRTSMessage = MessageFactory{}
