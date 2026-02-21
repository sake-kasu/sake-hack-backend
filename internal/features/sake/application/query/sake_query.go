package query

import (
	"context"

	"github.com/sake-kasu/sake-hack-backend/internal/features/sake/domain/entity"
)

// ListSakesFilter 酒一覧取得のフィルター条件
type ListSakesFilter struct {
	KindID    *int32
	BreweryID *int32
	LikeToken *string
	Offset    int32
	Limit     int32
}

// SakeQuery 酒クエリのインターフェース(読み取り操作)
type SakeQuery interface {
	List(ctx context.Context, filter ListSakesFilter) ([]entity.SakeListItem, entity.Pagination, error)
	GetDetail(ctx context.Context, id int32, likeToken *string) (*entity.SakeDetail, error)
	ListKinds(ctx context.Context) ([]entity.SakeKind, error)
	ListBreweries(ctx context.Context, keyword *string, limit int32) ([]entity.Brewery, error)
	ListDrinkStyles(ctx context.Context) ([]entity.DrinkStyle, error)
}
