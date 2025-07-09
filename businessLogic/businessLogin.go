package businesslogic

import (
	"chatapp/global"
	passwordhashing "chatapp/passwordHashing"
	"chatapp/response"
	securemiddleware "chatapp/secureMiddleware"

	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var clients = make(map[*websocket.Conn]bool)
var broadcast = make(chan global.Message)

type Resp struct {
	Response string `json:"response"`
}

func HandleConnections(w http.ResponseWriter, r *http.Request) {
	// when cross platform used like frontend and backend work on different platform
	global.Upgrader.CheckOrigin = func(r *http.Request) bool {
		origin := r.Header.Get("Origin")
		return origin == "https://chat-steel-zeta-49.vercel.app" || origin == "https://chat-sj5k-deepaks-projects-a5241927.vercel.app"
	}
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
	global.Upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	// global.Upgrader.CheckOrigin = func(r *http.Request) bool {
	// 	userUUid, err := securemiddleware.GetUserIDFromContext(r)
	// 	if err != nil {
	// 		w.Write([]byte("error in connection stablish"))
	// 		return false
	// 	}

	// 	if userUUid != uuid.Nil {
	// 		log.Println("connection stablish", userUUid)
	// 		return true
	// 	}
	// 	return false

	// }
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
	hashPassword, hashErr := passwordhashing.HashPasswordArgon2(data.Password)

	if hashErr != nil {
		http.Error(w, "failed to password hashing", http.StatusInternalServerError)
		return
	}
	data.Password = hashPassword
	data.UUID = uuid.New()
	if err := global.DBase.Create(&data).Error; err != nil {
		ws.WriteMessage(websocket.TextMessage, []byte("error accoured when create profile"))
		return
	}
	ws.WriteMessage(websocket.TextMessage, []byte("profile successfully created"))

}

func SignUp(w http.ResponseWriter, r *http.Request) {
	var input global.UserProfile
	err := json.NewDecoder(r.Body).Decode(&input)

	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if input.Username == "" || input.MobileNo == "" || input.EmailId == "" || input.Address == "" || input.Password == "" {
		http.Error(w, "required all field", http.StatusBadRequest)
		return
	}

	hashPassword, hashErr := passwordhashing.HashPasswordArgon2(input.Password)

	if hashErr != nil {
		http.Error(w, "failed to password hashing", http.StatusInternalServerError)
		return
	}
	input.Password = hashPassword
	input.UUID = uuid.New()
	if err := global.DBase.Create(&input).Error; err != nil {
		http.Error(w, "failed to create", http.StatusInternalServerError)
		return
	}

	response.MessagePassed(w, "signup successfully")

}

func LoginUser(w http.ResponseWriter, r *http.Request) {
	var input global.LoginCreds
	err := json.NewDecoder(r.Body).Decode(&input)

	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if input.MobileNo == "" || input.Password == "" {
		http.Error(w, "enter the correct details", http.StatusBadRequest)
		return
	}
	var data global.UserProfile

	if err := global.DBase.Model(&global.UserProfile{}).Where("mobile_no=?", input.MobileNo).Find(&data).Error; err != nil {
		http.Error(w, "failed to find in database", http.StatusInternalServerError)
		return
	}

	ok, err := passwordhashing.ComparePasswordArgon2(input.Password, data.Password)
	if err != nil || !ok {
		http.Error(w, "failed to login or incorrect password", http.StatusInternalServerError)
		return
	}

	if ok {
		log.Println("uuid", data.UUID)
		token, err := securemiddleware.GenerateJWT(data.UUID)
		if err != nil {
			http.Error(w, "Failed to generate token", http.StatusInternalServerError)
			return
		}

		if err := global.DBase.Model(&global.UserProfile{}).Where("mobile_no=?", input.MobileNo).Updates(map[string]interface{}{
			"IsLogin": true,
		}).Error; err != nil {
			http.Error(w, "failed to login", http.StatusInternalServerError)
			return
		}

		// Set JWT in secure cookie
		http.SetCookie(w, &http.Cookie{
			Name:     "Authorization",
			Value:    token,
			Path:     "/",
			HttpOnly: true,
			Secure:   true,                  // Only over HTTPS
			SameSite: http.SameSiteNoneMode, // For cross-origin
			MaxAge:   3600,                  // 1 day
		})

		// Send safe user data
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "Login successful",
			"user": map[string]interface{}{
				"uuid":     data.UUID,
				"id":       data.ID,
				"username": data.Username,
				"mobileNo": data.MobileNo,
				"emailId":  data.EmailId,
				"address":  data.Address,
			},
		})
	}

}

func ActiveUser(w http.ResponseWriter, r *http.Request) {

	var input global.MyBestHalfId
	err := json.NewDecoder(r.Body).Decode(&input)

	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if input.Id == 0 {
		http.Error(w, "something went wrong", http.StatusBadRequest)
		return
	}
	var data []global.UserFriend

	if err := global.DBase.Model(&global.UserFriend{}).
		Where("(user_profile_id = ? OR pd = ?) AND friend_req_status IN(?)", input.Id, input.Id, []string{"YES", "ACCEPTED"}).
		Preload("Request").
		Preload("UserProfile").
		Find(&data).Error; err != nil {
		http.Error(w, "failed to find in database", http.StatusInternalServerError)
		return
	}

	response.MessagePassed(w, data)
}

