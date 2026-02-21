package usecase

import (
	"context"

	"github.com/sake-kasu/sake-hack-backend/internal/features/sake/application/query"
	"github.com/sake-kasu/sake-hack-backend/internal/features/sake/domain/entity"
	"github.com/sake-kasu/sake-hack-backend/internal/logger"
)

// ListKindsOutput 酒の種類一覧取得の出力
type ListKindsOutput struct {
	Kinds []entity.SakeKind
}

// ListKindsUsecaseInterface 酒の種類一覧取得ユースケースのインターフェイス
type ListKindsUsecaseInterface interface {
	Execute(ctx context.Context) (*ListKindsOutput, error)
}

// ListKindsUsecase 酒の種類一覧取得ユースケース
type ListKindsUsecase struct {
	sakeQuery query.SakeQuery
}

// NewListKindsUsecase コンストラクタ
func NewListKindsUsecase(sakeQuery query.SakeQuery) *ListKindsUsecase {
	return &ListKindsUsecase{sakeQuery: sakeQuery}
}

// Execute 酒の種類一覧を取得する
func (u *ListKindsUsecase) Execute(ctx context.Context) (*ListKindsOutput, error) {
	defer logger.TraceMethodAuto(ctx, nil)()

	kinds, err := u.sakeQuery.ListKinds(ctx)
	if err != nil {
		return nil, err
	}

	return &ListKindsOutput{Kinds: kinds}, nil
}
