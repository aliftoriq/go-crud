package main

import (
	"github.com/aliftoriq/go-crud/controllers"
	"github.com/aliftoriq/go-crud/initializer"
	"github.com/aliftoriq/go-crud/middleware"
	"github.com/aliftoriq/go-crud/repositories"
	"github.com/gin-gonic/gin"
)

func init() {
	initializer.LoadEnvVariables()
	initializer.ConnectToDb()
	initializer.SyncDatabase()
	initializer.ConnectToMinio()
	initializer.ConnectToRedis()

}

func main() {
	r := gin.Default()

	userRepo := repositories.NewUserRepository()
	userController := controllers.NewUsersController(userRepo)

	r.POST("/signup", userController.Signup)
	r.POST("/login", userController.Login)
	r.GET("/validate", middleware.RequireAuth, userController.Validate)

	r.GET("/users/:id", middleware.RequireAuth, userController.GetUser)
	r.PUT("/users/:id", middleware.RequireAuth, userController.UpdateUser)
	r.DELETE("/users/:id", middleware.RequireAuth, userController.DeleteUser)

	r.POST("/articles", middleware.RequireAuth, controllers.CreateArticle)
	r.PUT("/articles/:id", middleware.RequireAuth, controllers.UpdateArticle)
	r.GET("/articles", middleware.RequireAuth, controllers.GetArticles)
	r.GET("/articles/:id", middleware.RequireAuth, controllers.GetArticleByID)
	r.DELETE("/articles/:id", middleware.RequireAuth, controllers.DeleteArticle)

	r.POST("/upload-image", middleware.RequireAuth, controllers.UploadImageToMinio)
	r.GET("/image/:id", middleware.RequireAuth, controllers.GetImage)
	r.DELETE("/image/:id", middleware.RequireAuth, controllers.DeleteImage)

	r.Run()
}
