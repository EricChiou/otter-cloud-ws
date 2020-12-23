package minio

import (
	"context"
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
