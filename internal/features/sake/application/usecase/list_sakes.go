package usecase

import (
	"context"

	"github.com/sake-kasu/sake-hack-backend/internal/features/sake/application/query"
	"github.com/sake-kasu/sake-hack-backend/internal/features/sake/domain/entity"
	"github.com/sake-kasu/sake-hack-backend/internal/logger"
)

// ListSakesInput 酒一覧取得の入力パラメータ
type ListSakesInput struct {
	Category *string
	Search   *string
	Offset   int32
	Limit    int32
}

// ListSakesOutput 酒一覧取得の出力
type ListSakesOutput struct {
	Sakes []entity.SakeDetail
	Total int64
}

// ListSakesUsecaseInterface 酒一覧取得ユースケースのインターフェイス
type ListSakesUsecaseInterface interface {
	Execute(ctx context.Context, input ListSakesInput) (*ListSakesOutput, error)
}

// ListSakesUsecase 酒一覧取得ユースケース
type ListSakesUsecase struct {
	sakeQuery query.SakeQuery
}

// NewListSakesUsecase コンストラクタ
func NewListSakesUsecase(sakeQuery query.SakeQuery) *ListSakesUsecase {
	return &ListSakesUsecase{sakeQuery: sakeQuery}
}

// Execute 酒一覧を取得する
func (u *ListSakesUsecase) Execute(ctx context.Context, input ListSakesInput) (*ListSakesOutput, error) {
	defer logger.TraceMethodAuto(ctx, input)()

	offset := input.Offset
	if offset < 0 {
		offset = 0
	}

	limit := input.Limit
	if limit < 1 || limit > 100 {
		limit = 20
	}

	sakes, pagination, err := u.sakeQuery.ListPublic(ctx, query.ListPublicSakesFilter{
		Category: input.Category,
		Search:   input.Search,
		Offset:   offset,
		Limit:    limit,
	})
	if err != nil {
		return nil, err
	}

	return &ListSakesOutput{
		Sakes: sakes,
		Total: pagination.Total,
	}, nil
}
