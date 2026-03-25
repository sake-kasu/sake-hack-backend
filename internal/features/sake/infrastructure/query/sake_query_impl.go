package query

import (
	"context"

	"github.com/google/uuid"
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

	total, err := q.queries.CountSakes(ctx, sqlc.NullSakeCategory{})
	if err != nil {
		logger.LogDatabaseError(ctx, "SELECT", "sakes", err, map[string]interface{}{"filter": filter})
		return nil, entity.Pagination{}, apperror.DatabaseError("酒の件数取得に失敗しました", err)
	}

	rows, err := q.queries.ListSakes(ctx, sqlc.ListSakesParams{
		Limit:    filter.Limit,
		Offset:   filter.Offset,
		Category: sqlc.NullSakeCategory{},
	})
	if err != nil {
		logger.LogDatabaseError(ctx, "SELECT", "sakes", err, map[string]interface{}{"filter": filter})
		return nil, entity.Pagination{}, apperror.DatabaseError("酒一覧の取得に失敗しました", err)
	}

	items := make([]entity.SakeListItem, 0, len(rows))
	for _, row := range rows {
		var id uuid.UUID
		if row.ID.Valid {
			id = uuid.UUID(row.ID.Bytes)
		}

		imagePreview := ""
		if row.ImageKey != nil {
			imagePreview = *row.ImageKey
		}

		items = append(items, entity.SakeListItem{
			ID:           id,
			Category:     entity.SakeCategory(row.Category),
			Name:         row.Name,
			ImagePreview: imagePreview,
		})
	}

	return items, entity.Pagination{
		Total:  total,
		Offset: filter.Offset,
		Limit:  filter.Limit,
	}, nil
}

// GetDetail はDB/OASの再設計に伴う再実装まで一時的に未提供。
func (q *sakeQueryImpl) GetDetail(ctx context.Context, id int32) (*entity.SakeDetail, error) {
	defer logger.TraceMethodAuto(ctx, id)()

	return nil, apperror.InternalServerError("sake detail query is temporarily disabled")
}
