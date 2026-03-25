package presentation

import (
	"github.com/go-playground/validator/v10"

	"github.com/sake-kasu/sake-hack-backend/internal/apperror"
)

var validate = validator.New()

// validateListParams はListSakes/ListStocksの共通パラメータをバリデーションする
func validateListParams(offset, limit int32) error {
	verr := apperror.NewValidationError("リクエストパラメータが不正です")

	if err := validate.Var(offset, "min=0"); err != nil {
		verr = verr.AddField("offset", "オフセットは0以上である必要があります")
	}

	if err := validate.Var(limit, "min=1,max=100"); err != nil {
		verr = verr.AddField("limit", "取得件数は1以上100以下である必要があります")
	}

	if verr.HasErrors() {
		return verr
	}
	return nil
}
