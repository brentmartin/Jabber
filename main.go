package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

// TODO:  Include an upgrader to upgrade the http connection to a websocket
//          create new connection object each time a connection is upgraded
//          pass new connection to hub to be stored
// TODO:  Build Hub to store connections and broadcast messages to them
//          initialize Hub object
//          create new Hub
//          build a function for Hub object to work in
//          set function to receive connections from clients and store them
//          set function to receive messages from client and broadcast back to all client
//          run the hub as a goroutine
// TODO:  Build Connection to store websocket connection and send/receive messages
//          initialize Connection object
//          create new Connection
//          build reader function for new Connections
//          build writer function for new Connections
//          run read and write as goroutines
//          update socket handler to send connections to hub to store
// TODO:  Update main func for hub and connections

var upgrader = websocket.Upgrader{}

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
