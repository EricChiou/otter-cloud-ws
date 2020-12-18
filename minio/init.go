package minio

import (
	"otter-cloud-ws/config"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var client *minio.Client

// Init minioClient
func Init() error {
	cfg := config.Get()

	endpoint := cfg.EndPoint
	accessKeyID := cfg.AccessKeyID
	secretAccessKey := cfg.SecretAccessKey
	useSSL := cfg.UseSSL

	var err error
	client, err = minio.New(
		endpoint,
		&minio.Options{
			Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
			Secure: useSSL,
		},
	)

	return err
}
