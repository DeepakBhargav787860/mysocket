package route

import (
	businesslogic "chatapp/businessLogic"
	"chatapp/cros"
	"net/http"
)

func Routes() {

	http.HandleFunc("/ws", businesslogic.HandleConnections)
	//for ping (for spin back up)
	http.Handle("/health", cros.EnableCORS(http.HandlerFunc(businesslogic.CheckHealth)))
	http.HandleFunc("/createProfile", businesslogic.CreateProfile)
}