func Logout(w http.ResponseWriter, r *http.Request) {
	var logoutId global.MyBestHalfId
	err := json.NewDecoder(r.Body).Decode(&logoutId)
	if err != nil {
		http.Error(w, "something went wrong", http.StatusBadRequest)
		return
	}

	if err := global.DBase.Model(&global.UserProfile{}).Where("id=?", logoutId.Id).Updates(map[string]interface{}{
		"IsLogin": false,
	}).Error; err != nil {
		http.Error(w, "failed to logout", http.StatusInternalServerError)
		return
	}
	//reset cookie
	// also expire token which is set 1 hour (pending task)
	http.SetCookie(w, &http.Cookie{
		Name:     "Authorization",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})

	// Optional response
	response.MessagePassed(w, "logout successfully")
}

func UserProfile(w http.ResponseWriter, r *http.Request) {

	var input global.MyBestHalf
	err := json.NewDecoder(r.Body).Decode(&input)

	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if input.Username == "" {
		http.Error(w, "something went wrong", http.StatusBadRequest)
		return
	}
	var data global.UserProfile

	if err := global.DBase.Model(&global.UserProfile{}).Where("username=?", input.Username).Find(&data).Error; err != nil {
		http.Error(w, "failed to find in database", http.StatusInternalServerError)
		return
	}

	response.MessagePassed(w, data)
}

func GetRequestSend(w http.ResponseWriter, r *http.Request) {

	var input global.MyBestHalfId
	err := json.NewDecoder(r.Body).Decode(&input)

	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if input.Id == 0 {
		http.Error(w, "something went wrong", http.StatusBadRequest)
		return
	}
	var data []global.UserFriend

	if err := global.DBase.Model(&global.UserFriend{}).Where("user_profile_id=?  AND friend_req_status NOT IN (?)", input.Id, []string{"BLOCKED"}).Preload("Request").Find(&data).Error; err != nil {
		http.Error(w, "failed to find in database", http.StatusInternalServerError)
		return
	}

	response.MessagePassed(w, data)
}

