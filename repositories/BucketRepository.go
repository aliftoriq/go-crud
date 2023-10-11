package repositories

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/aliftoriq/go-crud/initializer"
	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
)

//go:generate mockery --outpkg mocks --name BucketRepository
type BucketRepository interface {
	PutObject(c *gin.Context, bName string, oName string, file io.Reader, fileSize int64) error
	GetObject(c *gin.Context, bName string, oName string) (io.Reader, error)
	DeleteObject(c *gin.Context, bName string, oName string) error
}

type bucketRepository struct {
	minio *minio.Client
}

func NewBucketRepository() BucketRepository {
	return &bucketRepository{minio: initializer.Client}
}

func (br *bucketRepository) PutObject(c *gin.Context, bName string, oName string, file io.Reader, fileSize int64) error {
	_, err := br.minio.PutObject(context.Background(), bName, oName, file, fileSize, minio.PutObjectOptions{ContentType: "img/png"})
	if err != nil {
		return err
	}
	return nil
}

func (br *bucketRepository) GetObject(c *gin.Context, bName string, oName string) (io.Reader, error) {
	image, err := br.minio.GetObject(c, bName, oName, minio.GetObjectOptions{})

	if err != nil {
		fmt.Println(err)
		return nil, errors.New("FAILED TO GET OBJECT")
	}
	return image, nil
}

func (br *bucketRepository) DeleteObject(c *gin.Context, bName string, oName string) error {
	err := br.minio.RemoveObject(c, bName, oName, minio.RemoveObjectOptions{})
	if err != nil {
		return errors.New("FAILED TO DELETE OBJECT")
	}
	return nil
}
