package main

import (
    "fmt"
    "github.com/golang/protobuf/proto"
    "GRts/server/protocol"
)

/* == Player ==*/
type Player struct {
    id      uint32
    ws      *Connection
    login   string
}

/* == Player Messages == */
type PlayerMessage struct {
    playerId    uint32
    data        []byte
}

/* == Game ==*/
type Game struct {
    // Game configuration
    numberOfPlayers uint32

    // Game infos
    id              uint32
    name            string

    players         map[uint32]*Player

    newPlayer           chan *Player
    recv                chan PlayerMessage
    unregisterPlayer    chan uint32
}

func (this *Game) addPlayer(player *Player) {
    this.players[player.id] = player
    this.registerPlayerToGeneralRecv(player, this.recv, this.unregisterPlayer)
}

func (this *Game) removePlayer(id uint32) {
    // Kill recv->generalRecv routine
    this.unregisterPlayer <- id
    delete(this.players, id)
}

func (this *Game) broadcast(message []byte, excludedPlayerIds []uint32) {
    Info.Println("Broadcast")
    var excluded bool
    for id, player := range this.players {
        if len(excludedPlayerIds) > 0 {
            excluded = false
            for i := range excludedPlayerIds {
                if excludedPlayerIds[i] == id {
                    excluded = true
                    break
                }
            }
            if excluded {
                continue
            }
        }
        Info.Println("--->>> SEND")
        player.ws.send <- message
    }
}

func (this *Game) registerPlayerToGeneralRecv(player *Player, generalRecv chan<- PlayerMessage, unregister <-chan uint32) {
    go func () {
        Info.Println("Registering player recv chan")
        for {
            select {
            case userId := <- unregister:
                if userId == player.id {
                    Info.Println("Player recv channel has been kicked from the general recv chan")
                    break
                }
            case data := <- player.ws.recv:
                this.recv <- PlayerMessage{
                    playerId: player.id,
                    data: data,
                }
                if len(data) == 0 {
                    break
                }
            }
        }
        Info.Println("Player recv->generalRecv routine stopped.")
    }()
}

func (this *Game) ProcessPlayerMessage(playerId uint32, message *grtsproto.Message) {
    Info.Println("Receive message from Player", message.Type)

    if message.Type == "CHANGE_LOGIN" {
        data := &grtsproto.ChangeLogin{}
        err := proto.Unmarshal(message.Data, data)
        if err != nil { Error.Println("Error unmarshaling grtsproto.ChangeLogin") }
        this.players[playerId].login = data.Login
        this.broadcast(GRTSMessage.PlayerChangedLogin(this.players[playerId]), []uint32{})
    }
}

func (this *Game) Run(gameRunning chan<- *Game, gameStopped chan<- *Game) {
    Info.Println("Game", this.id, "init")

    this.players = make(map[uint32]*Player)
    this.recv = make(chan PlayerMessage)
    this.unregisterPlayer = make(chan uint32)
    // Waiting for player

    for {
        if uint32(len(this.players)) == this.numberOfPlayers {
            gameRunning <- this
            break;
        }

        select {
        case newPlayer := <- this.newPlayer:
            Info.Println("Game", this.id, "has a new player")
            // Notify other Players that someone joined the game
            this.addPlayer(newPlayer)
            this.broadcast(GRTSMessage.PlayerJoinedGame(newPlayer), []uint32{newPlayer.id})
            newPlayer.ws.send <- GRTSMessage.JoinedGame(this)
        case playerMessage := <- this.recv:
            message := &grtsproto.Message{}
            err := proto.Unmarshal(playerMessage.data, message)
            if err != nil { Error.Println("Error unmarshaling grtsproto.Message") }
            this.ProcessPlayerMessage(playerMessage.playerId, message)
        }
    }

    for {
        select{
        case playerMessage := <- this.recv:
            message := &grtsproto.Message{}
            err := proto.Unmarshal(playerMessage.data, message)
            if err != nil { Error.Println("Error unmarshaling grtsproto.Message") }
            this.ProcessPlayerMessage(playerMessage.playerId, message)
        }
    }
}

/* == Game Manager== */
type GameManager struct {
    newConnection chan *Connection

    // Map of games (bool: inGame?)
    games       map[*Game]bool

    nextGameId  uint32
    nextPlayerId uint32

    gameRunning chan *Game
    gameStopped chan *Game
}

func (this GameManager) Run() {
    Info.Println("GameManager: Run")

    this.nextGameId = 1
    this.nextPlayerId = 1

    this.games = make(map[*Game]bool)

    this.gameRunning = make(chan *Game)
    this.gameStopped = make(chan *Game)

    for {
        select {
        // New connection
        case connection := <- this.newConnection:
            Info.Println("New Player", this.nextPlayerId)
            player := &Player{
                ws: connection,
                id: this.nextPlayerId,
                login: fmt.Sprintf("Player #%d", this.nextPlayerId),
            }
            player.ws.send <- GRTSMessage.Connected(player)
            this.findMatchForPlayer(player)
            this.nextPlayerId++

        // A game just started
        case game := <- this.gameRunning:
            this.games[game] = true

        // A game just stopped, delete it
        case game := <- this.gameStopped:
            delete(this.games, game)
        }
    }
}

// New game configuration
func (this *GameManager) createGame(player *Player) {
    Info.Println("Creating a new room")

    // Create game and add it to the game list
    newGame := &Game{
        id: this.nextGameId,
        name: fmt.Sprintf("Game #%d", this.nextGameId),
        numberOfPlayers: 2,
    }
    newGame.newPlayer = make(chan *Player)

    this.games[newGame] = false

    // Launch new game
    go newGame.Run(this.gameRunning, this.gameStopped)

    // Add player
    newGame.newPlayer <- player

    this.nextGameId++;
}

func (this *GameManager) findMatchForPlayer(player *Player) {
    Info.Println("Looking for a room for the new player")

    // Looking for an existing game
    for game, isRunning := range this.games {
        if !isRunning {
            game.newPlayer <- player
            return
        }
    }

    // If no game found, create one
    this.createGame(player)
}
