package usecase

import (
	"context"

	"github.com/sake-kasu/sake-hack-backend/internal/features/like/domain/repository"
	"github.com/sake-kasu/sake-hack-backend/internal/logger"
)

// ToggleLikeInput いいね追加の入力
type ToggleLikeInput struct {
	SakeID int32
	Token  string
}

// ToggleLikeOutput いいね追加の出力
type ToggleLikeOutput struct {
	LikeCount int64
	IsLiked   bool
}

// ToggleLikeUsecaseInterface いいね追加ユースケースのインターフェース
type ToggleLikeUsecaseInterface interface {
	Execute(ctx context.Context, input ToggleLikeInput) (*ToggleLikeOutput, error)
}

// ToggleLikeUsecase いいね追加ユースケース
type ToggleLikeUsecase struct {
	likeRepo repository.LikeRepository
}

// NewToggleLikeUsecase コンストラクタ
func NewToggleLikeUsecase(likeRepo repository.LikeRepository) *ToggleLikeUsecase {
	return &ToggleLikeUsecase{likeRepo: likeRepo}
}

// Execute いいねを追加し、現在のいいね数を返す
func (u *ToggleLikeUsecase) Execute(ctx context.Context, input ToggleLikeInput) (*ToggleLikeOutput, error) {
	defer logger.TraceMethodAuto(ctx, input)()

	if err := u.likeRepo.Create(ctx, input.SakeID, input.Token); err != nil {
		return nil, err
	}

	count, err := u.likeRepo.GetCount(ctx, input.SakeID)
	if err != nil {
		return nil, err
	}

	isLiked, err := u.likeRepo.Exists(ctx, input.SakeID, input.Token)
	if err != nil {
		return nil, err
	}

	return &ToggleLikeOutput{
		LikeCount: count,
		IsLiked:   isLiked,
	}, nil
}
