package presentation

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/sake-kasu/sake-hack-backend/api/generated"
	"github.com/sake-kasu/sake-hack-backend/internal/apperror"
	"github.com/sake-kasu/sake-hack-backend/internal/features/sake/application/usecase"
	"github.com/sake-kasu/sake-hack-backend/internal/logger"
)

// SakeServerImpl 酒関連のServerInterface実装
type SakeServerImpl struct {
	listSakesUC     usecase.ListSakesUsecaseInterface
	getSakeDetailUC usecase.GetSakeDetailUsecaseInterface
	createStockUC   usecase.CreateStockUsecaseInterface
	updateStockUC   usecase.UpdateStockUsecaseInterface
	deleteStockUC   usecase.DeleteStockUsecaseInterface
}

// NewSakeServerImpl コンストラクタ
func NewSakeServerImpl(
	listSakesUC usecase.ListSakesUsecaseInterface,
	getSakeDetailUC usecase.GetSakeDetailUsecaseInterface,
	createStockUC usecase.CreateStockUsecaseInterface,
	updateStockUC usecase.UpdateStockUsecaseInterface,
	deleteStockUC usecase.DeleteStockUsecaseInterface,
) *SakeServerImpl {
	return &SakeServerImpl{
		listSakesUC:     listSakesUC,
		getSakeDetailUC: getSakeDetailUC,
		createStockUC:   createStockUC,
		updateStockUC:   updateStockUC,
		deleteStockUC:   deleteStockUC,
	}
}

// HealthCheck ヘルスチェック
// (GET /health)
func (s *SakeServerImpl) HealthCheck(c *gin.Context) {
	c.String(http.StatusOK, "ok")
}

// CreateSakeLike いいね登録
// (POST /sakes/{sakeId}/likes)
func (s *SakeServerImpl) CreateSakeLike(c *gin.Context, sakeId generated.SakeId, params generated.CreateSakeLikeParams) {
	c.JSON(http.StatusNotImplemented, notImplementedResponse("like API is temporarily disabled"))
}

// DeleteSakeLike いいね解除
// (DELETE /sakes/{sakeId}/likes)
func (s *SakeServerImpl) DeleteSakeLike(c *gin.Context, sakeId generated.SakeId, params generated.DeleteSakeLikeParams) {
	c.JSON(http.StatusNotImplemented, notImplementedResponse("like API is temporarily disabled"))
}

// ListSakes 酒一覧取得
// (GET /sakes)
func (s *SakeServerImpl) ListSakes(c *gin.Context, params generated.ListSakesParams) {
	ctx := c.Request.Context()
	defer logger.TraceMethodAuto(ctx, params)()

	if err := validateListParams(int32(params.Offset), int32(params.Limit)); err != nil {
		handleError(c, err)
		return
	}

	output, err := s.listSakesUC.Execute(ctx, toListSakesInput(int32(params.Offset), int32(params.Limit)))
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, toListSakesResponse(output))
}

// GetSakeDetail 酒詳細取得
// (GET /sakes/{sakeId})
func (s *SakeServerImpl) GetSakeDetail(c *gin.Context, sakeId generated.SakeId) {
	ctx := c.Request.Context()
	defer logger.TraceMethodAuto(ctx, sakeId)()

	output, err := s.getSakeDetailUC.Execute(ctx, uuid.UUID(sakeId))
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, toSakeDetailResponse(output.Detail))
}

// ListStocks 在庫一覧取得
// (GET /stocks)
func (s *SakeServerImpl) ListStocks(c *gin.Context, params generated.ListStocksParams) {
	ctx := c.Request.Context()
	defer logger.TraceMethodAuto(ctx, params)()

	if err := validateListParams(int32(params.Offset), int32(params.Limit)); err != nil {
		handleError(c, err)
		return
	}

	output, err := s.listSakesUC.Execute(ctx, toListSakesInput(int32(params.Offset), int32(params.Limit)))
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, toStockListResponse(output))
}

