package usecase

import (
	"context"

	"github.com/sake-kasu/sake-hack-backend/internal/apperror"
	"github.com/sake-kasu/sake-hack-backend/internal/features/like/domain/repository"
	"github.com/sake-kasu/sake-hack-backend/internal/logger"
)

// RemoveLikeInput いいね削除の入力
type RemoveLikeInput struct {
	SakeID int32
	Token  string
}

// RemoveLikeOutput いいね削除の出力
type RemoveLikeOutput struct {
	LikeCount int64
	IsLiked   bool
}

// RemoveLikeUsecaseInterface いいね削除ユースケースのインターフェース
type RemoveLikeUsecaseInterface interface {
	Execute(ctx context.Context, input RemoveLikeInput) (*RemoveLikeOutput, error)
}

// RemoveLikeUsecase いいね削除ユースケース
type RemoveLikeUsecase struct {
	likeRepo repository.LikeRepository
}

// NewRemoveLikeUsecase コンストラクタ
func NewRemoveLikeUsecase(likeRepo repository.LikeRepository) *RemoveLikeUsecase {
	return &RemoveLikeUsecase{likeRepo: likeRepo}
}

// Execute いいねを削除し、現在のいいね数を返す
func (u *RemoveLikeUsecase) Execute(ctx context.Context, input RemoveLikeInput) (*RemoveLikeOutput, error) {
	defer logger.TraceMethodAuto(ctx, input)()

	rowsAffected, err := u.likeRepo.Delete(ctx, input.SakeID, input.Token)
	if err != nil {
		return nil, err
	}

	if rowsAffected == 0 {
		return nil, apperror.NotFoundError("いいねが見つかりません")
	}

	count, err := u.likeRepo.GetCount(ctx, input.SakeID)
	if err != nil {
		return nil, err
	}

	return &RemoveLikeOutput{
		LikeCount: count,
		IsLiked:   false,
	}, nil
}
