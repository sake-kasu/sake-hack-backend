package presentation

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"github.com/sake-kasu/sake-hack-backend/api/generated"
	"github.com/sake-kasu/sake-hack-backend/internal/apperror"
	"github.com/sake-kasu/sake-hack-backend/internal/features/sake/application/usecase"
	"github.com/sake-kasu/sake-hack-backend/internal/features/sake/domain/entity"
	"github.com/sake-kasu/sake-hack-backend/internal/logger"
)

// validate はgo-playground/validatorのインスタンス
var validate = validator.New()

// SakeServerImpl 酒関連のServerInterface実装
type SakeServerImpl struct {
	listSakesUsecase usecase.ListSakesUsecaseInterface
}

// NewSakeServerImpl コンストラクタ
func NewSakeServerImpl(listSakesUsecase usecase.ListSakesUsecaseInterface) *SakeServerImpl {
	return &SakeServerImpl{
		listSakesUsecase: listSakesUsecase,
	}
}

// ListSakes 酒一覧取得
// (GET /sakes)
func (s *SakeServerImpl) ListSakes(c *gin.Context, params generated.ListSakesParams) {
	ctx := c.Request.Context()
	defer logger.TraceMethodAuto(ctx, params)()

	// バリデーション
	if err := validateListSakesParams(params); err != nil {
		handleError(c, err)
		return
	}

	// デフォルト値設定
	offset := int32(0)
	if params.Offset != nil {
		offset = *params.Offset
	}

	limit := int32(20)
	if params.Limit != nil {
		limit = *params.Limit
	}

	// Usecase実行
	output, err := s.listSakesUsecase.Execute(ctx, usecase.ListSakesInput{
		TypeID:    params.TypeId,
		BreweryID: params.BreweryId,
		Offset:    offset,
		Limit:     limit,
	})
	if err != nil {
		handleError(c, err)
		return
	}

	// レスポンス変換
	response := toListSakesResponse(output)
	c.JSON(http.StatusOK, response)
}

// validateListSakesParams はListSakesのパラメータをバリデーションする
func validateListSakesParams(params generated.ListSakesParams) error {
	verr := apperror.NewValidationError("リクエストパラメータが不正です")

	// go-playground/validator を使ったバリデーション
	// ポインタ型のため、nil チェック後に値をバリデーション
	if params.Offset != nil {
		if err := validate.Var(*params.Offset, "min=0"); err != nil {
			verr = verr.AddField("offset", "オフセットは0以上である必要があります")
		}
	}

	if params.Limit != nil {
		if err := validate.Var(*params.Limit, "min=1,max=100"); err != nil {
			verr = verr.AddField("limit", "取得件数は1以上100以下である必要があります")
		}
	}

	if params.TypeId != nil {
		if err := validate.Var(*params.TypeId, "min=1"); err != nil {
			verr = verr.AddField("type_id", "酒の種類IDは1以上である必要があります")
		}
	}

	if params.BreweryId != nil {
		if err := validate.Var(*params.BreweryId, "min=1"); err != nil {
			verr = verr.AddField("brewery_id", "酒造IDは1以上である必要があります")
		}
	}

	if verr.HasErrors() {
		return verr
	}
	return nil
}

// 【参考】手動バリデーションの例（サンプルとして残す）
//
// 手動バリデーションの実装例:
/*
func validateListSakesParamsManual(params generated.ListSakesParams) error {
	verr := apperror.NewValidationError("リクエストパラメータが不正です")

	// offset: minimum 0
	if params.Offset != nil && *params.Offset < 0 {
		verr = verr.AddField("offset", "オフセットは0以上である必要があります")
	}

	// limit: minimum 1, maximum 100
	if params.Limit != nil {
		if *params.Limit < 1 {
			verr = verr.AddField("limit", "取得件数は1以上である必要があります")
		}
		if *params.Limit > 100 {
			verr = verr.AddField("limit", "取得件数は100以下である必要があります")
		}
	}

	// type_id: minimum 1 (optional)
	if params.TypeId != nil && *params.TypeId < 1 {
		verr = verr.AddField("type_id", "酒の種類IDは1以上である必要があります")
	}

	// brewery_id: minimum 1 (optional)
	if params.BreweryId != nil && *params.BreweryId < 1 {
		verr = verr.AddField("brewery_id", "酒造IDは1以上である必要があります")
	}

	if verr.HasErrors() {
		return verr
	}
	return nil
}
*/

