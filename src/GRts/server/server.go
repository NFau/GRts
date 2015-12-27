package main

import (
    "fmt"
    "io/ioutil"
    "net/http"
    "github.com/gorilla/websocket"
)

/* -- Server -- */
type Server struct {
    gameManager GameManager
    wsUpgrader websocket.Upgrader
}

func (this Server) loadPage(filename string) ([]byte, error) {
    body, err := ioutil.ReadFile(filename)
    if err != nil {
        return nil, err;
    }
    return body, nil
}

func (this Server) viewHandler(writer http.ResponseWriter, request *http.Request) {
    Info.Println("Receive request : " + request.URL.Path)
    resource := request.URL.Path
    resource = "./client/html" + resource + ".html"
    page, err := this.loadPage(resource)
    if err != nil {
        Error.Println("Cannot load page :" + resource)
        http.Redirect(writer, request, "/index", http.StatusFound)
        return
    }
    fmt.Fprintf(writer, "%s", page)
}

func (this Server) socketConnectionHandler(writer http.ResponseWriter, request *http.Request) {
    Info.Println("Receive socket connection request")
    ws, err := this.wsUpgrader.Upgrade(writer, request, nil)
    if err != nil {
        Error.Println(err)
        return
    }

    Info.Println("New connection !")

    // Create Connection to send it to the gameManager
    connection := new(Connection)
    connection.ws = ws
    connection.uuid, _ = newUUID()
    connection.send = make(chan []byte, 256)
    connection.recv = make(chan []byte, 256)

    go connection.WritePump()
    go connection.ReadPump()

    this.gameManager.newConnection <- connection
}

func (this *Server) Start() {
    this.wsUpgrader = websocket.Upgrader{
        ReadBufferSize:  1024,
        WriteBufferSize: 1024,
    }

    // Create channel to receive new Connection
    this.gameManager.newConnection = make(chan *Connection)

    // Launch the gameManager goroutine
    go this.gameManager.Run()

    http.Handle("/resources/", http.StripPrefix("/resources/", http.FileServer(http.Dir("client/"))))
    http.HandleFunc("/socket", this.socketConnectionHandler)
    http.HandleFunc("/", this.viewHandler)

    Info.Println("Listening on port 8080")
    http.ListenAndServe(":8080", nil)
}
