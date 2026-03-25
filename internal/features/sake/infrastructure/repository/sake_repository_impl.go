package repository

import (
	"context"
	"strconv"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/sake-kasu/sake-hack-backend/internal/apperror"
	"github.com/sake-kasu/sake-hack-backend/internal/database/sqlc"
	"github.com/sake-kasu/sake-hack-backend/internal/features/sake/domain/entity"
	"github.com/sake-kasu/sake-hack-backend/internal/features/sake/domain/repository"
	"github.com/sake-kasu/sake-hack-backend/internal/logger"
)

// sakeRepositoryImpl 酒リポジトリの実装
type sakeRepositoryImpl struct {
	queries *sqlc.Queries
}

// NewSakeRepository コンストラクタ
func NewSakeRepository(db *pgxpool.Pool) repository.SakeRepository {
	return &sakeRepositoryImpl{
		queries: sqlc.New(db),
	}
}

// Create は現行DB schemaに在庫データを登録する。
func (r *sakeRepositoryImpl) Create(ctx context.Context, input repository.CreateSakeInput) (*entity.SakeListItem, error) {
	defer logger.TraceMethodAuto(ctx, input)()

	row, err := r.queries.CreateSake(ctx, sqlc.CreateSakeParams{
		Category:          sqlc.SakeCategory(input.Category),
		Name:              input.Name.Name,
		Phonetic:          emptyStringToNil(input.Name.Phonetic),
		AlcoholPercentage: numericFromFloat32Ptr(input.Abv),
		VolumeMax:         input.PurchaseVolume,
		VolumeRemain:      input.RemainingVolume,
		Region:            input.Brewery.OriginRegion,
		Price:             input.Price,
		Memo:              input.Memo,
	})
	if err != nil {
		logger.LogDatabaseError(ctx, "INSERT", "sakes", err, map[string]interface{}{"input": input})
		return nil, apperror.DatabaseError("酒の登録に失敗しました", err)
	}

	return &entity.SakeListItem{
		ID:           uuid.UUID(row.ID.Bytes),
		Category:     entity.SakeCategory(row.Category),
		Name:         row.Name,
		ImagePreview: "",
	}, nil
}

// Update はDB/OASの再設計に伴う再実装まで一時的に未提供。
func (r *sakeRepositoryImpl) Update(ctx context.Context, input repository.UpdateSakeInput) (*entity.SakeListItem, error) {
	defer logger.TraceMethodAuto(ctx, input)()

	return nil, apperror.InternalServerError("sake repository update is temporarily disabled")
}

// Delete はDB/OASの再設計に伴う再実装まで一時的に未提供。
func (r *sakeRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	defer logger.TraceMethodAuto(ctx, id)()

	return apperror.InternalServerError("sake repository delete is temporarily disabled")
}

// emptyStringToNil は空文字を nil に正規化する。
func emptyStringToNil(v string) *string {
	if v == "" {
		return nil
	}
	return &v
}

// numericFromFloat32Ptr は nullable な float を pgtype.Numeric に変換する。
func numericFromFloat32Ptr(v *float32) pgtype.Numeric {
	if v == nil {
		return pgtype.Numeric{}
	}

	var numeric pgtype.Numeric
	_ = numeric.ScanScientific(strconv.FormatFloat(float64(*v), 'f', -1, 32))
	return numeric
}
