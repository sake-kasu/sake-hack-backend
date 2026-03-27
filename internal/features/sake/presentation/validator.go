package presentation

import (
	"math"

	"github.com/go-playground/validator/v10"

	"github.com/sake-kasu/sake-hack-backend/api/generated"
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

// validateCreateStockRequest は在庫登録リクエストを検証する。
func validateCreateStockRequest(req generated.CreateStockRequest) error {
	verr := apperror.NewValidationError("在庫登録リクエストが不正です")

	if req.Name == "" {
		verr = verr.AddField("name", "酒名は必須です")
	}
	if len(req.TagNames) == 0 {
		verr = verr.AddField("tagNames", "タグは1件以上指定してください")
	}
	if len(req.ImageKeys) == 0 {
		verr = verr.AddField("imageKeys", "画像キーは1件以上指定してください")
	}
	if req.AlcoholPercentage != nil && (*req.AlcoholPercentage < 0 || *req.AlcoholPercentage > 100) {
		verr = verr.AddField("alcoholPercentage", "アルコール度数は0以上100以下である必要があります")
	}
	if req.VolumeMax != nil && *req.VolumeMax <= 0 {
		verr = verr.AddField("volumeMax", "最大容量は1以上である必要があります")
	}
	if req.VolumeRemain != nil && (*req.VolumeRemain < 0 || *req.VolumeRemain > 100) {
		verr = verr.AddField("volumeRemain", "残量は0以上100以下である必要があります")
	}
	if req.Price != nil && *req.Price < 0 {
		verr = verr.AddField("price", "価格は0以上である必要があります")
	}
	if req.Price != nil && *req.Price > generated.SakePrice(math.MaxInt32) {
		verr = verr.AddField("price", "価格はint32の範囲内で指定してください")
	}

	if verr.HasErrors() {
		return verr
	}
	return nil
}
