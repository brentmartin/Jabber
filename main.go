package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

// Initialize a variable that uses a websocket upgrader
var upgrader = websocket.Upgrader{}

// Startup main goroutine
func main() {
	// websocket handler
	http.HandleFunc("/socket", socketHandler)
	// html handler
	// start server
}

// Create handler for websocket
func socketHandler(w http.ResponseWriter, r *http.Request) {
	// upgrade the connection
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	for {
		// put in a reader for messages sent to server
		message, p, err := conn.ReadMessage()
		if err != nil {
			return
		}
		// put in a writer for messages printed back from server
		err = conn.WriteMessage(message, p)
		if err != nil {
			return
		}
	}
}

// Create a handler for the html
//      write HTML for messages
//      write JS to send messages
//      stylize a bit
// Start server to serve those handlers
