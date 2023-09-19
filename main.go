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

	arRepo := repositories.NewArticleRepository()
	cacheRepo := repositories.NewCacheRepository()
	arController := controllers.NewArticlesController(arRepo, cacheRepo)

	bucketRepo := repositories.NewBucketRepository()
	bucketController := controllers.NewBucketControllers(bucketRepo)

	r.POST("/signup", userController.Signup)
	r.POST("/login", userController.Login)
	r.GET("/validate", middleware.RequireAuth, userController.Validate)

	r.GET("/users/:id", middleware.RequireAuth, userController.GetUser)
	r.PUT("/users/:id", middleware.RequireAuth, userController.UpdateUser)
	r.DELETE("/users/:id", middleware.RequireAuth, userController.DeleteUser)

	r.POST("/articles", middleware.RequireAuth, arController.CreateArticle)
	r.PUT("/articles/:id", middleware.RequireAuth, arController.UpdateArticle)
	r.GET("/articles", middleware.RequireAuth, arController.GetArticles)
	r.GET("/articles/:id", middleware.RequireAuth, arController.GetArticleByID)
	r.DELETE("/articles/:id", middleware.RequireAuth, arController.DeleteArticle)

	r.POST("/upload-image", middleware.RequireAuth, bucketController.UploadImageToMinio)
	r.GET("/image/:id", middleware.RequireAuth, bucketController.GetImage)
	r.DELETE("/image/:id", middleware.RequireAuth, bucketController.DeleteImage)

	r.Run()
}
