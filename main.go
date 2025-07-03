package main

import (
	businesslogic "chatapp/businessLogic"
	"chatapp/db"
	"chatapp/global"
	"chatapp/route"
	"fmt"
	"log"
	"net/http"
)

func main() {
	// Database connect
	global.DBase = db.DB()
	if global.DBase == nil {
		log.Println("failed to connect database😔😔😔😔😔😔😔")
	} else {
		log.Println("database connected successfully 😘😘😘😘😘😘😘")
	}

	// Routes
	route.Routes()
	go businesslogic.HandleMessages()
	fmt.Println("Server started at :8080")
	http.ListenAndServe(":8080", nil)
}
