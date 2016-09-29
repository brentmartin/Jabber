package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

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
