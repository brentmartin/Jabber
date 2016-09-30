package main

import (
	"log"
	"net/http"

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
		case connection := <-hub.createConnection:
			hub.connections[connection] = true
		//      set function to receive disconnects and delete them
		case connection := <-hub.destroyConnection:
			if _, ok := hub.connections[connection]; ok {
				delete(hub.connections, connection)
				close(connection.send)
			}
		//      set function to receive messages from client and broadcast back to all client
		case message := <-hub.broadcast:
			for connection := range hub.connections {
				select {
				case connection.send <- message:
				default:
					close(connection.send)
					delete(hub.connections, connection)
				}
			}
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

// TODO:  Update main func for hub and connections
//          function for new hub
func newHub() *Hub {
	return &Hub{
		broadcast:         make(chan []byte),
		createConnection:  make(chan *Connection),
		destroyConnection: make(chan *Connection),
		connections:       make(map[*Connection]bool),
	}
}

// TODO:  Include an upgrader to upgrade the http connection to a websocket
//          create new connection object each time a connection is upgraded
//          pass new connection to hub to be stored

func main() {
	//          create new Hub
	hub := newHub()
	//          run the hub as a goroutine
	go hub.launch()
	// websocket handler
	//          update socket handler to send connections to hub to store
	http.HandleFunc("/socket", func(w http.ResponseWriter, r *http.Request) {
		socketChat(hub, w, r)
	})
	// index file handler
	http.Handle("/", http.FileServer(http.Dir(".")))
	// start server
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic("Error: " + err.Error())
	}
}

// write messages to server
//          build writer function for new Connections
func (c *Connection) writer() {
	defer func() {
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			c.conn.WriteMessage(websocket.TextMessage, message)
		}
	}
}

// read messages from server
//          build reader function for new Connections
func (c *Connection) reader() {
	defer func() {
		c.hub.destroyConnection <- c
		c.conn.Close()
	}()
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			c.hub.destroyConnection <- c
			c.conn.Close()
			break
		}
		c.hub.broadcast <- message
	}
}

func socketChat(hub *Hub, w http.ResponseWriter, r *http.Request) {
	// upgrade the connection
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	//          create new Connection
	connection := &Connection{hub: hub, conn: conn, send: make(chan []byte)}
	hub.createConnection <- connection

	//          run read and write as goroutines
	go connection.reader()
	go connection.writer()
}