func RequestSend(w http.ResponseWriter, r *http.Request) {
	var input global.RequestSend
	err := json.NewDecoder(r.Body).Decode(&input)

	if err != nil {

		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if input.UserProfileId == 0 || input.FriendId == 0 {
		http.Error(w, "something went wrong", http.StatusBadRequest)
		return
	}

	// get friend data
	var fData global.UserProfile
	if err := global.DBase.Model(&global.UserProfile{}).Where("id=?", input.FriendId).Find(&fData).Error; err != nil {
		http.Error(w, "failed to find in database", http.StatusInternalServerError)
		return
	}

	// create request data

	var rData = global.Request{
		UUID:     uuid.New(),
		Username: fData.Username,
		MobileNo: fData.MobileNo,
		Address:  fData.Address,
		EmailId:  fData.EmailId,
		Password: fData.Password,
	}

	//find already exist

	var alreadyExist global.Request
	if err := global.DBase.Model(&global.Request{}).Where("username=?", fData.Username).Find(&alreadyExist).Error; err != nil {
		http.Error(w, "failed to find in database", http.StatusInternalServerError)
		return
	}

	if alreadyExist.Username == "" {
		if err := global.DBase.Create(&rData).Error; err != nil {
			http.Error(w, "failed to find in database", http.StatusInternalServerError)
			return
		}
	}

	// get request data id
	var rId global.Request
	if err := global.DBase.Model(&global.Request{}).Where("username=?", fData.Username).Find(&rId).Error; err != nil {
		http.Error(w, "failed to find in database", http.StatusInternalServerError)
		return
	}

	// create request

	var rCreate = global.UserFriend{
		UUID:          uuid.New(),
		UserProfileId: input.UserProfileId,
		RequestId:     rId.ID,
		PD:            input.FriendId,
	}

	var findAlreadySendReq global.UserFriend

	if err := global.DBase.Model(&global.UserFriend{}).Where("user_profile_id=? AND pd=?", rCreate.UserProfileId, rCreate.RequestId).Find(&findAlreadySendReq).Error; err != nil {
		http.Error(w, "failed to find in database", http.StatusInternalServerError)
		return
	}

	if findAlreadySendReq.RequestId != 0 {
		http.Error(w, "request already send", http.StatusInternalServerError)
		return
	}

	if err := global.DBase.Debug().Model(&global.UserFriend{}).Create(&rCreate).Error; err != nil {
		http.Error(w, "failed to send request", http.StatusInternalServerError)
		return
	}

	response.MessagePassed(w, "request successfully send")
}

func RequestCome(w http.ResponseWriter, r *http.Request) {

	var input global.MyBestHalfId
	err := json.NewDecoder(r.Body).Decode(&input)

	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if input.Id == 0 {
		http.Error(w, "something went wrong", http.StatusBadRequest)
		return
	}
	var data []global.UserFriend

	if err := global.DBase.Model(&global.UserFriend{}).Where("request_id=? AND friend_req_status=?", input.Id, "NO").Preload("UserProfile").Find(&data).Error; err != nil {
		http.Error(w, "failed to find in database", http.StatusInternalServerError)
		return
	}

	response.MessagePassed(w, data)
}

func GetAllUser(w http.ResponseWriter, r *http.Request) {
	var input global.MyBestHalfId
	err := json.NewDecoder(r.Body).Decode(&input)

	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if input.Id == 0 {
		http.Error(w, "something went wrong", http.StatusBadRequest)
		return
	}
	var data []global.UserProfile
	if err := global.DBase.Model(&global.UserProfile{}).Where("id != ? ", input.Id).Find(&data).Error; err != nil {
		http.Error(w, "failed to find in database", http.StatusInternalServerError)
		return
	}

	response.MessagePassed(w, data)
}

func FindUserByMobileNo(w http.ResponseWriter, r *http.Request) {

	var input global.ReqMobileNo
	err := json.NewDecoder(r.Body).Decode(&input)

	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if input.MobileNo == "" {
		http.Error(w, "something went wrong", http.StatusBadRequest)
		return
	}
	var data global.UserProfile

	if err := global.DBase.Model(&global.UserProfile{}).Where("mobile_no=?", input.MobileNo).Find(&data).Error; err != nil {
		http.Error(w, "failed to find in database", http.StatusInternalServerError)
		return
	}

	if data.ID != 0 {
		response.MessagePassed(w, data)
	} else {
		http.Error(w, "no user found", http.StatusNotFound)
	}

}

func RequestARB(w http.ResponseWriter, r *http.Request) {

	var input global.RequestAcb
	err := json.NewDecoder(r.Body).Decode(&input)

	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if input.Status == "" || input.UserProfileId == 0 || input.RequestId == 0 {
		http.Error(w, "something went wrong", http.StatusBadRequest)
		return
	}

	if err := global.DBase.Model(&global.UserFriend{}).Where("user_profile_id=? AND request_id=?", input.UserProfileId, input.RequestId).Updates(map[string]interface{}{
		"FriendReqStatus": input.Status,
	}).Error; err != nil {
		http.Error(w, "failed to find in database", http.StatusInternalServerError)
		return
	}
	response.MessagePassed(w, fmt.Sprintf("Request Successfully: %s", input.Status))
}

func InComingRequest(w http.ResponseWriter, r *http.Request) {
	//open ws connection
	global.Upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	// global.Upgrader.CheckOrigin = func(r *http.Request) bool {
	// 	userUuid, err := securemiddleware.GetUserIDFromContext(r)
	// 	if err != nil {
	// 		http.Error(w, "failed to origin check", http.StatusBadRequest)
	// 		log.Println("error in open connection", err)
	// 		return false
	// 	}
	// 	if userUuid != uuid.Nil {
	// 		log.Println("connection open")
	// 		return true
	// 	}
	// 	http.Error(w, "failed to origin check1", http.StatusBadRequest)
	// 	log.Println("connection failed")
	// 	return false

	// }

	ws, err := global.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "Could not open websocket connection", http.StatusBadRequest)
		// ws.WriteMessage(websocket.TextMessage, []byte("failed to websocket setup"))
		fmt.Println(err)
		return
	}

	defer ws.Close()
	// var getInReq = make(chan []global.UserFriend)
	// read request data
	for {
		var input global.MyBestHalfId
		err := ws.ReadJSON(&input)
		if err != nil {
			ws.WriteMessage(websocket.TextMessage, []byte("failed to read data"))
			//always break not return
			break
		}
		go func() {
			var data []global.UserFriend
			if err := global.DBase.Model(&global.UserFriend{}).Where("pd=?", input.Id).Preload("UserProfile").Find(&data).Error; err != nil {
				ws.WriteMessage(websocket.TextMessage, []byte("failed to find in database"))
				return
			}
			// log.Println("data", data)
			ws.WriteJSON(data)
		}()
	}

	// go func() {
	// 	for data := range getInReq {
	// 		ws.WriteJSON(data)
	// 	}
	// }()

}

// func HandleConnections(w http.ResponseWriter, r *http.Request) {
// 	// when cross platform used like frontend and backend work on different platform
// 	global.Upgrader.CheckOrigin = func(r *http.Request) bool {
// 		origin := r.Header.Get("Origin")
// 		return origin == "https://chat-steel-zeta-49.vercel.app" || origin == "https://chat-sj5k-deepaks-projects-a5241927.vercel.app"
// 	}
// 	ws, err := global.Upgrader.Upgrade(w, r, nil)
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}
// 	defer ws.Close()

// 	clients[ws] = true

// 	for {
// 		var msg global.Message

// 		err := ws.ReadJSON(&msg)
// 		if err != nil {
// 			delete(clients, ws)
// 			break
// 		}
// 		global.DBase.Create(&msg)
// 		broadcast <- msg
// 	}
// }