// CreateStock 在庫登録
// (POST /stocks)
func (s *SakeServerImpl) CreateStock(c *gin.Context) {
	ctx := c.Request.Context()
	defer logger.TraceMethodAuto(ctx, nil)()

	var req generated.CreateStockRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handleError(c, apperror.BadRequestError("リクエストボディの解析に失敗しました"))
		return
	}

	if err := validateCreateStockRequest(req); err != nil {
		handleError(c, err)
		return
	}

	output, err := s.createStockUC.Execute(ctx, toCreateStockInput(req))
	if err != nil {
		handleError(c, err)
		return
	}

	detailOutput, err := s.getSakeDetailUC.Execute(ctx, output.Sake.ID)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, toStockDetailResponse(detailOutput.Detail))
}

// DeleteStock 在庫削除
// (DELETE /stocks/{sakeId})
func (s *SakeServerImpl) DeleteStock(c *gin.Context, sakeId generated.SakeId) {
	c.JSON(http.StatusNotImplemented, notImplementedResponse("stock delete API is temporarily disabled"))
}

// GetStockDetail 在庫詳細取得
// (GET /stocks/{sakeId})
func (s *SakeServerImpl) GetStockDetail(c *gin.Context, sakeId generated.SakeId) {
	ctx := c.Request.Context()
	defer logger.TraceMethodAuto(ctx, sakeId)()

	output, err := s.getSakeDetailUC.Execute(ctx, uuid.UUID(sakeId))
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, toStockDetailResponse(output.Detail))
}

// UpdateStock 在庫更新
// (PUT /stocks/{sakeId})
func (s *SakeServerImpl) UpdateStock(c *gin.Context, sakeId generated.SakeId) {
	ctx := c.Request.Context()
	defer logger.TraceMethodAuto(ctx, sakeId)()

	var req generated.CreateStockRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handleError(c, apperror.BadRequestError("リクエストボディの解析に失敗しました"))
		return
	}

	if err := validateCreateStockRequest(req); err != nil {
		handleError(c, err)
		return
	}

	if _, err := s.updateStockUC.Execute(ctx, toUpdateStockInput(uuid.UUID(sakeId), req)); err != nil {
		handleError(c, err)
		return
	}

	detailOutput, err := s.getSakeDetailUC.Execute(ctx, uuid.UUID(sakeId))
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, toStockDetailResponse(detailOutput.Detail))
}

// toListSakesInput は一覧APIの入力をユースケース入力に変換する。
func toListSakesInput(offset, limit int32) usecase.ListSakesInput {
	return usecase.ListSakesInput{
		KindID:    nil,
		BreweryID: nil,
		Offset:    offset,
		Limit:     limit,
	}
}

// notImplementedResponse は未実装API向けの共通エラーレスポンスを返す。
func notImplementedResponse(message string) generated.ErrorResponse {
	return generated.ErrorResponse{
		Data: nil,
		Errors: []generated.APIError{
			{
				Code:    "NOT_IMPLEMENTED",
				Message: message,
			},
		},
	}
}

// handleError エラーをHTTPレスポンスに変換
func handleError(c *gin.Context, err error) {
	if valErr, ok := err.(*apperror.ValidationError); ok {
		errors := make([]generated.APIError, 0, len(valErr.Fields))
		for field, msg := range valErr.Fields {
			errors = append(errors, generated.APIError{
				Code:    valErr.Code,
				Message: field + ": " + msg,
			})
		}
		c.JSON(valErr.Status, generated.ErrorResponse{
			Data:   nil,
			Errors: errors,
		})
		return
	}

	var appErr *apperror.AppError
	if appErr = apperror.As(err); appErr != nil {
		errors := []generated.APIError{
			{
				Code:    appErr.Code,
				Message: appErr.Message,
			},
		}
		c.JSON(appErr.Status, generated.ErrorResponse{
			Data:   nil,
			Errors: errors,
		})
		return
	}

	errors := []generated.APIError{
		{
			Code:    apperror.ErrCodeInternalError,
			Message: "内部エラーが発生しました",
		},
	}
	c.JSON(http.StatusInternalServerError, generated.ErrorResponse{
		Data:   nil,
		Errors: errors,
	})
}
