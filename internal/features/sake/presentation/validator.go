package presentation

import (
	"strings"

	"github.com/go-playground/validator/v10"

	"github.com/sake-kasu/sake-hack-backend/api/generated"
	"github.com/sake-kasu/sake-hack-backend/internal/apperror"
)

// validate はgo-playground/validatorのインスタンス
var validate = validator.New()

// validateListParams はListSakes/ListStocksの共通パラメータをバリデーションする
func validateListParams(offset, limit, typeID, breweryID *int32) error {
	verr := apperror.NewValidationError("リクエストパラメータが不正です")

	if offset != nil {
		if err := validate.Var(*offset, "min=0"); err != nil {
			verr = verr.AddField("offset", "オフセットは0以上である必要があります")
		}
	}

	if limit != nil {
		if err := validate.Var(*limit, "min=1,max=100"); err != nil {
			verr = verr.AddField("limit", "取得件数は1以上100以下である必要があります")
		}
	}

	if typeID != nil {
		if err := validate.Var(*typeID, "min=1"); err != nil {
			verr = verr.AddField("type_id", "酒の種類IDは1以上である必要があります")
		}
	}

	if breweryID != nil {
		if err := validate.Var(*breweryID, "min=1"); err != nil {
			verr = verr.AddField("brewery_id", "酒造IDは1以上である必要があります")
		}
	}

	if verr.HasErrors() {
		return verr
	}
	return nil
}

// validateCreateSakeRequest はCreateSakeRequestをバリデーションする
func validateCreateSakeRequest(req generated.CreateSakeRequest) error {
	verr := apperror.NewValidationError("在庫登録リクエストが不正です")

	if req.Name.Name == "" {
		verr = verr.AddField("name.name", "酒名は必須です")
	}

	if req.Brewery.Name == "" {
		verr = verr.AddField("brewery.name", "酒造名は必須です")
	}

	if req.Brewery.OriginCountry == "" {
		verr = verr.AddField("brewery.originCountry", "所在国は必須です")
	}

	if req.Kind.Name == "" {
		verr = verr.AddField("kind.name", "酒の種類名は必須です")
	}

	if req.Abv < 0 || req.Abv > 100 {
		verr = verr.AddField("abv", "アルコール度数は0以上100以下である必要があります")
	}

	if req.PurchaseVolume < 0 {
		verr = verr.AddField("purchaseVolume", "購入時容量は0以上である必要があります")
	}

	if req.RemainingVolume < 0 {
		verr = verr.AddField("remainingVolume", "残容量は0以上である必要があります")
	}

	if req.Price < 0 {
		verr = verr.AddField("price", "価格は0以上である必要があります")
	}

	if verr.HasErrors() {
		return verr
	}
	return nil
}

// allowedImageContentTypes は画像アップロードで許可されるContent-Type
var allowedImageContentTypes = map[string]bool{
	"image/jpeg": true,
	"image/png":  true,
	"image/gif":  true,
	"image/webp": true,
}

// maxFilenameLength はファイル名の最大長
const maxFilenameLength = 255

// validateUploadUrlRequest は画像アップロードURL発行リクエストをバリデーションする
func validateUploadUrlRequest(contentType, filename string) error {
	verr := apperror.NewValidationError("アップロードURLリクエストが不正です")

	if contentType == "" {
		verr = verr.AddField("contentType", "Content-Typeは必須です")
	} else if !allowedImageContentTypes[contentType] {
		verr = verr.AddField("contentType", "許可されていないContent-Typeです(image/jpeg, image/png, image/gif, image/webpのみ)")
	}

	if filename == "" {
		verr = verr.AddField("filename", "ファイル名は必須です")
	} else if len(filename) > maxFilenameLength {
		verr = verr.AddField("filename", "ファイル名は255文字以内である必要があります")
	} else if strings.Contains(filename, "..") || strings.ContainsAny(filename, "/\\") {
		verr = verr.AddField("filename", "ファイル名に不正な文字が含まれています")
	}

	if verr.HasErrors() {
		return verr
	}
	return nil
}

// validatePatchStockRequest はPatchStockRequestをバリデーションする
func validatePatchStockRequest(req generated.PatchStockRequest) error {
	verr := apperror.NewValidationError("在庫部分更新リクエストが不正です")

	if req.ObjectKey == "" {
		verr = verr.AddField("objectKey", "オブジェクトキーは必須です")
	} else if err := validateObjectKey(req.ObjectKey); err != nil {
		verr = verr.AddField("objectKey", err.Error())
	}

	if verr.HasErrors() {
		return verr
	}
	return nil
}

// validateObjectKey はオブジェクトキーの形式を検証する
func validateObjectKey(key string) error {
	if strings.Contains(key, "..") {
		return apperror.BadRequestError("オブジェクトキーに不正な文字列が含まれています")
	}
	if !strings.HasPrefix(key, "sakes/") {
		return apperror.BadRequestError("オブジェクトキーの形式が不正です")
	}
	return nil
}

// validateUpdateSakeRequest はUpdateSakeRequestをバリデーションする
func validateUpdateSakeRequest(req generated.UpdateSakeRequest) error {
	verr := apperror.NewValidationError("在庫更新リクエストが不正です")

	if req.Name.Name == "" {
		verr = verr.AddField("name.name", "酒名は必須です")
	}

	if req.Brewery.Name == "" {
		verr = verr.AddField("brewery.name", "酒造名は必須です")
	}

	if req.Brewery.OriginCountry == "" {
		verr = verr.AddField("brewery.originCountry", "所在国は必須です")
	}

	if req.Kind.Name == "" {
		verr = verr.AddField("kind.name", "酒の種類名は必須です")
	}

	if req.Abv < 0 || req.Abv > 100 {
		verr = verr.AddField("abv", "アルコール度数は0以上100以下である必要があります")
	}

	if req.PurchaseVolume < 0 {
		verr = verr.AddField("purchaseVolume", "購入時容量は0以上である必要があります")
	}

	if req.RemainingVolume < 0 {
		verr = verr.AddField("remainingVolume", "残容量は0以上である必要があります")
	}

	if req.Price < 0 {
		verr = verr.AddField("price", "価格は0以上である必要があります")
	}

	if verr.HasErrors() {
		return verr
	}
	return nil
}
