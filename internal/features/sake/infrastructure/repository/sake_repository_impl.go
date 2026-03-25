package repository

import (
	"context"
	"errors"
	"strconv"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
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
	db      *pgxpool.Pool
	queries *sqlc.Queries
}

// NewSakeRepository コンストラクタ
func NewSakeRepository(db *pgxpool.Pool) repository.SakeRepository {
	return &sakeRepositoryImpl{
		db:      db,
		queries: sqlc.New(db),
	}
}

// Create は現行DB schemaに在庫データを登録する。
func (r *sakeRepositoryImpl) Create(ctx context.Context, input repository.CreateSakeInput) (*entity.SakeListItem, error) {
	defer logger.TraceMethodAuto(ctx, input)()

	tx, err := r.db.Begin(ctx)
	if err != nil {
		logger.LogDatabaseError(ctx, "BEGIN", "sakes", err, nil)
		return nil, apperror.DatabaseError("トランザクション開始に失敗しました", err)
	}
	defer tx.Rollback(ctx) //nolint:errcheck

	qtx := r.queries.WithTx(tx)

	row, err := qtx.CreateSake(ctx, sqlc.CreateSakeParams{
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

	sakeID := uuid.UUID(row.ID.Bytes)
	if err := r.syncTags(ctx, tx, sakeID, input.TagNames); err != nil {
		return nil, err
	}
	if err := r.syncImages(ctx, tx, sakeID, input.ImageKeys); err != nil {
		return nil, err
	}
	if err := tx.Commit(ctx); err != nil {
		logger.LogDatabaseError(ctx, "COMMIT", "sakes", err, nil)
		return nil, apperror.DatabaseError("トランザクションのコミットに失敗しました", err)
	}

	return &entity.SakeListItem{
		ID:           sakeID,
		Category:     entity.SakeCategory(row.Category),
		Name:         row.Name,
		ImagePreview: firstImageKey(input.ImageKeys),
	}, nil
}

// Update は現行DB schemaに在庫データを更新する。
func (r *sakeRepositoryImpl) Update(ctx context.Context, input repository.UpdateSakeInput) (*entity.SakeListItem, error) {
	defer logger.TraceMethodAuto(ctx, input)()

	tx, err := r.db.Begin(ctx)
	if err != nil {
		logger.LogDatabaseError(ctx, "BEGIN", "sakes", err, nil)
		return nil, apperror.DatabaseError("トランザクション開始に失敗しました", err)
	}
	defer tx.Rollback(ctx) //nolint:errcheck

	qtx := r.queries.WithTx(tx)

	row, err := qtx.UpdateSake(ctx, sqlc.UpdateSakeParams{
		ID:                pgtype.UUID{Bytes: [16]byte(input.ID), Valid: true},
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
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperror.NotFoundError("酒が見つかりません").WithDetails("sake_id", input.ID.String())
		}

		logger.LogDatabaseError(ctx, "UPDATE", "sakes", err, map[string]interface{}{"input": input})
		return nil, apperror.DatabaseError("酒の更新に失敗しました", err)
	}

	if err := r.syncTags(ctx, tx, input.ID, input.TagNames); err != nil {
		return nil, err
	}
	if err := r.syncImages(ctx, tx, input.ID, input.ImageKeys); err != nil {
		return nil, err
	}
	if err := tx.Commit(ctx); err != nil {
		logger.LogDatabaseError(ctx, "COMMIT", "sakes", err, nil)
		return nil, apperror.DatabaseError("トランザクションのコミットに失敗しました", err)
	}

	return &entity.SakeListItem{
		ID:           uuid.UUID(row.ID.Bytes),
		Category:     entity.SakeCategory(row.Category),
		Name:         row.Name,
		ImagePreview: firstImageKey(input.ImageKeys),
	}, nil
}

// Delete はUUIDで指定された在庫データを削除する。
func (r *sakeRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	defer logger.TraceMethodAuto(ctx, id)()

	rowsAffected, err := r.queries.DeleteSake(ctx, pgtype.UUID{Bytes: [16]byte(id), Valid: true})
	if err != nil {
		logger.LogDatabaseError(ctx, "DELETE", "sakes", err, map[string]interface{}{"sake_id": id.String()})
		return apperror.DatabaseError("酒の削除に失敗しました", err)
	}
	if rowsAffected == 0 {
		return apperror.NotFoundError("酒が見つかりません").WithDetails("sake_id", id.String())
	}

	return nil
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

// syncTags はお酒のタグ関連を現在の入力で置き換える。
func (r *sakeRepositoryImpl) syncTags(ctx context.Context, tx pgx.Tx, sakeID uuid.UUID, tagNames []string) error {
	if _, err := tx.Exec(ctx, `DELETE FROM sake_tag_links WHERE sake_id = $1`, sakeID); err != nil {
		logger.LogDatabaseError(ctx, "DELETE", "sake_tag_links", err, map[string]interface{}{"sake_id": sakeID.String()})
		return apperror.DatabaseError("タグの削除に失敗しました", err)
	}

	for _, tagName := range tagNames {
		var tagID uuid.UUID
		if err := tx.QueryRow(ctx, `
			INSERT INTO sake_tags (tag)
			VALUES ($1)
			ON CONFLICT (tag) DO UPDATE SET tag = EXCLUDED.tag
			RETURNING id
		`, tagName).Scan(&tagID); err != nil {
			logger.LogDatabaseError(ctx, "UPSERT", "sake_tags", err, map[string]interface{}{"tag": tagName})
			return apperror.DatabaseError("タグの登録に失敗しました", err)
		}

		if _, err := tx.Exec(ctx, `
			INSERT INTO sake_tag_links (sake_id, sake_tag_id)
			VALUES ($1, $2)
		`, sakeID, tagID); err != nil {
			logger.LogDatabaseError(ctx, "INSERT", "sake_tag_links", err, map[string]interface{}{
				"sake_id": sakeID.String(),
				"tag_id":  tagID.String(),
			})
			return apperror.DatabaseError("タグの紐付けに失敗しました", err)
		}
	}

	return nil
}

// syncImages はお酒画像を現在の入力で置き換える。
func (r *sakeRepositoryImpl) syncImages(ctx context.Context, tx pgx.Tx, sakeID uuid.UUID, imageKeys []string) error {
	if _, err := tx.Exec(ctx, `DELETE FROM sake_images WHERE sake_id = $1`, sakeID); err != nil {
		logger.LogDatabaseError(ctx, "DELETE", "sake_images", err, map[string]interface{}{"sake_id": sakeID.String()})
		return apperror.DatabaseError("画像の削除に失敗しました", err)
	}

	for idx, imageKey := range imageKeys {
		if _, err := tx.Exec(ctx, `
			INSERT INTO sake_images (sake_id, image_key, sort_order)
			VALUES ($1, $2, $3)
		`, sakeID, imageKey, idx); err != nil {
			logger.LogDatabaseError(ctx, "INSERT", "sake_images", err, map[string]interface{}{
				"sake_id":    sakeID.String(),
				"image_key":  imageKey,
				"sort_order": idx,
			})
			return apperror.DatabaseError("画像の登録に失敗しました", err)
		}
	}

	return nil
}

// firstImageKey は先頭画像キーを一覧表示用に返す。
func firstImageKey(imageKeys []string) string {
	if len(imageKeys) == 0 {
		return ""
	}
	return imageKeys[0]
}
