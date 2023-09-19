package controllers

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/aliftoriq/go-crud/repositories"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type BucketControllers interface {
	UploadImageToMinio(c *gin.Context)
	GetImage(c *gin.Context)
	DeleteImage(c *gin.Context)
}

type bucketControllers struct {
	bucketRepository repositories.BucketRepository
}

func NewBucketControllers(buckerRepo repositories.BucketRepository) BucketControllers {
	return &bucketControllers{bucketRepository: buckerRepo}
}

func (bc *bucketControllers) UploadImageToMinio(c *gin.Context) {

	// untuk testing
	// bucketName := "test-minio"
	// objectName := "test.jpg"

	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newUUID := uuid.NewString()
	bucketName := os.Getenv("BUCKETNAME")
	objectName := newUUID + ".jpg"

	fileContent, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer fileContent.Close()

	// _, err = minioClient.bucket.PutObject(context.Background(), bucketName, objectName, fileContent, file.Size, minio.PutObjectOptions{ContentType: "img/png"})

	err = bc.bucketRepository.PutObject(c, bucketName, objectName, fileContent, file.Size)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "Image uploaded successfully",
		"fileName": objectName,
	})
}

func (bc *bucketControllers) GetImage(c *gin.Context) {
	objectName := c.Param("id")
	bucketName := os.Getenv("BUCKETNAME")

	// unit testing
	// bucketName := "test-minio"

	// image, err := minioClient.bucket.GetObject(c, bucketName, objectName, minio.GetObjectOptions{})

	image, err := bc.bucketRepository.GetObject(c, bucketName, objectName)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Connection Refused / File Name Not Found",
			"error":   err.Error()})
		fmt.Println(err)
	}

	// defer image.Close()

	_, err = io.Copy(c.Writer, image)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Image Not Found"})
		return
	}

	c.Header("Content-Type", "image/jpeg")
}

func (bc *bucketControllers) DeleteImage(c *gin.Context) {
	bucketName := os.Getenv("BUCKETNAME")
	objectName := c.Param("id")

	// err := minioClient.bucket.RemoveObject(c, bucketName, objectName, minio.RemoveObjectOptions{})
	err := bc.bucketRepository.DeleteObject(c, bucketName, objectName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Image deleted successfully",
	})
}
