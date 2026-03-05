package config

import (
	"context"
	"log"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func ConnectMinIO(cfg *Config) *minio.Client {
	client, err := minio.New(cfg.MinioEndpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.MinioAccessKey, cfg.MinioSecretKey, ""),
		Secure: false,
	})
	if err != nil {
		log.Fatalf("❌ Failed to connect to MinIO: %v", err)
	}

	ctx := context.Background()
	exists, err := client.BucketExists(ctx, cfg.MinioBucket)
	if err != nil {
		log.Fatalf("❌ Failed to check bucket: %v", err)
	}
	if !exists {
		err = client.MakeBucket(ctx, cfg.MinioBucket, minio.MakeBucketOptions{})
		if err != nil {
			log.Fatalf("❌ Failed to create bucket: %v", err)
		}
		log.Printf("✅ MinIO bucket '%s' created", cfg.MinioBucket)
	}

	log.Println("✅ MinIO connected successfully!")
	return client
}
