package usecase

import (
	"context"

	"go.uber.org/zap"

	"github.com/sake-kasu/sake-hack-backend/internal/apperror"
	"github.com/sake-kasu/sake-hack-backend/internal/features/sake/application/port"
	"github.com/sake-kasu/sake-hack-backend/internal/features/sake/application/query"
	"github.com/sake-kasu/sake-hack-backend/internal/features/sake/domain/entity"
	"github.com/sake-kasu/sake-hack-backend/internal/features/sake/domain/repository"
	"github.com/sake-kasu/sake-hack-backend/internal/logger"
)

// UpdateStockInput 在庫更新の入力パラメータ
type UpdateStockInput struct {
	ID              int32
	Category        entity.SakeCategory
	Kind            entity.SakeKind
	Brewery         entity.Brewery
	Name            entity.SakeName
	Abv             float32
	PurchaseVolume  float32
	RemainingVolume float32
	Memo            *string
	DrinkStyles     []entity.DrinkStyle
	Price           int32
	ObjectKey       *string
}

// UpdateStockOutput 在庫更新の出力
type UpdateStockOutput struct {
	Sake *entity.SakeListItem
}

// UpdateStockUsecaseInterface 在庫更新ユースケースのインターフェイス
type UpdateStockUsecaseInterface interface {
	Execute(ctx context.Context, input UpdateStockInput) (*UpdateStockOutput, error)
}

// UpdateStockUsecase 在庫更新ユースケース
type UpdateStockUsecase struct {
	sakeRepo    repository.SakeRepository
	sakeQuery   query.SakeQuery
	fileDeleter port.FileDeleter
}

// NewUpdateStockUsecase コンストラクタ
func NewUpdateStockUsecase(
	sakeRepo repository.SakeRepository,
	sakeQuery query.SakeQuery,
	fileDeleter port.FileDeleter,
) *UpdateStockUsecase {
	return &UpdateStockUsecase{
		sakeRepo:    sakeRepo,
		sakeQuery:   sakeQuery,
		fileDeleter: fileDeleter,
	}
}

// Execute 在庫を更新する
func (u *UpdateStockUsecase) Execute(ctx context.Context, input UpdateStockInput) (*UpdateStockOutput, error) {
	defer logger.TraceMethodAuto(ctx, input)()

	if !input.Category.IsValid() {
		return nil, apperror.BadRequestError("無効なカテゴリです")
	}

	// 更新前のobjectKeyを取得(差し替え時の旧ファイル削除用)
	oldObjectKey, err := u.sakeQuery.GetObjectKey(ctx, input.ID)
	if err != nil {
		return nil, err
	}

	sake, err := u.sakeRepo.Update(ctx, repository.UpdateSakeInput{
		ID:              input.ID,
		Category:        input.Category,
		Kind:            input.Kind,
		Brewery:         input.Brewery,
		Name:            input.Name,
		Abv:             input.Abv,
		PurchaseVolume:  input.PurchaseVolume,
		RemainingVolume: input.RemainingVolume,
		Memo:            input.Memo,
		DrinkStyles:     input.DrinkStyles,
		Price:           input.Price,
		ObjectKey:       input.ObjectKey,
	})
	if err != nil {
		return nil, err
	}

	// 旧ファイルが存在し、かつobjectKeyが変更された場合はストレージから削除する
	if oldObjectKey != nil && isObjectKeyChanged(oldObjectKey, input.ObjectKey) {
		if err := u.fileDeleter.DeleteObject(ctx, *oldObjectKey); err != nil {
			// ファイル削除失敗はDB更新成功後なので警告ログのみ(業務継続優先)
			logger.Warn(ctx, "旧ファイルのストレージ削除に失敗しました",
				zap.String("object_key", *oldObjectKey),
				zap.Error(err),
			)
		}
	}

	return &UpdateStockOutput{Sake: sake}, nil
}

// isObjectKeyChanged は旧objectKeyと新objectKeyが異なるかを判定する
func isObjectKeyChanged(oldKey *string, newKey *string) bool {
	if newKey == nil {
		return true
	}
	return *oldKey != *newKey
}
