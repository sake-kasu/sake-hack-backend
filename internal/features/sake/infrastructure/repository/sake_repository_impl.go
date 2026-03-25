package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/sake-kasu/sake-hack-backend/internal/apperror"
	"github.com/sake-kasu/sake-hack-backend/internal/features/sake/domain/entity"
	"github.com/sake-kasu/sake-hack-backend/internal/features/sake/domain/repository"
	"github.com/sake-kasu/sake-hack-backend/internal/logger"
)

// sakeRepositoryImpl 酒リポジトリの実装
type sakeRepositoryImpl struct{}

// NewSakeRepository コンストラクタ
func NewSakeRepository(_ *pgxpool.Pool) repository.SakeRepository {
	return &sakeRepositoryImpl{}
}

// Create はDB/OASの再設計に伴う再実装まで一時的に未提供。
func (r *sakeRepositoryImpl) Create(ctx context.Context, input repository.CreateSakeInput) (*entity.SakeListItem, error) {
	defer logger.TraceMethodAuto(ctx, input)()

	return nil, apperror.InternalServerError("sake repository create is temporarily disabled")
}

// Update はDB/OASの再設計に伴う再実装まで一時的に未提供。
func (r *sakeRepositoryImpl) Update(ctx context.Context, input repository.UpdateSakeInput) (*entity.SakeListItem, error) {
	defer logger.TraceMethodAuto(ctx, input)()

	return nil, apperror.InternalServerError("sake repository update is temporarily disabled")
}

// Delete はDB/OASの再設計に伴う再実装まで一時的に未提供。
func (r *sakeRepositoryImpl) Delete(ctx context.Context, id int32) error {
	defer logger.TraceMethodAuto(ctx, id)()

	return apperror.InternalServerError("sake repository delete is temporarily disabled")
}
