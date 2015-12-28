package main

import (
    "bytes"
    "fmt"
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
    players         map[uint32]*Player

    newPlayer           chan *Player
    recv                chan PlayerMessage
    unregisterPlayer    chan uint32
}

func (this *Game) addPlayer(player *Player) {
    this.players[player.id] = player
    this.registerPlayerToGeneralRecv(player, this.recv, this.unregisterPlayer)

    player.ws.send <- []byte("You joined the game " + string(this.id))
}

func (this *Game) removePlayer(id uint32) {
    // Kill recv->generalRecv routine
    this.unregisterPlayer <- id
    delete(this.players, id)
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
        case player := <- this.newPlayer:
            Info.Println("Game", this.id, "has a new player")
            this.addPlayer(player)
        case playerMessage := <- this.recv:
            Info.Println("Receive message from Player", playerMessage.playerId, ":", string(playerMessage.data))
            if len(playerMessage.data) == 0 || bytes.Equal(playerMessage.data, []byte("QUIT GAME")) {
                this.removePlayer(playerMessage.playerId)
            }
        }
    }

    for {
        select{
        case playerMessage := <- this.recv:
            Info.Println("Receive message from Player", playerMessage.playerId, ":", string(playerMessage.data))
            if bytes.Equal(playerMessage.data, []byte("QUIT GAME")) {
                this.removePlayer(playerMessage.playerId)
            }
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
    player.ws.send <- []byte("I'm looking for your game bro")
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
