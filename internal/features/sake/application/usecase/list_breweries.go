package usecase

import (
	"context"

	"github.com/sake-kasu/sake-hack-backend/internal/features/sake/application/query"
	"github.com/sake-kasu/sake-hack-backend/internal/features/sake/domain/entity"
	"github.com/sake-kasu/sake-hack-backend/internal/logger"
)

// ListBreweriesInput 酒造一覧取得の入力パラメータ
type ListBreweriesInput struct {
	Keyword *string
	Limit   int32
}

// ListBreweriesOutput 酒造一覧取得の出力
type ListBreweriesOutput struct {
	Breweries []entity.Brewery
}

// ListBreweriesUsecaseInterface 酒造一覧取得ユースケースのインターフェイス
type ListBreweriesUsecaseInterface interface {
	Execute(ctx context.Context, input ListBreweriesInput) (*ListBreweriesOutput, error)
}

// ListBreweriesUsecase 酒造一覧取得ユースケース
type ListBreweriesUsecase struct {
	sakeQuery query.SakeQuery
}

// NewListBreweriesUsecase コンストラクタ
func NewListBreweriesUsecase(sakeQuery query.SakeQuery) *ListBreweriesUsecase {
	return &ListBreweriesUsecase{sakeQuery: sakeQuery}
}

// Execute 酒造一覧を取得する
func (u *ListBreweriesUsecase) Execute(ctx context.Context, input ListBreweriesInput) (*ListBreweriesOutput, error) {
	defer logger.TraceMethodAuto(ctx, input)()

	limit := input.Limit
	if limit < 1 || limit > 100 {
		limit = 50
	}

	breweries, err := u.sakeQuery.ListBreweries(ctx, input.Keyword, limit)
	if err != nil {
		return nil, err
	}

	return &ListBreweriesOutput{Breweries: breweries}, nil
}
