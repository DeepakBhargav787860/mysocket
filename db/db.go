package db

import (
	"chatapp/global"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func DB() *gorm.DB {

	dsn := "host=dpg-d1hnjgvfte5s73ag8n2g-a user=apilab password=PHDZ0Y8yyr5IxaCZy1OLUXL8MoJJkM6W dbname=socket_fjdf port=5432 sslmode=disable"
	// dsn := "host=localhost user=postgres password=deepak123 dbname=socket port=5432 sslmode=disable"
	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("DB connection failed")
	}
	db.AutoMigrate(&global.Message{}, &global.UserProfile{})

	return db
}
