package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/sake-kasu/sake-hack-backend/internal/apperror"
	"github.com/sake-kasu/sake-hack-backend/internal/database/sqlc"
	domainRepo "github.com/sake-kasu/sake-hack-backend/internal/features/like/domain/repository"
	"github.com/sake-kasu/sake-hack-backend/internal/logger"
)

type likeRepositoryImpl struct {
	queries *sqlc.Queries
}

// NewLikeRepository LikeRepositoryの実装を作成
func NewLikeRepository(db *pgxpool.Pool) domainRepo.LikeRepository {
	return &likeRepositoryImpl{
		queries: sqlc.New(db),
	}
}

// Create いいねを作成する
func (r *likeRepositoryImpl) Create(ctx context.Context, sakeID int32, token string) error {
	defer logger.TraceMethodAuto(ctx, map[string]interface{}{"sake_id": sakeID})()

	if err := r.queries.CreateLike(ctx, sqlc.CreateLikeParams{
		SakeID: sakeID,
		Token:  token,
	}); err != nil {
		logger.LogDatabaseError(ctx, "INSERT", "likes", err, map[string]interface{}{
			"sake_id": sakeID,
		})
		return apperror.DatabaseError("いいねの登録に失敗しました", err)
	}
	return nil
}

// Delete いいねを削除する
func (r *likeRepositoryImpl) Delete(ctx context.Context, sakeID int32, token string) (int64, error) {
	defer logger.TraceMethodAuto(ctx, map[string]interface{}{"sake_id": sakeID})()

	rows, err := r.queries.DeleteLike(ctx, sqlc.DeleteLikeParams{
		SakeID: sakeID,
		Token:  token,
	})
	if err != nil {
		logger.LogDatabaseError(ctx, "DELETE", "likes", err, map[string]interface{}{
			"sake_id": sakeID,
		})
		return 0, apperror.DatabaseError("いいねの削除に失敗しました", err)
	}
	return rows, nil
}

// GetCount 酒のいいね数を取得する
func (r *likeRepositoryImpl) GetCount(ctx context.Context, sakeID int32) (int64, error) {
	defer logger.TraceMethodAuto(ctx, map[string]interface{}{"sake_id": sakeID})()

	count, err := r.queries.GetLikeCount(ctx, sakeID)
	if err != nil {
		logger.LogDatabaseError(ctx, "SELECT", "likes", err, map[string]interface{}{
			"sake_id": sakeID,
		})
		return 0, apperror.DatabaseError("いいね数の取得に失敗しました", err)
	}
	return count, nil
}

// Exists いいねの存在確認
func (r *likeRepositoryImpl) Exists(ctx context.Context, sakeID int32, token string) (bool, error) {
	defer logger.TraceMethodAuto(ctx, map[string]interface{}{"sake_id": sakeID})()

	exists, err := r.queries.ExistsLike(ctx, sqlc.ExistsLikeParams{
		SakeID: sakeID,
		Token:  token,
	})
	if err != nil {
		logger.LogDatabaseError(ctx, "SELECT", "likes", err, map[string]interface{}{
			"sake_id": sakeID,
		})
		return false, apperror.DatabaseError("いいね存在確認に失敗しました", err)
	}
	return exists, nil
}
