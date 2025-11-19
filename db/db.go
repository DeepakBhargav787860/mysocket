package db

import (
	"chatapp/global"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func DB() *gorm.DB {
	dsn := "host=dpg-d4erkmeuk2gs739nngng-a user=socket_fjdf_21qj_user password=A42GGSzypjV11VeHxRnceCYCJRUGeRVM dbname=socket_fjdf_21qj port=5432 sslmode=disable"
	// dsn := "host=localhost user=postgres password=deepak123 dbname=socket port=5432 sslmode=disable"
	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("DB connection is failed")
	}
	db.AutoMigrate(&global.Message{}, &global.UserProfile{}, &global.UserFriend{}, &global.Request{})

	return db
}
