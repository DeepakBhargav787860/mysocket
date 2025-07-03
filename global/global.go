package global

import (
	"github.com/gorilla/websocket"
	"gorm.io/gorm"
)

var Upgrader = websocket.Upgrader{}

type Message struct {
	gorm.Model
	Username string `json:"username"`
	Content  string `json:"content"`
}

type UserProfile struct {
	gorm.Model
	Username string `json:"username" binding:"required" gorm:"unique;not null"`
	MobileNo string `json:"mobileNo" binding:"required"`
	Address  string `json:"address" binding:"required"`
	EmailId  string `json:"emailId" binding:"required"`
}

var (
	DBase       *gorm.DB
	UserMessage *Message
)
