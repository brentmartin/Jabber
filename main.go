package main

// Initialize a variable that uses a websocket upgrader
// Startup main goroutine
// Create handler for websocket
//      upgrade the connection
//      put in a reader for messages sent to server
//      put in a writer for messages printed back from server
// Create a handler for the html
//      write HTML for messages
//      write JS to send messages
//      stylize a bit
// Start server to serve those handlers
