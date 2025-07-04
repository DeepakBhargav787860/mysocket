package main

import (
	"bytes"
	businesslogic "chatapp/businessLogic"
	"chatapp/db"
	"chatapp/global"
	"chatapp/route"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

type LogtailWriter struct {
	endpoint string
}

func (w *LogtailWriter) Write(p []byte) (n int, err error) {
	logMsg := []byte(`{"message": "` + string(bytes.TrimSpace(p)) + `"}`)
	resp, err := http.Post(w.endpoint, "application/json", bytes.NewBuffer(logMsg))
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	return len(p), nil
}

func main() {
	// Database connect
	global.DBase = db.DB()
	if global.DBase == nil {
		log.Println("failed to connect database😔😔😔😔😔😔😔")
	} else {
		log.Println("database connected successfully 😘😘😘😘😘😘😘")
	}
	//logs
	logtailEndpoint := "https://s1369505.eu-nbg-2-vec.betterstackdata.com"

	// Logtail + Terminal dono me logs bhejna
	log.SetOutput(io.MultiWriter(os.Stdout, &LogtailWriter{endpoint: logtailEndpoint}))
	// Routes
	route.Routes()
	go businesslogic.HandleMessages()
	fmt.Println("Server started at :8080")
	http.ListenAndServe(":8080", nil)
}
