package main

import (
	"fmt"
	"net/http"
)

func main() {
	// Database connect
	DB()
	// Routes
	Routes()
	go handleMessages()
	fmt.Println("Server started at :8080")
	http.ListenAndServe(":8080", nil)
}
