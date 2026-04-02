package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/sake-kasu/sake-hack-backend/internal/features/sake/domain/repository"
	"github.com/sake-kasu/sake-hack-backend/internal/logger"
)

// DeleteStockUsecaseInterface 在庫削除ユースケースのインターフェイス
type DeleteStockUsecaseInterface interface {
	Execute(ctx context.Context, id uuid.UUID) error
}

// DeleteStockUsecase 在庫削除ユースケース
type DeleteStockUsecase struct {
	sakeRepo repository.SakeRepository
}

// NewDeleteStockUsecase コンストラクタ
func NewDeleteStockUsecase(sakeRepo repository.SakeRepository) *DeleteStockUsecase {
	return &DeleteStockUsecase{sakeRepo: sakeRepo}
}

// Execute 在庫を削除する
func (u *DeleteStockUsecase) Execute(ctx context.Context, id uuid.UUID) error {
	defer logger.TraceMethodAuto(ctx, id)()

	return u.sakeRepo.Delete(ctx, id)
}
