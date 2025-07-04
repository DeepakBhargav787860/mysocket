package businesslogic

import (
	"chatapp/global"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var clients = make(map[*websocket.Conn]bool)
var broadcast = make(chan global.Message)

type Resp struct {
	Response string `json:"response"`
}

func HandleConnections(w http.ResponseWriter, r *http.Request) {
	global.Upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	ws, err := global.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer ws.Close()

	clients[ws] = true

	for {
		var msg global.Message

		err := ws.ReadJSON(&msg)
		if err != nil {
			delete(clients, ws)
			break
		}
		global.DBase.Create(&msg)
		broadcast <- msg
	}
}

func HandleMessages() {
	for {
		msg := <-broadcast
		for client := range clients {
			client.WriteJSON(msg)
		}
	}

}

// func CheckHealth(w http.ResponseWriter, r *http.Request) {

//		log.Println("information", r.Host, r.Method, r.RemoteAddr)
//		// if r.URL.Path == "/health" {
//		// 	log.Println("hey")
//		// }
//		// ws, _ := global.Upgrader.Upgrade(w, r, nil)
//		// defer ws.Close()
//		// ws.WriteMessage(websocket.TextMessage, []byte("profile successfully created"))
//		// w.Header().Set("Access-Control-Allow-Origin", "*")
//		w.Header().Set("Content-Type", "application/json") // important
//		w.WriteHeader(http.StatusOK)
//		res := Resp{Response: "anjali i love uuu tum best ho yrr mere liye best darling 😘😘😘😘😘😘😘😘😘😘😘😘😘😘😘😘😘😘😘😘😘😘😘😘😘😘😘😘😘😘😘😘😘😘😘😘😘😘😘😘😘😘😘"}
//		err := json.NewEncoder(w).Encode(res)
//		if err != nil {
//			log.Println("Error encoding response:", err)
//		}
//	}
func CheckHealth(w http.ResponseWriter, r *http.Request) {
	clientIP := r.Header.Get("X-Forwarded-For")
	if clientIP == "" {
		clientIP = r.RemoteAddr
	}

	log.Println("Client IP:", clientIP)
	log.Println("Method:", r.Method)
	log.Println("RemoteAddr (raw):", r.RemoteAddr)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	res := Resp{Response: "Hello from server!"}
	json.NewEncoder(w).Encode(res)
}

func CreateProfile(w http.ResponseWriter, r *http.Request) {
	ws, err := global.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("error come in create profile", err)
		panic(err)
	}
	defer ws.Close()
	var data global.UserProfile
	readErr := ws.ReadJSON(&data)

	if readErr != nil {
		ws.WriteMessage(websocket.TextMessage, []byte("error when msg read"))
		return
	}
	if err := global.DBase.Create(&data).Error; err != nil {
		ws.WriteMessage(websocket.TextMessage, []byte("error accoured when create profile"))
		return
	}
	ws.WriteMessage(websocket.TextMessage, []byte("profile successfully created"))

}
