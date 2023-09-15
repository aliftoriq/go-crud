package initializer

import (
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectToDb() {
	//postgres://vibknctl:SVMnCJb4KJnh4OCSMnuVCr9LTtiiH6Ve@rosie.db.elephantsql.com/vibknctl
	var err error
	dsn := os.Getenv("DB")
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("FAILED TO CONNECT THE DATABASE")
	}
}
