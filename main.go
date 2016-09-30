package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

// TODO:  Build Hub to store connections and broadcast messages to them
//          initialize Hub object
type Hub struct {
	connections       map[*Connection]bool
	broadcast         chan []byte
	createConnection  chan *Connection
	destroyConnection chan *Connection
}

//          build a function for Hub object to work in
func (hub *Hub) launch() {
	for {
		select {
		//      set function to receive connections from clients and store them
		case conn := <-hub.createConnection:
			hub.connections[conn] = true
		//      set function to receive disconnects and delete them
		case conn := <-hub.destroyConnection:
			if _, ok := hub.connections[conn]; ok {
				delete(hub.connections, conn)
				close(conn.send)
			}
		//      set function to receive messages from client and broadcast back to all client
		}
	}
}

// TODO:  Build Connection to store websocket connection and send/receive messages
//          initialize Connection object
type Connection struct {
	hub  *Hub
	ws   *websocket.Conn
	send chan []byte
}

//          create new Connection
//          build reader function for new Connections
//          build writer function for new Connections
//          run read and write as goroutines
// TODO:  Update main func for hub and connections
//          create new Hub
func newHub() *Hub {
	return &Hub{
		broadcast:         make(chan []byte),
		createConnection:  make(chan *Connection),
		destroyConnection: make(chan *Connection),
		connections:       make(map[*Connection]bool),
	}
}

//          run the hub as a goroutine
//          update socket handler to send connections to hub to store
// TODO:  Include an upgrader to upgrade the http connection to a websocket
//          create new connection object each time a connection is upgraded
//          pass new connection to hub to be stored

func main() {
	// websocket handler
	http.HandleFunc("/socket", socketHandler)
	// index file handler
	http.Handle("/", http.FileServer(http.Dir(".")))
	// start server
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic("Error: " + err.Error())
	}
}

func socketHandler(w http.ResponseWriter, r *http.Request) {
	// upgrade the connection
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	for {
		// read messages from server
		message, p, err := conn.ReadMessage()
		if err != nil {
			return
		}

		// timeout between actions
		time.Sleep(time.Second * 1)

		// write messages to server
		err = conn.WriteMessage(message, p)
		if err != nil {
			return
		}
	}
}
