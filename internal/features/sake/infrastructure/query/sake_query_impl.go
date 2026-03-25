package query

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/sake-kasu/sake-hack-backend/internal/apperror"
	appQuery "github.com/sake-kasu/sake-hack-backend/internal/features/sake/application/query"
	"github.com/sake-kasu/sake-hack-backend/internal/features/sake/domain/entity"
	"github.com/sake-kasu/sake-hack-backend/internal/logger"
)

type sakeQueryImpl struct{}

// NewSakeQuery SakeQueryの実装を作成
func NewSakeQuery(_ *pgxpool.Pool) appQuery.SakeQuery {
	return &sakeQueryImpl{}
}

// List はDB/OASの再設計に伴う再実装まで一時的に未提供。
func (q *sakeQueryImpl) List(ctx context.Context, filter appQuery.ListSakesFilter) ([]entity.SakeListItem, entity.Pagination, error) {
	defer logger.TraceMethodAuto(ctx, filter)()

	return nil, entity.Pagination{}, apperror.InternalServerError("sake query is temporarily disabled")
}

// GetDetail はDB/OASの再設計に伴う再実装まで一時的に未提供。
func (q *sakeQueryImpl) GetDetail(ctx context.Context, id int32) (*entity.SakeDetail, error) {
	defer logger.TraceMethodAuto(ctx, id)()

	return nil, apperror.InternalServerError("sake detail query is temporarily disabled")
}
