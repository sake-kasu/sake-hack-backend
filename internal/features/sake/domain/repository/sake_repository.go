package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/sake-kasu/sake-hack-backend/internal/features/sake/domain/entity"
)

// CreateSakeInput 酒の作成に必要な入力
type CreateSakeInput struct {
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

// UpdateSakeInput 酒の更新に必要な入力
type UpdateSakeInput struct {
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

// SakeRepository 酒リポジトリのインターフェース(書き込み操作)
type SakeRepository interface {
	Create(ctx context.Context, input CreateSakeInput) (*entity.SakeListItem, error)
	Update(ctx context.Context, input UpdateSakeInput) (*entity.SakeListItem, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
