package query

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/sake-kasu/sake-hack-backend/internal/apperror"
	"github.com/sake-kasu/sake-hack-backend/internal/database/sqlc"
	appQuery "github.com/sake-kasu/sake-hack-backend/internal/features/sake/application/query"
	"github.com/sake-kasu/sake-hack-backend/internal/features/sake/domain/entity"
	"github.com/sake-kasu/sake-hack-backend/internal/logger"
)

type sakeQueryImpl struct {
	db      *pgxpool.Pool
	queries *sqlc.Queries
}

// NewSakeQuery SakeQueryの実装を作成
func NewSakeQuery(db *pgxpool.Pool) appQuery.SakeQuery {
	return &sakeQueryImpl{
		db:      db,
		queries: sqlc.New(db),
	}
}

// List フィルター条件に基づいて酒一覧を取得
func (q *sakeQueryImpl) List(ctx context.Context, filter appQuery.ListSakesFilter) ([]entity.SakeDetail, entity.Pagination, error) {
	defer logger.TraceMethodAuto(ctx, filter)()

	// カウント用のクエリを手動で実行
	countQuery := `
		SELECT COUNT(*) AS total
		FROM sakes s
		WHERE
			($1::text IS NULL OR s.category::text = $1)
			AND ($2::text IS NULL OR s.name ILIKE '%' || $2 || '%')
	`
	var total int64
	err := q.db.QueryRow(ctx, countQuery, filter.Category, filter.Search).Scan(&total)
	if err != nil {
		logger.LogDatabaseError(ctx, "SELECT", "sakes", err, map[string]interface{}{"filter": filter})
		return nil, entity.Pagination{}, apperror.DatabaseError("酒の件数取得に失敗しました", err)
	}

	// リスト取得用のクエリを手動で実行
	listQuery := `
		SELECT
			s.id,
			s.category,
			s.name,
			s.image_url,
			s.abv,
			s.memo,
			s.created_at,
			s.updated_at,
			sk.name AS kind_name,
			b.origin_region AS brewery_region
		FROM sakes s
		INNER JOIN sake_kinds sk ON s.kind_id = sk.id
		INNER JOIN breweries b ON s.brewery_id = b.id
		WHERE
			($1::text IS NULL OR s.category::text = $1)
			AND ($2::text IS NULL OR s.name ILIKE '%' || $2 || '%')
		ORDER BY s.created_at DESC
		LIMIT $3 OFFSET $4
	`
	rows, err := q.db.Query(ctx, listQuery, filter.Category, filter.Search, filter.Limit, filter.Offset)
	if err != nil {
		logger.LogDatabaseError(ctx, "SELECT", "sakes", err, map[string]interface{}{"filter": filter})
		return nil, entity.Pagination{}, apperror.DatabaseError("酒一覧の取得に失敗しました", err)
	}
	defer rows.Close()

	items := make([]entity.SakeDetail, 0)
	for rows.Next() {
		var item struct {
			ID            int32
			Category      string
			Name          string
			ImageURL      *string
			Abv           float32
			Memo          *string
			CreatedAt     pgtype.Timestamptz
			UpdatedAt     pgtype.Timestamptz
			KindName      string
			BreweryRegion *string
		}
		
		err := rows.Scan(
			&item.ID,
			&item.Category,
			&item.Name,
			&item.ImageURL,
			&item.Abv,
			&item.Memo,
			&item.CreatedAt,
			&item.UpdatedAt,
			&item.KindName,
			&item.BreweryRegion,
		)
		if err != nil {
			logger.LogDatabaseError(ctx, "SCAN", "sakes", err, nil)
			return nil, entity.Pagination{}, apperror.DatabaseError("酒一覧のスキャンに失敗しました", err)
		}

		items = append(items, entity.SakeDetail{
			ID:       item.ID,
			Category: entity.SakeCategory(item.Category),
			Kind: entity.SakeKind{
				ID:   0,
				Name: item.KindName,
			},
			Brewery: entity.Brewery{
				ID:            0,
				Name:          "",
				OriginCountry: "",
				OriginRegion:  item.BreweryRegion,
			},
			Name: entity.SakeName{
				Name:     item.Name,
				Phonetic: "",
			},
			Abv:             item.Abv,
			PurchaseVolume:  0,
			RemainingVolume: 0,
			Memo:            item.Memo,
			DrinkStyles:     []entity.DrinkStyle{},
			Price:           0,
			ImageUrl:        item.ImageURL,
			CreatedAt:       item.CreatedAt.Time,
			UpdatedAt:       item.UpdatedAt.Time,
		})
	}

	if err := rows.Err(); err != nil {
		logger.LogDatabaseError(ctx, "ROWS", "sakes", err, nil)
		return nil, entity.Pagination{}, apperror.DatabaseError("酒一覧の取得に失敗しました", err)
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
