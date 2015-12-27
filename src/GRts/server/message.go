package main

import (
//    "encoding/json"
    "GRts/server/protocol"
    "github.com/golang/protobuf/proto"
)

/* Message */
type MessageFactory struct {}

func (this MessageFactory) Connected(player *Player) ([]byte) {
    message := &grtsproto.Connected {
        Event: "CONNECTED",
        PlayerId: player.id,
    }
    data, err := proto.Marshal(message)
    if err != nil { Error.Println("Problem marshalling message", err) }
    return data
}

var GRTSMessage = MessageFactory{}
