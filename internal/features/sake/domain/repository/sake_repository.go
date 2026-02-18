package repository

import (
	"context"

	"github.com/sake-kasu/sake-hack-backend/internal/features/sake/domain/entity"
)

// CreateSakeInput 酒の作成に必要な入力
type CreateSakeInput struct {
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

// UpdateSakeInput 酒の更新に必要な入力
type UpdateSakeInput struct {
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
	ImageUrl        *string
}

// SakeRepository 酒リポジトリのインターフェース(書き込み操作)
type SakeRepository interface {
	Create(ctx context.Context, input CreateSakeInput) (*entity.SakeListItem, error)
	Update(ctx context.Context, input UpdateSakeInput) (*entity.SakeListItem, error)
	Delete(ctx context.Context, id int32) error
}
