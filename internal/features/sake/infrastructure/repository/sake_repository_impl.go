package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

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

// Create 酒を新規作成する
func (r *sakeRepositoryImpl) Create(ctx context.Context, input repository.CreateSakeInput) (*entity.SakeListItem, error) {
	defer logger.TraceMethodAuto(ctx, input)()

	tx, err := r.db.Begin(ctx)
	if err != nil {
		logger.LogDatabaseError(ctx, "BEGIN", "sakes", err, nil)
		return nil, apperror.DatabaseError("トランザクション開始に失敗しました", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	qtx := r.queries.WithTx(tx)

	breweryID, err := r.upsertBrewery(ctx, qtx, input.Brewery)
	if err != nil {
		return nil, err
	}

	kindID, err := r.upsertSakeKind(ctx, qtx, input.Kind)
	if err != nil {
		return nil, err
	}

	row, err := qtx.CreateSake(ctx, sqlc.CreateSakeParams{
		Category:        sqlc.SakeCategory(input.Category),
		KindID:          kindID,
		BreweryID:       breweryID,
		Name:            input.Name.Name,
		Phonetic:        input.Name.Phonetic,
		Abv:             input.Abv,
		PurchaseVolume:  input.PurchaseVolume,
		RemainingVolume: input.RemainingVolume,
		Memo:            input.Memo,
		Price:           input.Price,
		ObjectKey:       input.ObjectKey,
	})
	if err != nil {
		logger.LogDatabaseError(ctx, "INSERT", "sakes", err, map[string]interface{}{"input": input})
		return nil, apperror.DatabaseError("酒の登録に失敗しました", err)
	}

	if err := r.syncDrinkStyles(ctx, qtx, row.ID, input.DrinkStyles); err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		logger.LogDatabaseError(ctx, "COMMIT", "sakes", err, nil)
		return nil, apperror.DatabaseError("トランザクションのコミットに失敗しました", err)
	}

	return &entity.SakeListItem{
		ID:        row.ID,
		Category:  entity.SakeCategory(row.Category),
		Name:      row.Name,
		ObjectKey: row.ObjectKey,
	}, nil
}

// Update 酒を更新する
func (r *sakeRepositoryImpl) Update(ctx context.Context, input repository.UpdateSakeInput) (*entity.SakeListItem, error) {
	defer logger.TraceMethodAuto(ctx, input)()

	tx, err := r.db.Begin(ctx)
	if err != nil {
		logger.LogDatabaseError(ctx, "BEGIN", "sakes", err, nil)
		return nil, apperror.DatabaseError("トランザクション開始に失敗しました", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	qtx := r.queries.WithTx(tx)

	breweryID, err := r.upsertBrewery(ctx, qtx, input.Brewery)
	if err != nil {
		return nil, err
	}

	kindID, err := r.upsertSakeKind(ctx, qtx, input.Kind)
	if err != nil {
		return nil, err
	}

	row, err := qtx.UpdateSake(ctx, sqlc.UpdateSakeParams{
		ID:              input.ID,
		Category:        sqlc.SakeCategory(input.Category),
		KindID:          kindID,
		BreweryID:       breweryID,
		Name:            input.Name.Name,
		Phonetic:        input.Name.Phonetic,
		Abv:             input.Abv,
		PurchaseVolume:  input.PurchaseVolume,
		RemainingVolume: input.RemainingVolume,
		Memo:            input.Memo,
		Price:           input.Price,
		ObjectKey:       input.ObjectKey,
	})
	if err != nil {
		logger.LogDatabaseError(ctx, "UPDATE", "sakes", err, map[string]interface{}{"sake_id": input.ID})
		return nil, apperror.DatabaseError("酒の更新に失敗しました", err)
	}

	if err := r.syncDrinkStyles(ctx, qtx, row.ID, input.DrinkStyles); err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		logger.LogDatabaseError(ctx, "COMMIT", "sakes", err, nil)
		return nil, apperror.DatabaseError("トランザクションのコミットに失敗しました", err)
	}

	return &entity.SakeListItem{
		ID:        row.ID,
		Category:  entity.SakeCategory(row.Category),
		Name:      row.Name,
		ObjectKey: row.ObjectKey,
	}, nil
}

// UpdateObjectKey 酒のオブジェクトキーを更新する
func (r *sakeRepositoryImpl) UpdateObjectKey(ctx context.Context, id int32, objectKey string) error {
	defer logger.TraceMethodAuto(ctx, map[string]interface{}{"sake_id": id, "object_key": objectKey})()

	rowsAffected, err := r.queries.UpdateSakeObjectKey(ctx, sqlc.UpdateSakeObjectKeyParams{
		ID:        id,
		ObjectKey: &objectKey,
	})
	if err != nil {
		logger.LogDatabaseError(ctx, "UPDATE", "sakes", err, map[string]interface{}{"sake_id": id})
		return apperror.DatabaseError("オブジェクトキーの更新に失敗しました", err)
	}
	if rowsAffected == 0 {
		return apperror.NotFoundError("酒が見つかりません").WithDetails("sake_id", id)
	}
	return nil
}

// Delete 酒を削除する
func (r *sakeRepositoryImpl) Delete(ctx context.Context, id int32) error {
	defer logger.TraceMethodAuto(ctx, id)()

	rowsAffected, err := r.queries.DeleteSake(ctx, id)
	if err != nil {
		logger.LogDatabaseError(ctx, "DELETE", "sakes", err, map[string]interface{}{"sake_id": id})
		return apperror.DatabaseError("酒の削除に失敗しました", err)
	}
	if rowsAffected == 0 {
		return apperror.NotFoundError("酒が見つかりません").WithDetails("sake_id", id)
	}

	return nil
}

// upsertBrewery 酒造をUPSERTしてIDを返す
func (r *sakeRepositoryImpl) upsertBrewery(ctx context.Context, qtx *sqlc.Queries, brewery entity.Brewery) (int32, error) {
	id, err := qtx.UpsertBrewery(ctx, sqlc.UpsertBreweryParams{
		Name:          brewery.Name,
		OriginCountry: brewery.OriginCountry,
		OriginRegion:  brewery.OriginRegion,
		Latitude:      brewery.Latitude,
		Longitude:     brewery.Longitude,
	})
	if err != nil {
		logger.LogDatabaseError(ctx, "UPSERT", "breweries", err, map[string]interface{}{"brewery": brewery.Name})
		return 0, apperror.DatabaseError("酒造の登録に失敗しました", err)
	}
	return id, nil
}

// upsertSakeKind 酒の種類をUPSERTしてIDを返す
func (r *sakeRepositoryImpl) upsertSakeKind(ctx context.Context, qtx *sqlc.Queries, kind entity.SakeKind) (int32, error) {
	id, err := qtx.UpsertSakeKind(ctx, kind.Name)
	if err != nil {
		logger.LogDatabaseError(ctx, "UPSERT", "sake_kinds", err, map[string]interface{}{"kind": kind.Name})
		return 0, apperror.DatabaseError("酒の種類の登録に失敗しました", err)
	}
	return id, nil
}

// syncDrinkStyles 飲み方を同期する(既存を削除して再挿入)
func (r *sakeRepositoryImpl) syncDrinkStyles(ctx context.Context, qtx *sqlc.Queries, sakeID int32, drinkStyles []entity.DrinkStyle) error {
	if err := qtx.DeleteSakeDrinkStyles(ctx, sakeID); err != nil {
		logger.LogDatabaseError(ctx, "DELETE", "sake_drink_styles", err, map[string]interface{}{"sake_id": sakeID})
		return apperror.DatabaseError("飲み方の削除に失敗しました", err)
	}

	for _, ds := range drinkStyles {
		dsID, err := qtx.UpsertDrinkStyle(ctx, sqlc.UpsertDrinkStyleParams{
			Name:        ds.Name,
			Description: ds.Description,
		})
		if err != nil {
			logger.LogDatabaseError(ctx, "UPSERT", "drink_styles", err, map[string]interface{}{"drink_style": ds.Name})
			return apperror.DatabaseError("飲み方の登録に失敗しました", err)
		}

		if err := qtx.InsertSakeDrinkStyle(ctx, sqlc.InsertSakeDrinkStyleParams{
			SakeID:       sakeID,
			DrinkStyleID: dsID,
		}); err != nil {
			logger.LogDatabaseError(ctx, "INSERT", "sake_drink_styles", err, map[string]interface{}{
				"sake_id":        sakeID,
				"drink_style_id": dsID,
			})
			return apperror.DatabaseError("酒と飲み方の紐付けに失敗しました", err)
		}
	}

	return nil
}
