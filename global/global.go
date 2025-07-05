package global

import (
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"gorm.io/gorm"
)

var Upgrader = websocket.Upgrader{}

type Message struct {
	gorm.Model
	UserProfileId uint   `json:"userProfileId" binding:"required"`
	Content       string `json:"content"`
	FriendId      uint   `json:"friendId" binding:"required"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

type MyBestHalfId struct {
	Id uint `json:"id" binding:"required" `
}

type MyBestHalf struct {
	Username string `json:"username" binding:"required" `
}

type RequestSend struct {
	UserProfileId uint `json:"userProfileId" binding:"required"`
	FriendId      uint `json:"friendId" binding:"required"`
}

type UserProfile struct {
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
