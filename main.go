package main

import (
	"bytes"
	businesslogic "chatapp/businessLogic"
	"chatapp/cros"
	"chatapp/db"
	"chatapp/global"
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

	// middleware.InitLogger()

	// mux := http.NewServeMux()
	http.HandleFunc("/ws", businesslogic.HandleConnections)
	http.Handle("/health", cros.EnableCORS(http.HandlerFunc(businesslogic.CheckHealth)))
	//with socket
	http.HandleFunc("/createProfile", businesslogic.CreateProfile)
	//without socket
	http.Handle("/signUpUser", cros.EnableCORS(http.HandlerFunc(businesslogic.SignUp)))
	http.Handle("/loginUser", cros.EnableCORS(http.HandlerFunc(businesslogic.LoginUser)))
	http.Handle("/logout", cros.EnableCORS(securemiddleware.AuthMiddleware(http.HandlerFunc(businesslogic.Logout))))
	http.Handle("/userProfile", securemiddleware.AuthMiddleware(http.HandlerFunc(businesslogic.UserProfile)))
	http.Handle("/requestSend", cros.EnableCORS(securemiddleware.AuthMiddleware(http.HandlerFunc(businesslogic.RequestSend))))
	http.Handle("/getRequestSend", cros.EnableCORS(securemiddleware.AuthMiddleware(http.HandlerFunc(businesslogic.GetRequestSend))))
	//not use
	http.Handle("/requestCome", cros.EnableCORS(securemiddleware.AuthMiddleware(http.HandlerFunc(businesslogic.RequestCome))))
	//batter use
	http.HandleFunc("/getIncomingRequest", businesslogic.InComingRequest)

	http.Handle("/activeUser", securemiddleware.AuthMiddleware(http.HandlerFunc(businesslogic.ActiveUser)))
	http.Handle("/getAllUser", securemiddleware.AuthMiddleware(http.HandlerFunc(businesslogic.GetAllUser)))
	http.Handle("/findUserByMobileNo", cros.EnableCORS(securemiddleware.AuthMiddleware(http.HandlerFunc(businesslogic.FindUserByMobileNo))))
	http.Handle("/requestARB", cros.EnableCORS(securemiddleware.AuthMiddleware(http.HandlerFunc(businesslogic.RequestARB))))

	// Wrap your mux with logging middleware
	// loggedMux := middleware.LoggingMiddleware(mux)
	go businesslogic.HandleMessages()
	http.ListenAndServe(":8080", nil)
}
