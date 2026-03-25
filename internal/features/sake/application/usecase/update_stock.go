package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/sake-kasu/sake-hack-backend/internal/apperror"
	"github.com/sake-kasu/sake-hack-backend/internal/features/sake/domain/entity"
	"github.com/sake-kasu/sake-hack-backend/internal/features/sake/domain/repository"
	"github.com/sake-kasu/sake-hack-backend/internal/logger"
)

// UpdateStockInput 在庫更新の入力パラメータ
type UpdateStockInput struct {
	ID              uuid.UUID
	Category        entity.SakeCategory
	Kind            entity.SakeKind
	Brewery         entity.Brewery
	Name            entity.SakeName
	Abv             *float32
	PurchaseVolume  *int32
	RemainingVolume *int32
	Memo            *string
	DrinkStyles     []entity.DrinkStyle
	TagNames        []string
	Price           *int32
	ImageUrl        *string
	ImageKeys       []string
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
	sakeRepo repository.SakeRepository
}

// NewUpdateStockUsecase コンストラクタ
func NewUpdateStockUsecase(sakeRepo repository.SakeRepository) *UpdateStockUsecase {
	return &UpdateStockUsecase{sakeRepo: sakeRepo}
}

// Execute 在庫を更新する
func (u *UpdateStockUsecase) Execute(ctx context.Context, input UpdateStockInput) (*UpdateStockOutput, error) {
	defer logger.TraceMethodAuto(ctx, input)()

	if !input.Category.IsValid() {
		return nil, apperror.BadRequestError("無効なカテゴリです")
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
		TagNames:        input.TagNames,
		Price:           input.Price,
		ImageUrl:        input.ImageUrl,
		ImageKeys:       input.ImageKeys,
	})
	if err != nil {
		return nil, err
	}

	return &UpdateStockOutput{Sake: sake}, nil
}
