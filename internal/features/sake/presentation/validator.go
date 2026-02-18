package presentation

import (
	"fmt"

	"github.com/sake-kasu/sake-hack-backend/api/generated"
	"github.com/sake-kasu/sake-hack-backend/internal/apperror"
)

// ValidateGetSakesParams GET /sakes のパラメータをバリデーション
func ValidateGetSakesParams(params generated.GetSakesParams) error {
	validationErr := apperror.NewValidationError("パラメータのバリデーションに失敗しました")

	// Offset のバリデーション
	if params.Offset != nil && *params.Offset < 0 {
		validationErr.AddField("offset", fmt.Sprintf("offset must be greater than or equal to 0 (got: %d)", *params.Offset))
	}

	// Limit のバリデーション
	if params.Limit != nil {
		if *params.Limit < 1 {
			validationErr.AddField("limit", fmt.Sprintf("limit must be greater than or equal to 1 (got: %d)", *params.Limit))
		} else if *params.Limit > 100 {
			validationErr.AddField("limit", fmt.Sprintf("limit must be less than or equal to 100 (got: %d)", *params.Limit))
		}
	}

	// Category のバリデーション（enumは自動生成された型で保証されている）
	// Search のバリデーション（必要に応じて追加）

	// エラーがある場合のみ返す
	if validationErr.HasErrors() {
		return validationErr
	}

	return nil
}
