package usecase

import (
	"context"

	"github.com/sake-kasu/sake-hack-backend/internal/features/sake/application/query"
	"github.com/sake-kasu/sake-hack-backend/internal/features/sake/domain/entity"
	"github.com/sake-kasu/sake-hack-backend/internal/logger"
)

// ListDrinkStylesOutput 飲み方一覧取得の出力
type ListDrinkStylesOutput struct {
	DrinkStyles []entity.DrinkStyle
}

// ListDrinkStylesUsecaseInterface 飲み方一覧取得ユースケースのインターフェイス
type ListDrinkStylesUsecaseInterface interface {
	Execute(ctx context.Context) (*ListDrinkStylesOutput, error)
}

// ListDrinkStylesUsecase 飲み方一覧取得ユースケース
type ListDrinkStylesUsecase struct {
	sakeQuery query.SakeQuery
}

// NewListDrinkStylesUsecase コンストラクタ
func NewListDrinkStylesUsecase(sakeQuery query.SakeQuery) *ListDrinkStylesUsecase {
	return &ListDrinkStylesUsecase{sakeQuery: sakeQuery}
}

// Execute 飲み方一覧を取得する
func (u *ListDrinkStylesUsecase) Execute(ctx context.Context) (*ListDrinkStylesOutput, error) {
	defer logger.TraceMethodAuto(ctx, nil)()

	drinkStyles, err := u.sakeQuery.ListDrinkStyles(ctx)
	if err != nil {
		return nil, err
	}

	return &ListDrinkStylesOutput{DrinkStyles: drinkStyles}, nil
}