// toListSakesResponse DomainエンティティをAPIレスポンスに変換
func toListSakesResponse(output *usecase.ListSakesOutput) generated.ListSakesResponse {
	sakes := make([]generated.Sake, 0, len(output.Sakes))
	for _, sake := range output.Sakes {
		sakes = append(sakes, toSakeResponse(sake))
	}

	meta := generated.SakeListMeta{
		Total:  output.Pagination.Total,
		Offset: output.Pagination.Offset,
		Limit:  output.Pagination.Limit,
	}

	return generated.ListSakesResponse{
		Data:   &sakes,
		Meta:   &meta,
		Errors: nil,
	}
}

// toSakeResponse DomainエンティティをAPIレスポンスに変換
// NOTE: entity.Sake にはまだ Category/ImagePreview フィールドがないため暫定マッピング
func toSakeResponse(sake entity.Sake) generated.Sake {
	return generated.Sake{
		Id:           sake.ID,
		Category:     generated.SakeCategory(sake.Type.Name),
		Name:         sake.Name,
		ImagePreview: "",
	}
}

// GetSakeDetail 酒詳細取得
// (GET /sakes/{id})
func (s *SakeServerImpl) GetSakeDetail(c *gin.Context, id int32) {
	c.JSON(http.StatusNotImplemented, generated.ErrorResponse{
		Errors: &[]generated.APIError{{Code: "NOT_IMPLEMENTED", Message: "未実装です"}},
	})
}

// ListStocks 在庫一覧取得
// (GET /stocks)
func (s *SakeServerImpl) ListStocks(c *gin.Context, params generated.ListStocksParams) {
	c.JSON(http.StatusNotImplemented, generated.ErrorResponse{
		Errors: &[]generated.APIError{{Code: "NOT_IMPLEMENTED", Message: "未実装です"}},
	})
}

// CreateStock 在庫登録
// (POST /stocks)
func (s *SakeServerImpl) CreateStock(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, generated.ErrorResponse{
		Errors: &[]generated.APIError{{Code: "NOT_IMPLEMENTED", Message: "未実装です"}},
	})
}

// GetStockDetail 在庫詳細取得
// (GET /stocks/{id})
func (s *SakeServerImpl) GetStockDetail(c *gin.Context, id int32) {
	c.JSON(http.StatusNotImplemented, generated.ErrorResponse{
		Errors: &[]generated.APIError{{Code: "NOT_IMPLEMENTED", Message: "未実装です"}},
	})
}

// UpdateStock 在庫更新
// (PUT /stocks/{id})
func (s *SakeServerImpl) UpdateStock(c *gin.Context, id int32) {
	c.JSON(http.StatusNotImplemented, generated.ErrorResponse{
		Errors: &[]generated.APIError{{Code: "NOT_IMPLEMENTED", Message: "未実装です"}},
	})
}

// DeleteStock 在庫削除
// (DELETE /stocks/{id})
func (s *SakeServerImpl) DeleteStock(c *gin.Context, id int32) {
	c.JSON(http.StatusNotImplemented, generated.ErrorResponse{
		Errors: &[]generated.APIError{{Code: "NOT_IMPLEMENTED", Message: "未実装です"}},
	})
}

// handleError エラーをHTTPレスポンスに変換
func handleError(c *gin.Context, err error) {
	// ValidationErrorのチェック
	if valErr, ok := err.(*apperror.ValidationError); ok {
		// ValidationErrorの場合、Fieldsも含めてレスポンス
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

	// 通常のAppErrorのチェック
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

	// 予期しないエラー
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
