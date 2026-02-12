package presentation

import (
	"net/http"

	"github.com/gin-gonic/gin"

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

// ListSakes 酒一覧取得
// (GET /sakes)
func (s *SakeServerImpl) ListSakes(c *gin.Context, params generated.ListSakesParams) {
	ctx := c.Request.Context()
	defer logger.TraceMethodAuto(ctx, params)()

	if err := validateListParams(params.Offset, params.Limit, params.TypeId, params.BreweryId); err != nil {
		handleError(c, err)
		return
	}

	output, err := s.listSakesUC.Execute(ctx, toListSakesInput(params.Offset, params.Limit, params.TypeId, params.BreweryId))
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, toListSakesResponse(output))
}

// GetSakeDetail 酒詳細取得
// (GET /sakes/{id})
func (s *SakeServerImpl) GetSakeDetail(c *gin.Context, id int32) {
	ctx := c.Request.Context()
	defer logger.TraceMethodAuto(ctx, id)()

	output, err := s.getSakeDetailUC.Execute(ctx, id)
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

	if err := validateListParams(params.Offset, params.Limit, params.TypeId, params.BreweryId); err != nil {
		handleError(c, err)
		return
	}

	output, err := s.listSakesUC.Execute(ctx, toListSakesInput(params.Offset, params.Limit, params.TypeId, params.BreweryId))
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, toListSakesResponse(output))
}

// CreateStock 在庫登録
// (POST /stocks)
func (s *SakeServerImpl) CreateStock(c *gin.Context) {
	ctx := c.Request.Context()
	defer logger.TraceMethodAuto(ctx, nil)()

	var req generated.CreateSakeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handleError(c, apperror.BadRequestError("リクエストボディの解析に失敗しました"))
		return
	}

	if err := validateCreateSakeRequest(req); err != nil {
		handleError(c, err)
		return
	}

	input := toCreateStockInput(req)
	output, err := s.createStockUC.Execute(ctx, input)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, toCreateSakeResponse(output.Sake))
}

// GetStockDetail 在庫詳細取得
// (GET /stocks/{id})
func (s *SakeServerImpl) GetStockDetail(c *gin.Context, id int32) {
	ctx := c.Request.Context()
	defer logger.TraceMethodAuto(ctx, id)()

	output, err := s.getSakeDetailUC.Execute(ctx, id)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, toSakeDetailResponse(output.Detail))
}

// UpdateStock 在庫更新
// (PUT /stocks/{id})
func (s *SakeServerImpl) UpdateStock(c *gin.Context, id int32) {
	ctx := c.Request.Context()
	defer logger.TraceMethodAuto(ctx, id)()

	var req generated.UpdateSakeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handleError(c, apperror.BadRequestError("リクエストボディの解析に失敗しました"))
		return
	}

	if err := validateUpdateSakeRequest(req); err != nil {
		handleError(c, err)
		return
	}

	input := toUpdateStockInput(id, req)
	if _, err := s.updateStockUC.Execute(ctx, input); err != nil {
		handleError(c, err)
		return
	}

	detailOutput, err := s.getSakeDetailUC.Execute(ctx, id)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, toSakeDetailResponse(detailOutput.Detail))
}

// DeleteStock 在庫削除
// (DELETE /stocks/{id})
func (s *SakeServerImpl) DeleteStock(c *gin.Context, id int32) {
	ctx := c.Request.Context()
	defer logger.TraceMethodAuto(ctx, id)()

	if err := s.deleteStockUC.Execute(ctx, id); err != nil {
		handleError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

// toListSakesInput パラメータをListSakesInputに変換する
func toListSakesInput(offset, limit, typeID, breweryID *int32) usecase.ListSakesInput {
	o := int32(0)
	if offset != nil {
		o = *offset
	}

	l := int32(20)
	if limit != nil {
		l = *limit
	}

	return usecase.ListSakesInput{
		KindID:    typeID,
		BreweryID: breweryID,
		Offset:    o,
		Limit:     l,
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
			Errors: &errors,
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
			Errors: &errors,
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
		Errors: &errors,
	})
}
