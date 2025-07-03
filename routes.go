package main

import "net/http"

func Routes() {
	http.HandleFunc("/ws", handleConnections)
	//for ping (for spin back up)
	http.HandleFunc("/health", CheckHealth)
}
