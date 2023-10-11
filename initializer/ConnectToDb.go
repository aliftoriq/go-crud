package initializer

import (
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// type Database interface {
// }

// type database struct {
// 	DB *gorm.DB
// }

// func NewDatabase() Database {
// 	db := ConnectToDb()
// 	return &database{DB: db}
// }

var DB *gorm.DB

func ConnectToDb() {
	var err error
	dsn := os.Getenv("DB")
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("FAILED TO CONNECT THE DATABASE")
	}
}
