package main

// Code from the chat exemple of gorilla/websocket

import (
    "github.com/gorilla/websocket"
    "time"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)


/* -- Websocket wrapper -- */
type Connection struct {
    ws *websocket.Conn

    uuid string

    // Buffered channel of outgoing message
    send chan []byte
    // Channel to forward client cmd
    recv chan []byte
}

// readPump pumps messages from the websocket connection to the hub.
func (this *Connection) ReadPump() {
	defer func() {
        Warning.Println(this.uuid + ": Something wrong happened")
        this.recv <- []byte("")
		this.ws.Close()
	}()
	this.ws.SetReadLimit(maxMessageSize)
	this.ws.SetReadDeadline(time.Now().Add(pongWait))
	this.ws.SetPongHandler(func(string) error {
        this.ws.SetReadDeadline(time.Now().Add(pongWait));
        return nil
    })
	for {
		_, message, err := this.ws.ReadMessage()
		if err != nil {
			break
		}
		this.recv <- message
	}
}

// write writes a message with the given message type and payload.
func (this *Connection) write(mt int, payload []byte) error {
	this.ws.SetWriteDeadline(time.Now().Add(writeWait))
	return this.ws.WriteMessage(mt, payload)
}

// writePump pumps messages from the MatchManager to the websocket connection.
func (this *Connection) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		this.ws.Close()
	}()
	for {
		select {
		case message, ok := <-this.send:
			if !ok {
				this.write(websocket.CloseMessage, []byte{})
				return
			}
			if err := this.write(websocket.BinaryMessage, message); err != nil {
				return
			}
		case <-ticker.C:
			if err := this.write(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}
