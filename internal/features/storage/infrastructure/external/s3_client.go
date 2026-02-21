package external

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	"github.com/sake-kasu/sake-hack-backend/internal/config"
)

// コンパイル時にインターフェース実装を検証
var _ interface {
	GeneratePresignedPutURL(ctx context.Context, objectKey string, contentType string) (string, error)
	GeneratePresignedGetURL(ctx context.Context, objectKey string) (string, error)
	ResolveURL(ctx context.Context, objectKey string) (string, error)
} = (*S3Client)(nil)

// S3Client はS3/RustFS互換のストレージクライアント
type S3Client struct {
	client         *s3.Client
	presignClient  *s3.PresignClient
	bucketName     string
	uploadExpiry   time.Duration
	downloadExpiry time.Duration
}

// NewS3Client はS3クライアントを初期化する
func NewS3Client(ctx context.Context, cfg config.StorageConfig) (*S3Client, error) {
	optFns := []func(*awsconfig.LoadOptions) error{
		awsconfig.WithRegion(cfg.Region),
	}

	if cfg.AccessKey != "" && cfg.SecretKey != "" {
		optFns = append(optFns, awsconfig.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(cfg.AccessKey, cfg.SecretKey, ""),
		))
	}

	awsCfg, err := awsconfig.LoadDefaultConfig(ctx, optFns...)
	if err != nil {
		return nil, fmt.Errorf("AWS設定の読み込みに失敗しました: %w", err)
	}

	s3OptFns := []func(*s3.Options){}
	if cfg.Endpoint != "" {
		s3OptFns = append(s3OptFns, func(o *s3.Options) {
			o.BaseEndpoint = aws.String(cfg.Endpoint)
		})
	}
	if cfg.UsePathStyle {
		s3OptFns = append(s3OptFns, func(o *s3.Options) {
			o.UsePathStyle = true
		})
	}

	client := s3.NewFromConfig(awsCfg, s3OptFns...)
	presignClient := s3.NewPresignClient(client)

	return &S3Client{
		client:         client,
		presignClient:  presignClient,
		bucketName:     cfg.BucketName,
		uploadExpiry:   cfg.UploadExpiry,
		downloadExpiry: cfg.DownloadExpiry,
	}, nil
}

// GeneratePresignedPutURL は署名付きアップロードURLを発行する
func (c *S3Client) GeneratePresignedPutURL(ctx context.Context, objectKey string, contentType string) (string, error) {
	req, err := c.presignClient.PresignPutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(c.bucketName),
		Key:         aws.String(objectKey),
		ContentType: aws.String(contentType),
	}, s3.WithPresignExpires(c.uploadExpiry))
	if err != nil {
		return "", fmt.Errorf("署名付きPUT URLの生成に失敗しました: %w", err)
	}

	return req.URL, nil
}

// GeneratePresignedGetURL は署名付きダウンロードURLを発行する
func (c *S3Client) GeneratePresignedGetURL(ctx context.Context, objectKey string) (string, error) {
	req, err := c.presignClient.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(c.bucketName),
		Key:    aws.String(objectKey),
	}, s3.WithPresignExpires(c.downloadExpiry))
	if err != nil {
		return "", fmt.Errorf("署名付きGET URLの生成に失敗しました: %w", err)
	}

	return req.URL, nil
}

// ResolveURL はオブジェクトキーから署名付きダウンロードURLを生成する
// sake featureのURLResolverインターフェースを暗黙的に満たす
func (c *S3Client) ResolveURL(ctx context.Context, objectKey string) (string, error) {
	return c.GeneratePresignedGetURL(ctx, objectKey)
}
