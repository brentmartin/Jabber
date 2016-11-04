package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

// initialize Hub
type Hub struct {
	connections       map[*Connection]bool
	broadcast         chan []byte
	createConnection  chan *Connection
	destroyConnection chan *Connection
}

// initialize connections
type Connection struct {
	hub  *Hub
	conn *websocket.Conn
	send chan []byte
}

func main() {
	hub := newHub()
	go hub.launch()

	// websocket handler
	http.HandleFunc("/socket", func(w http.ResponseWriter, r *http.Request) {
		socketChat(hub, w, r)
	})

	// index file handler
	http.Handle("/", http.FileServer(http.Dir("./public")))

	// start server
	err := http.ListenAndServe(":"+os.Getenv("PORT"), nil)
	if err != nil {
		panic("Error: " + err.Error())
	}
}

// creates a new hub object when called
func newHub() *Hub {
	return &Hub{
		broadcast:         make(chan []byte),
		createConnection:  make(chan *Connection),
		destroyConnection: make(chan *Connection),
		connections:       make(map[*Connection]bool),
	}
}

// runs a goroutine for Hub when called
func (hub *Hub) launch() {
	for {
		select {
		// receive connections from client and pass to hub to store
		case connection := <-hub.createConnection:
			hub.connections[connection] = true
		// receive disconnects from client and pass to hub to delete
		case connection := <-hub.destroyConnection:
			if _, ok := hub.connections[connection]; ok {
				delete(hub.connections, connection)
				close(connection.send)
			}
		//  receive messages from client and pass to hub to broadcast
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

// write messages to server (goroutine)
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

// read messages from server (goroutine)
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

	// new Connection object
	connection := &Connection{hub: hub, conn: conn, send: make(chan []byte)}
	hub.createConnection <- connection

	go connection.reader()
	go connection.writer()
}
