package query

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/sake-kasu/sake-hack-backend/internal/apperror"
	"github.com/sake-kasu/sake-hack-backend/internal/database/sqlc"
	appQuery "github.com/sake-kasu/sake-hack-backend/internal/features/sake/application/query"
	"github.com/sake-kasu/sake-hack-backend/internal/features/sake/domain/entity"
	"github.com/sake-kasu/sake-hack-backend/internal/logger"
)

type sakeQueryImpl struct {
	queries *sqlc.Queries
}

// NewSakeQuery SakeQueryの実装を作成
func NewSakeQuery(db *pgxpool.Pool) appQuery.SakeQuery {
	return &sakeQueryImpl{
		queries: sqlc.New(db),
	}
}

// List フィルター条件に基づいて酒一覧を取得
func (q *sakeQueryImpl) List(ctx context.Context, filter appQuery.ListSakesFilter) ([]entity.SakeListItem, entity.Pagination, error) {
	defer logger.TraceMethodAuto(ctx, filter)()

	total, err := q.queries.CountSakes(ctx, sqlc.CountSakesParams{
		KindID:    filter.KindID,
		BreweryID: filter.BreweryID,
	})
	if err != nil {
		logger.LogDatabaseError(ctx, "SELECT", "sakes", err, map[string]interface{}{"filter": filter})
		return nil, entity.Pagination{}, apperror.DatabaseError("酒の件数取得に失敗しました", err)
	}

	rows, err := q.queries.ListSakes(ctx, sqlc.ListSakesParams{
		Limit:     filter.Limit,
		Offset:    filter.Offset,
		KindID:    filter.KindID,
		BreweryID: filter.BreweryID,
	})
	if err != nil {
		logger.LogDatabaseError(ctx, "SELECT", "sakes", err, map[string]interface{}{"filter": filter})
		return nil, entity.Pagination{}, apperror.DatabaseError("酒一覧の取得に失敗しました", err)
	}

	items := make([]entity.SakeListItem, 0, len(rows))
	for _, row := range rows {
		imagePreview := ""
		if row.ImageUrl != nil {
			imagePreview = *row.ImageUrl
		}
		items = append(items, entity.SakeListItem{
			ID:           row.ID,
			Category:     entity.SakeCategory(row.Category),
			Name:         row.Name,
			ImagePreview: imagePreview,
		})
	}

	pagination := entity.Pagination{
		Total:  total,
		Offset: filter.Offset,
		Limit:  filter.Limit,
	}

	return items, pagination, nil
}

// GetDetail 酒の詳細を取得
func (q *sakeQueryImpl) GetDetail(ctx context.Context, id int32) (*entity.SakeDetail, error) {
	defer logger.TraceMethodAuto(ctx, id)()

	row, err := q.queries.GetSakeDetailByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperror.NotFoundError("酒が見つかりません").WithDetails("sake_id", id)
		}
		logger.LogDatabaseError(ctx, "SELECT", "sakes", err, map[string]interface{}{"sake_id": id})
		return nil, apperror.DatabaseError("酒の詳細取得に失敗しました", err)
	}

	drinkStyleRows, err := q.queries.GetDrinkStylesBySakeID(ctx, id)
	if err != nil {
		logger.LogDatabaseError(ctx, "SELECT", "drink_styles", err, map[string]interface{}{"sake_id": id})
		return nil, apperror.DatabaseError("飲み方の取得に失敗しました", err)
	}

	drinkStyles := make([]entity.DrinkStyle, 0, len(drinkStyleRows))
	for _, ds := range drinkStyleRows {
		drinkStyles = append(drinkStyles, entity.DrinkStyle{
			ID:          ds.ID,
			Name:        ds.Name,
			Description: ds.Description,
		})
	}

	return &entity.SakeDetail{
		ID:       row.ID,
		Category: entity.SakeCategory(row.Category),
		Kind: entity.SakeKind{
			ID:   row.KindID,
			Name: row.KindName,
		},
		Brewery: entity.Brewery{
			ID:            row.BreweryID,
			Name:          row.BreweryName,
			OriginCountry: row.BreweryOriginCountry,
			OriginRegion:  row.BreweryOriginRegion,
			Latitude:      row.BreweryLatitude,
			Longitude:     row.BreweryLongitude,
		},
		Name: entity.SakeName{
			Name:     row.Name,
			Phonetic: row.Phonetic,
		},
		Abv:             row.Abv,
		PurchaseVolume:  row.PurchaseVolume,
		RemainingVolume: row.RemainingVolume,
		Memo:            row.Memo,
		DrinkStyles:     drinkStyles,
		Price:           row.Price,
		ImageUrl:        row.ImageUrl,
		CreatedAt:       row.CreatedAt.Time,
		UpdatedAt:       row.UpdatedAt.Time,
	}, nil
}
