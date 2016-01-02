package main

import (
    "GRts/server/protocol"
    "github.com/golang/protobuf/proto"
)

/* Message */
type MessageFactory struct {}

func (this MessageFactory) serialize(eventType string, data []byte) []byte {
    message := &grtsproto.Message {
        Type: eventType,
        Data: data,
    }
    serializedMessage, err := proto.Marshal(message)
    if err != nil { Error.Println("Problem marshalling network Message", err) }
    return serializedMessage
}

func (this MessageFactory) Connected(player *Player) ([]byte) {
    connected := &grtsproto.Connected {
        PlayerId: player.id,
        PlayerLogin: player.login,
    }
    data, err := proto.Marshal(connected)
    if err != nil { Error.Println("Problem marshalling proto message", err) }
    return this.serialize("CONNECTED", data)
}

func (this MessageFactory) JoinedGame(game *Game) ([]byte) {
    message := &grtsproto.JoinedGame {
        GameId: game.id,
        GameName: game.name,
    }
    for _, player := range game.players {
        message.Players = append(message.Players, &grtsproto.PlayerDefinition {
            Id: player.id,
            Login: player.login,
        })
    }
    data, err := proto.Marshal(message)
    if err != nil { Error.Println("Problem marshalling proto message", err) }
    return this.serialize("JOINED_GAME", data)
}

func (this MessageFactory) PlayerJoinedGame(player *Player) ([]byte) {
    pjg := &grtsproto.PlayerJoinedGame {
        Player: &grtsproto.PlayerDefinition {
            Id: player.id,
            Login: player.login,
        },
    }
    data, err := proto.Marshal(pjg)
    if err != nil { Error.Println("Problem marshalling proto message", err) }
    return this.serialize("PLAYER_JOINED_GAME", data)
}

func (this MessageFactory) PlayerLeftGame(playerId uint32) ([]byte) {
    message := &grtsproto.PlayerLeftGame {
        PlayerId: playerId,
    }
    data, err := proto.Marshal(message)
    if err != nil { Error.Println("Problem marshalling proto message", err) }
    return this.serialize("PLAYER_LEFT_GAME", data)
}

func (this MessageFactory) PlayerChangedLogin(player *Player) ([]byte) {
    message := &grtsproto.PlayerChangedLogin {
        Player: &grtsproto.PlayerDefinition {
            Id: player.id,
            Login: player.login,
        },
    }
    data, err := proto.Marshal(message)
    if err != nil { Error.Println("Problem marshalling proto message", err) }
    return this.serialize("PLAYER_CHANGED_LOGIN", data)
}

var GRTSMessage = MessageFactory{}
