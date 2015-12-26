package main

import (
)

/* == Player ==*/
type Player struct {
    ws *Connection
}

/* == Game ==*/
type Game struct {
    // Game configuration
    numberOfPlayers int

    // Game infos
    id          uint
    players     []*Player

    newPlayer   chan *Player
}

func (this *Game) addPlayer(player *Player) {
    // Insert new game into games array
    players := make([]*Player, len(this.players) + 1)
    copy(players, this.players[:])
    players[len(this.players)] = player
    this.players = players

    player.ws.send <- []byte("You joined the game " + string(this.id))
}

func (this *Game) Run(gameRunning chan<- *Game, gameStopped chan<- *Game) {
    // Waiting for player
    Info.Println("Game", this.id, "running")
    for {
        if len(this.players) == this.numberOfPlayers {
            gameRunning <- this
            break;
        }
        Info.Println("Game", this.id, "waiting for more players")
        player := <- this.newPlayer
        Info.Println("Game", this.id, "has a new player")
        this.addPlayer(player)
    }

    for {

    }
}

/* == Game Manager== */
type GameManager struct {
    newConnection chan *Connection

    // Map of games (bool: inGame?)
    games       map[*Game]bool
    nextGameid  uint

    gameRunning chan *Game
    gameStopped chan *Game
}

func (this GameManager) Run() {
    Info.Println("GameManager: Run")
    this.nextGameid = 1
    this.games = make(map[*Game]bool)

    this.gameRunning = make(chan *Game)
    this.gameStopped = make(chan *Game)


    for {
        select {
        // New connection
        case connection := <- this.newConnection:
            Info.Println("New Player ;)")
            player := new(Player)
            player.ws = connection
            this.findMatchForPlayer(player)

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
    newGame := &Game{id: this.nextGameid, numberOfPlayers: 2}
    newGame.newPlayer = make(chan *Player)

    this.games[newGame] = false

    // Launch new game
    go newGame.Run(this.gameRunning, this.gameStopped)

    // Add player
    newGame.newPlayer <- player

    this.nextGameid++;
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
