package usecase

import (
	"context"

	"github.com/sake-kasu/sake-hack-backend/internal/apperror"
	"github.com/sake-kasu/sake-hack-backend/internal/features/sake/domain/entity"
	"github.com/sake-kasu/sake-hack-backend/internal/features/sake/domain/repository"
	"github.com/sake-kasu/sake-hack-backend/internal/logger"
)

// CreateStockInput 在庫登録の入力パラメータ
type CreateStockInput struct {
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
	ImageUrl        *string
}

// CreateStockOutput 在庫登録の出力
type CreateStockOutput struct {
	Sake *entity.SakeListItem
}

// CreateStockUsecaseInterface 在庫登録ユースケースのインターフェイス
type CreateStockUsecaseInterface interface {
	Execute(ctx context.Context, input CreateStockInput) (*CreateStockOutput, error)
}

// CreateStockUsecase 在庫登録ユースケース
type CreateStockUsecase struct {
	sakeRepo repository.SakeRepository
}

// NewCreateStockUsecase コンストラクタ
func NewCreateStockUsecase(sakeRepo repository.SakeRepository) *CreateStockUsecase {
	return &CreateStockUsecase{sakeRepo: sakeRepo}
}

// Execute 在庫を登録する
func (u *CreateStockUsecase) Execute(ctx context.Context, input CreateStockInput) (*CreateStockOutput, error) {
	defer logger.TraceMethodAuto(ctx, input)()

	if !input.Category.IsValid() {
		return nil, apperror.BadRequestError("無効なカテゴリです")
	}

	sake, err := u.sakeRepo.Create(ctx, repository.CreateSakeInput{
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
		ImageUrl:        input.ImageUrl,
	})
	if err != nil {
		return nil, err
	}

	return &CreateStockOutput{Sake: sake}, nil
}
