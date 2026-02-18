package query

import (
	"context"

	"github.com/sake-kasu/sake-hack-backend/internal/features/sake/domain/entity"
)

// ListSakesFilter 酒一覧取得のフィルター条件
type ListSakesFilter struct {
	Category *string
	Search   *string
	Offset   int32
	Limit    int32
}

// SakeQuery 酒クエリのインターフェース(読み取り操作)
type SakeQuery interface {
	List(ctx context.Context, filter ListSakesFilter) ([]entity.SakeDetail, entity.Pagination, error)
	GetDetail(ctx context.Context, id int32) (*entity.SakeDetail, error)
}
