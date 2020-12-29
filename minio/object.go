package minio

import (
	"context"
	"mime/multipart"
	"net/url"
	"time"

	"github.com/minio/minio-go/v7"
)

// Object file struct
type Object struct {
	ContentType  string    `json:"contentType"`
	Name         string    `json:"name"`
	Size         int64     `json:"size"`
	LastModified time.Time `json:"lastModified"`
}

// ListObjects lists objects in a bucket.
func ListObjects(bucketName, prefix string) []Object {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if len(prefix) > 0 && prefix[0:1] == "/" {
		prefix = prefix[1:]
	}
	opts := minio.ListObjectsOptions{
		Prefix:    prefix,
		Recursive: false,
	}

	var objectList []Object
	for object := range client.ListObjects(ctx, bucketName, opts) {
		if object.Err == nil {
			var objectInfo minio.ObjectInfo
			if object.Size > 0 {
				objectInfo, _ = client.StatObject(ctx, bucketName, object.Key, minio.StatObjectOptions{})
			}

			objectList = append(objectList, Object{
				ContentType:  objectInfo.ContentType,
				Name:         object.Key,
				Size:         object.Size,
				LastModified: object.LastModified,
			})
		}
	}

	return objectList
}

// PutObject upload file
func PutObject(bucketName, prefix string, fileHeader *multipart.FileHeader) error {
	ctx := context.Background()
	putObjectOptions := minio.PutObjectOptions{ContentType: fileHeader.Header.Get("content-type")}
	file, _ := fileHeader.Open()
	defer file.Close()

	_, err := client.PutObject(ctx, bucketName, prefix+fileHeader.Filename, file, fileHeader.Size, putObjectOptions)
	return err
}

// PresignedGetObject generates a presigned URL for HTTP GET operations
func PresignedGetObject(bucketName, prefix, fileName string, exp time.Duration) (*url.URL, error) {
	ctx := context.Background()

	reqParams := make(url.Values)
	reqParams.Set("response-content-disposition", "attachment; filename=\""+fileName+"\"")

	url, err := client.PresignedGetObject(ctx, bucketName, prefix+fileName, exp, reqParams)
	return url, err
}

// GetObject get object
func GetObject(bucketName, prefix, fileName string) (*minio.Object, error) {
	ctx := context.Background()

	return client.GetObject(ctx, bucketName, prefix+fileName, minio.GetObjectOptions{})
}
