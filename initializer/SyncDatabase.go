package initializer

import Models "github.com/aliftoriq/go-crud/models"

func SyncDatabase() {
	// DB.Migrator().DropTable(&Models.Article{})
	DB.AutoMigrate(&Models.User{}, &Models.Article{})
}
