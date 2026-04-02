package query

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/pgtype"

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
func (q *sakeQueryImpl) List(ctx context.Context, filter appQuery.ListSakesFilter) ([]entity.SakeListItem, entity.Pagination, error) {
	defer logger.TraceMethodAuto(ctx, filter)()

	total, err := q.queries.CountSakes(ctx, sqlc.NullSakeCategoryType{})
	if err != nil {
		logger.LogDatabaseError(ctx, "SELECT", "sakes", err, map[string]interface{}{"filter": filter})
		return nil, entity.Pagination{}, apperror.DatabaseError("酒の件数取得に失敗しました", err)
	}

	rows, err := q.queries.ListSakes(ctx, sqlc.ListSakesParams{
		Limit:    filter.Limit,
		Offset:   filter.Offset,
		Category: sqlc.NullSakeCategoryType{},
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

// GetDetail はUUIDで指定された酒詳細を取得する。
func (q *sakeQueryImpl) GetDetail(ctx context.Context, id uuid.UUID) (*entity.SakeDetail, error) {
	defer logger.TraceMethodAuto(ctx, id)()

	row, err := q.queries.GetSakeDetailByID(ctx, pgtype.UUID{Bytes: [16]byte(id), Valid: true})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperror.NotFoundError("酒が見つかりません").WithDetails("sake_id", id.String())
		}

		logger.LogDatabaseError(ctx, "SELECT", "sakes", err, map[string]interface{}{"sake_id": id.String()})
		return nil, apperror.DatabaseError("酒の詳細取得に失敗しました", err)
	}

	abv, err := numericToFloat32Ptr(row.AlcoholPercentage)
	if err != nil {
		return nil, apperror.DatabaseError("アルコール度数の変換に失敗しました", err)
	}
	tags, err := q.loadTags(ctx, id)
	if err != nil {
		return nil, err
	}
	images, err := q.loadImages(ctx, id)
	if err != nil {
		return nil, err
	}

	return &entity.SakeDetail{
		ID:       id,
		Category: entity.SakeCategory(row.Category),
		Name: entity.SakeName{
			Name:     row.Name,
			Phonetic: stringValue(row.Phonetic),
		},
		Abv:             abv,
		PurchaseVolume:  row.VolumeMax,
		RemainingVolume: row.VolumeRemain,
		Brewery: entity.Brewery{
			OriginRegion: row.Region,
		},
		Memo:        row.Memo,
		DrinkStyles: []entity.DrinkStyle{},
		Tags:        tags,
		Price:       row.Price,
		ImageUrl:    firstDetailImage(images),
		Images:      images,
		CreatedAt:   row.CreatedAt.Time,
		UpdatedAt:   row.UpdatedAt.Time,
	}, nil
}

// numericToFloat32Ptr は pgtype.Numeric を *float32 に変換する。
func numericToFloat32Ptr(v pgtype.Numeric) (*float32, error) {
	floatValue, err := v.Float64Value()
	if err != nil {
		return nil, err
	}
	if !floatValue.Valid {
		return nil, nil
	}

	result := float32(floatValue.Float64)
	return &result, nil
}

// stringValue は nil を空文字に正規化する。
func stringValue(v *string) string {
	if v == nil {
		return ""
	}
	return *v
}

// loadTags は酒に紐づくタグを読み込む。
func (q *sakeQueryImpl) loadTags(ctx context.Context, sakeID uuid.UUID) ([]entity.SakeTag, error) {
	rows, err := q.db.Query(ctx, `
		SELECT st.id, st.tag
		FROM sake_tags st
		INNER JOIN sake_tag_links stl ON stl.sake_tag_id = st.id
		WHERE stl.sake_id = $1
		ORDER BY st.tag ASC
	`, sakeID)
	if err != nil {
		logger.LogDatabaseError(ctx, "SELECT", "sake_tags", err, map[string]interface{}{"sake_id": sakeID.String()})
		return nil, apperror.DatabaseError("タグの取得に失敗しました", err)
	}
	defer rows.Close()

	tags := []entity.SakeTag{}
	for rows.Next() {
		var id uuid.UUID
		var name string
		if err := rows.Scan(&id, &name); err != nil {
			return nil, apperror.DatabaseError("タグの取得に失敗しました", err)
		}
		tags = append(tags, entity.SakeTag{ID: id, Name: name})
	}
	if err := rows.Err(); err != nil {
		return nil, apperror.DatabaseError("タグの取得に失敗しました", err)
	}
	return tags, nil
}

// loadImages は酒に紐づく画像を表示順で読み込む。
func (q *sakeQueryImpl) loadImages(ctx context.Context, sakeID uuid.UUID) ([]entity.SakeImage, error) {
	rows, err := q.db.Query(ctx, `
		SELECT id, image_key, sort_order
		FROM sake_images
		WHERE sake_id = $1
		ORDER BY sort_order ASC
	`, sakeID)
	if err != nil {
		logger.LogDatabaseError(ctx, "SELECT", "sake_images", err, map[string]interface{}{"sake_id": sakeID.String()})
		return nil, apperror.DatabaseError("画像の取得に失敗しました", err)
	}
	defer rows.Close()

	images := []entity.SakeImage{}
	for rows.Next() {
		var image entity.SakeImage
		if err := rows.Scan(&image.ID, &image.ImageKey, &image.SortOrder); err != nil {
			return nil, apperror.DatabaseError("画像の取得に失敗しました", err)
		}
		images = append(images, image)
	}
	if err := rows.Err(); err != nil {
		return nil, apperror.DatabaseError("画像の取得に失敗しました", err)
	}
	return images, nil
}

// firstDetailImage は先頭画像を detail の後方互換フィールドに入れる。
func firstDetailImage(images []entity.SakeImage) *string {
	if len(images) == 0 {
		return nil
	}
	result := images[0].ImageKey
	return &result
}
