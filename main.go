package main

import (
	"bytes"
	businesslogic "chatapp/businessLogic"
	"chatapp/cros"
	"chatapp/db"
	"chatapp/global"
	"chatapp/middleware"
	securemiddleware "chatapp/secureMiddleware"
	"log"
	"net/http"
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

	middleware.InitLogger()

	mux := http.NewServeMux()
	mux.HandleFunc("/ws", businesslogic.HandleConnections)
	mux.Handle("/health", cros.EnableCORS(http.HandlerFunc(businesslogic.CheckHealth)))
	//with socket
	mux.HandleFunc("/createProfile", businesslogic.CreateProfile)
	//without socket
	mux.Handle("/signUpUser", cros.EnableCORS(http.HandlerFunc(businesslogic.SignUp)))
	mux.Handle("/loginUser", cros.EnableCORS(http.HandlerFunc(businesslogic.LoginUser)))
	mux.Handle("/userProfile", securemiddleware.AuthMiddleware(http.HandlerFunc(businesslogic.UserProfile)))
	mux.Handle("/requestSend", securemiddleware.AuthMiddleware(http.HandlerFunc(businesslogic.RequestSend)))
	mux.Handle("/getRequestSend", securemiddleware.AuthMiddleware(http.HandlerFunc(businesslogic.GetRequestSend)))
	mux.Handle("/requestCome", securemiddleware.AuthMiddleware(http.HandlerFunc(businesslogic.RequestCome)))
	mux.Handle("/activeUser", securemiddleware.AuthMiddleware(http.HandlerFunc(businesslogic.ActiveUser)))
	mux.Handle("/getAllUser", securemiddleware.AuthMiddleware(http.HandlerFunc(businesslogic.GetAllUser)))
		mux.Handle("/findUserByMobileNo", cros.EnableCORS(securemiddleware.AuthMiddleware(http.HandlerFunc(businesslogic.FindUserByMobileNo))))
	// Wrap your mux with logging middleware
	loggedMux := middleware.LoggingMiddleware(mux)
	go businesslogic.HandleMessages()
	http.ListenAndServe(":8080", loggedMux)
}
