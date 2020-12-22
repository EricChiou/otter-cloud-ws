package minio

import (
	"context"
	"time"

	"github.com/minio/minio-go/v7"
)

// Object file struct
type Object struct {
	ContentType  string
	Name         string
	Size         int64
	LastModified time.Time
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
		Recursive: true,
	}

	var objectList []Object
	for object := range client.ListObjects(ctx, bucketName, opts) {
		if object.Err == nil {
			objectInfo, _ := client.StatObject(ctx, bucketName, object.Key, minio.StatObjectOptions{})
			objectList = append(objectList, Object{
				ContentType:  objectInfo.ContentType,
				Name:         objectInfo.Key,
				Size:         objectInfo.Size,
				LastModified: objectInfo.LastModified,
			})
		}
	}

	return objectList
}
