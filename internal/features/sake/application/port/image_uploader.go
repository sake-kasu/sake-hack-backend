package port

import "context"

// ImageUploader は画像アップロード用署名付きURL発行のインターフェース
// sake featureがConsumerとして定義し、S3クライアントをDI経由で利用する
type ImageUploader interface {
	// GeneratePresignedPutURL は署名付きアップロードURLを発行する
	GeneratePresignedPutURL(ctx context.Context, objectKey string, contentType string) (string, error)
}
