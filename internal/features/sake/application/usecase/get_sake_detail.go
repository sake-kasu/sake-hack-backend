package usecase

import (
	"context"

	"github.com/sake-kasu/sake-hack-backend/internal/features/sake/application/query"
	"github.com/sake-kasu/sake-hack-backend/internal/features/sake/domain/entity"
	"github.com/sake-kasu/sake-hack-backend/internal/logger"
)

// GetSakeDetailOutput 酒詳細取得の出力
type GetSakeDetailOutput struct {
	Detail *entity.SakeDetail
}

// GetSakeDetailInput 酒詳細取得の入力パラメータ
type GetSakeDetailInput struct {
	ID        int32
	LikeToken *string
}

// GetSakeDetailUsecaseInterface 酒詳細取得ユースケースのインターフェイス
type GetSakeDetailUsecaseInterface interface {
	Execute(ctx context.Context, input GetSakeDetailInput) (*GetSakeDetailOutput, error)
}

// GetSakeDetailUsecase 酒詳細取得ユースケース
type GetSakeDetailUsecase struct {
	sakeQuery query.SakeQuery
}

// NewGetSakeDetailUsecase コンストラクタ
func NewGetSakeDetailUsecase(sakeQuery query.SakeQuery) *GetSakeDetailUsecase {
	return &GetSakeDetailUsecase{sakeQuery: sakeQuery}
}

// Execute 酒の詳細を取得する
func (u *GetSakeDetailUsecase) Execute(ctx context.Context, input GetSakeDetailInput) (*GetSakeDetailOutput, error) {
	defer logger.TraceMethodAuto(ctx, input)()

	detail, err := u.sakeQuery.GetDetail(ctx, input.ID, input.LikeToken)
	if err != nil {
		return nil, err
	}

	return &GetSakeDetailOutput{Detail: detail}, nil
}
