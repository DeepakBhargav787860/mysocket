package global

import (
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"gorm.io/gorm"
)

var Upgrader = websocket.Upgrader{}

type ChatLoad struct {
	UserProfileId uint `json:"userProfileId" binding:"required"`
	FriendId      uint `json:"friendId" binding:"required"`
}

type ChatWindow struct {
	UserProfileId uint   `json:"userProfileId" binding:"required"`
	Content       string `json:"content" `
	FriendId      uint   `json:"friendId" binding:"required"`
	Type          string `json:"type"`
	AudioData     string `json:"audioData"`
}

type Event struct {
	Type string `json:"type"` // "typing" or "stop_typing"
	From uint   `json:"from"`
	To   uint   `json:"to"`
}

type Message struct {
	gorm.Model
	UserProfileId uint   `json:"userProfileId" binding:"required"`
	Content       string `json:"content"`
	FriendId      uint   `json:"friendId" binding:"required"`
	FilePath      string `json:"filePath" binding:"required"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

type MyBestHalfId struct {
	Id uint `json:"id" binding:"required" `
}

type ReqMobileNo struct {
	MobileNo string `json:"mobileNo" binding:"required" `
}

type RequestAcb struct {
	Status        string `json:"status" binding:"required" `
	UserProfileId uint   `json:"userProfileId" binding:"required"`
	RequestId     uint   `json:"requestId" binding:"required"`
}

type MyBestHalf struct {
	Username string `json:"username" binding:"required" `
}

type RequestSend struct {
	UserProfileId uint `json:"userProfileId" binding:"required"`
	FriendId      uint `json:"friendId" binding:"required"`
}

type UserProfile struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UUID      uuid.UUID `json:"uuid"` // auto handled by Postgres
	Username  string    `json:"username" binding:"required" gorm:"unique;not null"`
	MobileNo  string    `json:"mobileNo" binding:"required" gorm:"unique;not null"`
	Address   string    `json:"address" binding:"required"`
	EmailId   string    `json:"emailId" binding:"required"`
	Password  string    `json:"password" binding:"required"`
	IsLogin   bool      `json:"isLogin" gorm:"default:false"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type BlockedAccount struct {
	ID   uint      `gorm:"primaryKey"`
	UUID uuid.UUID `json:"uuid"`
}

type Request struct {
	ID        uint      `gorm:"primaryKey"`
	UUID      uuid.UUID `json:"uuid"` // auto handled by Postgres
	Username  string    `json:"username" binding:"required" gorm:"unique;not null"`
	MobileNo  string    `json:"mobileNo" binding:"required" gorm:"unique;not null"`
	Address   string    `json:"address" binding:"required"`
	EmailId   string    `json:"emailId" binding:"required"`
	Password  string    `json:"password" binding:"required"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type UserFriend struct {
	UUID            uuid.UUID   `json:"uuid"`
	ID              uint        `gorm:"primaryKey"`
	UserProfileId   uint        `json:"userProfileId" binding:"required"`
	PD              uint        `json:"pd"`
	RequestId       uint        `json:"requestId" binding:"required"`
	FriendReqStatus string      `json:"friendReqStatus" gorm:"default:'NO'"` //YES NO
	UserProfile     UserProfile `json:"userProfileData"`
	Request         Request     `json:"requestData"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type LoginCreds struct {
	Password string `json:"password" binding:"required"`
	MobileNo string `json:"mobileNo" binding:"required"`
}

var (
	DBase       *gorm.DB
	UserMessage *Message
)
