package presentation

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/sake-kasu/sake-hack-backend/api/generated"
	"github.com/sake-kasu/sake-hack-backend/internal/apperror"
	"github.com/sake-kasu/sake-hack-backend/internal/features/sake/application/usecase"
	"github.com/sake-kasu/sake-hack-backend/internal/logger"
	"github.com/sake-kasu/sake-hack-backend/internal/utils"
)

// SakeServerImpl 酒関連のServerInterface実装
type SakeServerImpl struct {
	listSakesUC     usecase.ListSakesUsecaseInterface
	getSakeDetailUC usecase.GetSakeDetailUsecaseInterface
}

// NewSakeServerImpl コンストラクタ
func NewSakeServerImpl(
	listSakesUC usecase.ListSakesUsecaseInterface,
	getSakeDetailUC usecase.GetSakeDetailUsecaseInterface,
) *SakeServerImpl {
	return &SakeServerImpl{
		listSakesUC:     listSakesUC,
		getSakeDetailUC: getSakeDetailUC,
	}
}

// GetSakes 酒一覧取得（公開用）
// (GET /sakes)
func (s *SakeServerImpl) GetSakes(c *gin.Context, params generated.GetSakesParams) {
	ctx := c.Request.Context()
	defer logger.TraceMethodAuto(ctx, params)()

	offset := int32(0)
	if params.Offset != nil {
		offset = *params.Offset
	}

	limit := int32(20)
	if params.Limit != nil {
		limit = *params.Limit
	}

	var category *string
	if params.Category != nil {
		cat := string(*params.Category)
		category = &cat
	}

	input := usecase.ListSakesInput{
		Category: category,
		Search:   params.Q,
		Offset:   offset,
		Limit:    limit,
	}

	output, err := s.listSakesUC.Execute(ctx, input)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, toListSakesResponse(output))
}

// GetSakeDetail 酒詳細取得（公開用）
// (GET /sakes/{id})
func (s *SakeServerImpl) GetSakeDetail(c *gin.Context, id generated.SakeIDPathParameter) {
	ctx := c.Request.Context()
	defer logger.TraceMethodAuto(ctx, id)()

	// UUIDをint32に変換
	id32, err := convertUUIDToInt32(id)
	if err != nil {
		handleError(c, err)
		return
	}

	output, err := s.getSakeDetailUC.Execute(ctx, id32)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, toSakeDetail(*output.Detail))
}

// GetStocks 在庫一覧取得
// (GET /stocks)
func (s *SakeServerImpl) GetStocks(c *gin.Context, params generated.GetStocksParams) {
	ctx := c.Request.Context()
	defer logger.TraceMethodAuto(ctx, params)()

	// TODO: 在庫一覧取得の実装
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented yet"})
}

// CreateStock 在庫登録
// (POST /stocks)
func (s *SakeServerImpl) CreateStock(c *gin.Context) {
	ctx := c.Request.Context()
	defer logger.TraceMethodAuto(ctx, nil)()

	// TODO: 在庫登録の実装
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented yet"})
}

// GetStockDetail 在庫詳細取得
// (GET /stocks/{id})
func (s *SakeServerImpl) GetStockDetail(c *gin.Context, id generated.StockIDPathParameter) {
	ctx := c.Request.Context()
	defer logger.TraceMethodAuto(ctx, id)()

	// TODO: 在庫詳細取得の実装
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented yet"})
}

// UpdateStock 在庫更新
// (PATCH /stocks/{id})
func (s *SakeServerImpl) UpdateStock(c *gin.Context, id generated.StockIDPathParameter) {
	ctx := c.Request.Context()
	defer logger.TraceMethodAuto(ctx, id)()

	// TODO: 在庫更新の実装
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented yet"})
}

// DeleteStock 在庫削除
// (DELETE /stocks/{id})
func (s *SakeServerImpl) DeleteStock(c *gin.Context, id generated.StockIDPathParameter) {
	ctx := c.Request.Context()
	defer logger.TraceMethodAuto(ctx, id)()

	// TODO: 在庫削除の実装
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented yet"})
}

// convertUUIDToInt32 UUIDをint32に変換（ヘルパー関数）
func convertUUIDToInt32(id uuid.UUID) (int32, error) {
	return utils.UUIDToInt32(id)
}

// handleError エラーをHTTPレスポンスに変換
func handleError(c *gin.Context, err error) {
	if valErr, ok := err.(*apperror.ValidationError); ok {
		c.JSON(valErr.Status, generated.ErrorResponse{
			ErrorCode: generated.ErrorCode(valErr.Code),
			Message:   valErr.Message,
			Reason:    nil,
		})
		return
	}

	var appErr *apperror.AppError
	if appErr = apperror.As(err); appErr != nil {
		c.JSON(appErr.Status, generated.ErrorResponse{
			ErrorCode: generated.ErrorCode(appErr.Code),
			Message:   appErr.Message,
			Reason:    nil,
		})
		return
	}

	c.JSON(http.StatusInternalServerError, generated.ErrorResponse{
		ErrorCode: generated.ErrorCodeE9000,
		Message:   "内部エラーが発生しました",
		Reason:    nil,
	})
}
