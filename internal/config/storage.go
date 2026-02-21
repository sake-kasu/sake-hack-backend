package config

import "time"

// StorageConfig はS3/RustFSストレージの設定
type StorageConfig struct {
	BucketName     string        `env:"STORAGE_BUCKET_NAME" envDefault:"sake-hack-bucket"`
	Region         string        `env:"STORAGE_REGION" envDefault:"ap-northeast-1"`
	Endpoint       string        `env:"STORAGE_ENDPOINT"`
	AccessKey      string        `env:"STORAGE_ACCESS_KEY"`
	SecretKey      string        `env:"STORAGE_SECRET_KEY"`
	UsePathStyle   bool          `env:"STORAGE_USE_PATH_STYLE" envDefault:"false"`
	UploadExpiry   time.Duration `env:"STORAGE_UPLOAD_EXPIRY" envDefault:"15m"`
	DownloadExpiry time.Duration `env:"STORAGE_DOWNLOAD_EXPIRY" envDefault:"1h"`
}
