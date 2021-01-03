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
	file, _ := fileHeader.Open()
	defer file.Close()

	putObjectOptions := minio.PutObjectOptions{ContentType: fileHeader.Header.Get("content-type")}

	_, err := client.PutObject(
		context.Background(),
		bucketName,
		prefix+fileHeader.Filename,
		file,
		fileHeader.Size,
		putObjectOptions,
	)
	return err
}

// PresignedGetObject generates a presigned URL for HTTP GET operations
func PresignedGetObject(bucketName, prefix, fileName string, exp time.Duration) (*url.URL, error) {
	reqParams := make(url.Values)
	reqParams.Set("response-content-disposition", "attachment; filename=\""+fileName+"\"")

	return client.PresignedGetObject(
		context.Background(),
		bucketName,
		prefix+fileName,
		exp,
		reqParams,
	)
}

// GetObject get object
func GetObject(bucketName, prefix, fileName string) (*minio.Object, error) {
	opts := minio.GetObjectOptions{}

	return client.GetObject(context.Background(), bucketName, prefix+fileName, opts)
}

// RemoveObject removes an object with some specified options
func RemoveObject(bucketName, prefix, fileName string) error {
	opts := minio.RemoveObjectOptions{
		GovernanceBypass: true,
	}

	return client.RemoveObject(context.Background(), bucketName, prefix+fileName, opts)
}

// RemoveObjects removes a list of objects obtained from an input channel
func RemoveObjects(bucketName, prefix string) error {
	objectsCh := make(chan minio.ObjectInfo)
	go func() {
		defer close(objectsCh)

		opts := minio.ListObjectsOptions{
			Prefix:    prefix,
			Recursive: true,
		}
		for object := range client.ListObjects(context.Background(), bucketName, opts) {
			if object.Err == nil {
				objectsCh <- object
			}
		}
	}()

	opts := minio.RemoveObjectsOptions{
		GovernanceBypass: true,
	}

	for err := range client.RemoveObjects(context.Background(), bucketName, objectsCh, opts) {
		if err.Err != nil {
			return err.Err
		}
	}

	return nil
}
