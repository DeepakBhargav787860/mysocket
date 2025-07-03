package main

import "net/http"

func Routes() {
	http.HandleFunc("/ws", handleConnections)
}
