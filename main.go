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

	middlewareAuth := middleware.NewAuth()

	userRepo := repositories.NewUserRepository()
	userController := controllers.NewUsersController(userRepo)

	arRepo := repositories.NewArticleRepository()
	cacheRepo := repositories.NewCacheRepository()
	arController := controllers.NewArticlesController(arRepo, cacheRepo)

	bucketRepo := repositories.NewBucketRepository()
	bucketController := controllers.NewBucketControllers(bucketRepo)

	r.POST("/signup", userController.Signup)
	r.POST("/login", userController.Login)
	r.GET("/validate", middlewareAuth.RequireAuth, userController.Validate)

	r.GET("/users/:id", middlewareAuth.RequireAuth, userController.GetUser)
	r.PUT("/users/:id", middlewareAuth.RequireAuth, userController.UpdateUser)
	r.DELETE("/users/:id", middlewareAuth.RequireAuth, userController.DeleteUser)

	r.POST("/articles", middlewareAuth.RequireAuth, arController.CreateArticle)
	r.PUT("/articles/:id", middlewareAuth.RequireAuth, arController.UpdateArticle)
	r.GET("/articles", middlewareAuth.RequireAuth, arController.GetArticles)
	r.GET("/articles/:id", middlewareAuth.RequireAuth, arController.GetArticleByID)
	r.DELETE("/articles/:id", middlewareAuth.RequireAuth, arController.DeleteArticle)

	r.POST("/upload-image", middlewareAuth.RequireAuth, bucketController.UploadImageToMinio)
	r.GET("/image/:id", middlewareAuth.RequireAuth, bucketController.GetImage)
	r.DELETE("/image/:id", middlewareAuth.RequireAuth, bucketController.DeleteImage)

	r.Run()
}
