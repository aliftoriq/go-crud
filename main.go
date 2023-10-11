package main

import (
	_ "github.com/aliftoriq/go-crud/docs"

	"github.com/aliftoriq/go-crud/controllers"
	"github.com/aliftoriq/go-crud/initializer"
	"github.com/aliftoriq/go-crud/middleware"
	"github.com/aliftoriq/go-crud/repositories"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func init() {
	initializer.LoadEnvVariables()
	initializer.ConnectToDb()
	initializer.SyncDatabase()
	initializer.ConnectToMinio()
	initializer.ConnectToRedis()

}

// @title Tag Go Crud Service API
// @version 1.0
// @description A golang Restfull API

// @host localhost:4001

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

	// add swagger
	r.GET("/docs/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

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
