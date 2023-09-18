package controllers

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/aliftoriq/go-crud/initializer"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
)

func UploadImageToMinio(c *gin.Context) {
	bucketName := os.Getenv("BUCKETNAME")

	newUUID := uuid.NewString()

	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	objectName := newUUID + ".jpg"

	fileContent, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer fileContent.Close()

	_, err = initializer.Client.PutObject(context.Background(), bucketName, objectName, fileContent, file.Size, minio.PutObjectOptions{ContentType: "img/png"})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "Image uploaded successfully",
		"fileName": objectName,
	})
}

func GetImage(c *gin.Context) {
	bucketName := os.Getenv("BUCKETNAME")
	objectName := c.Param("id")

	image, err := initializer.Client.GetObject(c, bucketName, objectName, minio.GetObjectOptions{})

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Connection Refused / File Name Not Found",
			"error":   err.Error()})
		fmt.Println(err)
	}

	defer image.Close()

	_, err = io.Copy(c.Writer, image)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Image Not Found"})
		return
	}

	c.Header("Content-Type", "image/jpeg")
}

func DeleteImage(c *gin.Context) {
	bucketName := os.Getenv("BUCKETNAME")
	objectName := c.Param("id")

	err := initializer.Client.RemoveObject(c, bucketName, objectName, minio.RemoveObjectOptions{})
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
