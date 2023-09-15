package initializer

import (
	"fmt"
	"log"
	"os"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var Client *minio.Client

func ConnectToMinio() {
	endpoint := os.Getenv("ENDPOINT")
	accessKey := os.Getenv("ACCESKEY")
	secretKey := os.Getenv("SECRETKEY")

	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: false,
	})
	if err != nil {
		fmt.Println("minio client")
		log.Fatalln(err)
	}

	Client = minioClient
}
