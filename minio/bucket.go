package minio

import (
	"context"
	"otter-cloud-ws/config"
	"otter-cloud-ws/service/sha3"

	"github.com/minio/minio-go/v7"
)

// GetUserBucketID get user minio bucket id
func GetUserBucketID(acc string) string {
	cfg := config.Get()

	return sha3.Encrypt(acc+"_"+cfg.BucketHashKey, 224)
}

// CreateUserBucket create user minio bucket
func CreateUserBucket(acc string) error {
	bucketID := GetUserBucketID(acc)

	return client.MakeBucket(
		context.Background(),
		bucketID,
		minio.MakeBucketOptions{Region: "ap-east-1"},
	)
}
