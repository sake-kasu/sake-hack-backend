package usecase

import (
	"context"

	"go.uber.org/zap"

	"github.com/sake-kasu/sake-hack-backend/internal/features/sake/application/port"
	"github.com/sake-kasu/sake-hack-backend/internal/features/sake/application/query"
	"github.com/sake-kasu/sake-hack-backend/internal/features/sake/domain/repository"
	"github.com/sake-kasu/sake-hack-backend/internal/logger"
)

// DeleteStockUsecaseInterface 在庫削除ユースケースのインターフェイス
type DeleteStockUsecaseInterface interface {
	Execute(ctx context.Context, id int32) error
}

// DeleteStockUsecase 在庫削除ユースケース
type DeleteStockUsecase struct {
	sakeRepo    repository.SakeRepository
	sakeQuery   query.SakeQuery
	fileDeleter port.FileDeleter
}

// NewDeleteStockUsecase コンストラクタ
func NewDeleteStockUsecase(
	sakeRepo repository.SakeRepository,
	sakeQuery query.SakeQuery,
	fileDeleter port.FileDeleter,
) *DeleteStockUsecase {
	return &DeleteStockUsecase{
		sakeRepo:    sakeRepo,
		sakeQuery:   sakeQuery,
		fileDeleter: fileDeleter,
	}
}

// Execute 在庫を削除する
func (u *DeleteStockUsecase) Execute(ctx context.Context, id int32) error {
	defer logger.TraceMethodAuto(ctx, id)()

	// 削除前にobjectKeyを取得(ストレージのファイル削除用)
	objectKey, err := u.sakeQuery.GetObjectKey(ctx, id)
	if err != nil {
		return err
	}

	if err := u.sakeRepo.Delete(ctx, id); err != nil {
		return err
	}

	// DB削除成功後、ストレージのファイルを削除する
	if objectKey != nil {
		if err := u.fileDeleter.DeleteObject(ctx, *objectKey); err != nil {
			// ファイル削除失敗はDB削除成功後なので警告ログのみ(業務継続優先)
			logger.Warn(ctx, "ストレージのファイル削除に失敗しました",
				zap.String("object_key", *objectKey),
				zap.Error(err),
			)
		}
	}

	return nil
}
