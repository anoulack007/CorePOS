package config

import (
	"context"
	"fmt"
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

	// Make the bucket public (Read-Only) so images can be fetched via URL
	policy := fmt.Sprintf(`{"Version": "2012-10-17","Statement": [{"Action": ["s3:GetObject"],"Effect": "Allow","Principal": {"AWS": ["*"]},"Resource": ["arn:aws:s3:::%s/*"]}]}`, cfg.MinioBucket)
	err = client.SetBucketPolicy(ctx, cfg.MinioBucket, policy)
	if err != nil {
		log.Printf("⚠️ Warning: Failed to set bucket public policy: %v", err)
	}

	log.Println("✅ MinIO connected successfully!")
	return client
}
