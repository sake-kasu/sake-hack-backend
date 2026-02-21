package usecase

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/google/uuid"

	"github.com/sake-kasu/sake-hack-backend/internal/apperror"
	"github.com/sake-kasu/sake-hack-backend/internal/features/sake/application/port"
	"github.com/sake-kasu/sake-hack-backend/internal/logger"
)

// CreateUploadUrlInput 画像アップロードURL発行の入力
type CreateUploadUrlInput struct {
	SakeID      int32
	ContentType string
	Filename    string
}

// CreateUploadUrlOutput 画像アップロードURL発行の出力
type CreateUploadUrlOutput struct {
	UploadURL string
	ObjectKey string
}

// CreateUploadUrlUsecaseInterface 画像アップロードURL発行ユースケースのインターフェース
type CreateUploadUrlUsecaseInterface interface {
	Execute(ctx context.Context, input CreateUploadUrlInput) (*CreateUploadUrlOutput, error)
}

// CreateUploadUrlUsecase 画像アップロードURL発行ユースケース
type CreateUploadUrlUsecase struct {
	imageUploader port.ImageUploader
}

// NewCreateUploadUrlUsecase コンストラクタ
func NewCreateUploadUrlUsecase(imageUploader port.ImageUploader) *CreateUploadUrlUsecase {
	return &CreateUploadUrlUsecase{
		imageUploader: imageUploader,
	}
}

// Execute 画像アップロード用署名付きURLを発行する
// DBへのobjectKey保存は行わない。クライアントがS3アップロード成功後にPATCH /stocks/{id}で更新する。
func (u *CreateUploadUrlUsecase) Execute(ctx context.Context, input CreateUploadUrlInput) (*CreateUploadUrlOutput, error) {
	defer logger.TraceMethodAuto(ctx, input)()

	ext := filepath.Ext(input.Filename)
	if ext == "" {
		ext = extensionFromContentType(input.ContentType)
	}

	objectKey := fmt.Sprintf("sakes/%d/%s%s", input.SakeID, uuid.New().String(), ext)

	uploadURL, err := u.imageUploader.GeneratePresignedPutURL(ctx, objectKey, input.ContentType)
	if err != nil {
		return nil, apperror.InternalServerError("署名付きアップロードURLの生成に失敗しました").WithErr(err)
	}

	return &CreateUploadUrlOutput{
		UploadURL: uploadURL,
		ObjectKey: objectKey,
	}, nil
}

// extensionFromContentType はContent-Typeから拡張子を取得する
func extensionFromContentType(contentType string) string {
	switch contentType {
	case "image/jpeg":
		return ".jpg"
	case "image/png":
		return ".png"
	case "image/gif":
		return ".gif"
	case "image/webp":
		return ".webp"
	default:
		return ""
	}
}
