package presentation

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/sake-kasu/sake-hack-backend/api/generated"
	"github.com/sake-kasu/sake-hack-backend/internal/apperror"
	"github.com/sake-kasu/sake-hack-backend/internal/features/like/application/usecase"
	"github.com/sake-kasu/sake-hack-backend/internal/logger"
)

// LikeServerImpl いいね関連のServerInterface実装
type LikeServerImpl struct {
	toggleLikeUC usecase.ToggleLikeUsecaseInterface
	removeLikeUC usecase.RemoveLikeUsecaseInterface
}

// NewLikeServerImpl コンストラクタ
func NewLikeServerImpl(
	toggleLikeUC usecase.ToggleLikeUsecaseInterface,
	removeLikeUC usecase.RemoveLikeUsecaseInterface,
) *LikeServerImpl {
	return &LikeServerImpl{
		toggleLikeUC: toggleLikeUC,
		removeLikeUC: removeLikeUC,
	}
}

// CreateSakeLike いいねをつける
// (POST /sakes/{id}/likes)
func (s *LikeServerImpl) CreateSakeLike(c *gin.Context, id int32, params generated.CreateSakeLikeParams) {
	ctx := c.Request.Context()
	defer logger.TraceMethodAuto(ctx, map[string]interface{}{"sake_id": id})()

	if err := validateLikeToken(params.XLikeToken); err != nil {
		handleError(c, err)
		return
	}

	output, err := s.toggleLikeUC.Execute(ctx, usecase.ToggleLikeInput{
		SakeID: id,
		Token:  params.XLikeToken,
	})
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, generated.LikeResponse{
		LikeCount: output.LikeCount,
		IsLiked:   output.IsLiked,
	})
}

// DeleteSakeLike いいねを外す
// (DELETE /sakes/{id}/likes)
func (s *LikeServerImpl) DeleteSakeLike(c *gin.Context, id int32, params generated.DeleteSakeLikeParams) {
	ctx := c.Request.Context()
	defer logger.TraceMethodAuto(ctx, map[string]interface{}{"sake_id": id})()

	if err := validateLikeToken(params.XLikeToken); err != nil {
		handleError(c, err)
		return
	}

	output, err := s.removeLikeUC.Execute(ctx, usecase.RemoveLikeInput{
		SakeID: id,
		Token:  params.XLikeToken,
	})
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, generated.LikeResponse{
		LikeCount: output.LikeCount,
		IsLiked:   output.IsLiked,
	})
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
