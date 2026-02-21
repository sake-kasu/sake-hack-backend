package usecase

import (
	"context"

	"github.com/sake-kasu/sake-hack-backend/internal/features/sake/domain/repository"
	"github.com/sake-kasu/sake-hack-backend/internal/logger"
)

// PatchStockInput 在庫部分更新の入力
type PatchStockInput struct {
	ID        int32
	ObjectKey string
}

// PatchStockUsecaseInterface 在庫部分更新ユースケースのインターフェース
type PatchStockUsecaseInterface interface {
	Execute(ctx context.Context, input PatchStockInput) error
}

// PatchStockUsecase 在庫部分更新ユースケース
type PatchStockUsecase struct {
	sakeRepo repository.SakeRepository
}

// NewPatchStockUsecase コンストラクタ
func NewPatchStockUsecase(sakeRepo repository.SakeRepository) *PatchStockUsecase {
	return &PatchStockUsecase{sakeRepo: sakeRepo}
}

// Execute 在庫情報を部分更新する
func (u *PatchStockUsecase) Execute(ctx context.Context, input PatchStockInput) error {
	defer logger.TraceMethodAuto(ctx, input)()

	return u.sakeRepo.UpdateObjectKey(ctx, input.ID, input.ObjectKey)
}
