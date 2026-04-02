package presentation

import (
	"github.com/google/uuid"
	"github.com/sake-kasu/sake-hack-backend/api/generated"
	"github.com/sake-kasu/sake-hack-backend/internal/features/sake/application/usecase"
	"github.com/sake-kasu/sake-hack-backend/internal/features/sake/domain/entity"
	"github.com/sake-kasu/sake-hack-backend/internal/utils"
)

// toListSakesResponse ListSakesOutputをAPIレスポンスに変換
func toListSakesResponse(output *usecase.ListSakesOutput) generated.ListSakesResponse {
	sakes := make([]generated.Sake, 0, len(output.Sakes))
	for _, item := range output.Sakes {
		sakes = append(sakes, toSakeResponse(item))
	}

	return generated.ListSakesResponse{
		Data: sakes,
		Meta: generated.SakeListMeta{
			Total:  int(output.Pagination.Total),
			Offset: int(output.Pagination.Offset),
			Limit:  int(output.Pagination.Limit),
		},
	}
}

// toStockListResponse ListSakesOutputを在庫一覧レスポンスに変換
func toStockListResponse(output *usecase.ListSakesOutput) generated.ListStockResponse {
	stocks := make([]generated.Stock, 0, len(output.Sakes))
	for _, item := range output.Sakes {
		stocks = append(stocks, toStockResponse(item))
	}

	return generated.ListStockResponse{
		Data: stocks,
		Meta: generated.SakeListMeta{
			Total:  int(output.Pagination.Total),
			Offset: int(output.Pagination.Offset),
			Limit:  int(output.Pagination.Limit),
		},
	}
}

func toSakeResponse(item entity.SakeListItem) generated.Sake {
	return generated.Sake{
		SakeId:        generated.SakeID(item.ID),
		Category:      toGeneratedCategory(item.Category),
		Name:          generated.SakeName(item.Name),
		Phonetic:      nil,
		FirstImageKey: emptyToStringPtr(item.ImagePreview),
		LikeCount:     generated.LikeCount(0),
	}
}

func toStockResponse(item entity.SakeListItem) generated.Stock {
	return generated.Stock{
		SakeId:        generated.SakeID(item.ID),
		Category:      toGeneratedCategory(item.Category),
		Name:          generated.SakeName(item.Name),
		Phonetic:      nil,
		FirstImageKey: emptyToStringPtr(item.ImagePreview),
		LikeCount:     generated.LikeCount(0),
	}
}

// toSakeDetailResponse はドメインの詳細モデルを酒詳細レスポンスに変換する。
func toSakeDetailResponse(detail *entity.SakeDetail) generated.SakeDetail {
	tags := make([]generated.SakeTag, 0, len(detail.Tags))
	for _, tag := range detail.Tags {
		tags = append(tags, generated.SakeTag{
			TagId:   generated.SakeTagID(tag.ID),
			TagName: generated.SakeTagName(tag.Name),
		})
	}
	images := make([]generated.SakeImage, 0, len(detail.Images))
	for _, image := range detail.Images {
		images = append(images, generated.SakeImage{
			ImageId:   image.ID,
			ImageKey:  generated.SakeImageKey(image.ImageKey),
			SortOrder: generated.SakeSortOrder(image.SortOrder),
		})
	}
	return generated.SakeDetail{
		SakeId:            generated.SakeID(detail.ID),
		Name:              generated.SakeName(detail.Name.Name),
		Phonetic:          emptyToSakePhoneticPtr(detail.Name.Phonetic),
		Category:          toGeneratedCategory(detail.Category),
		AlcoholPercentage: float32ToAlcoholPercentagePtr(detail.Abv),
		Region:            stringToSakeRegionPtr(detail.Brewery.OriginRegion),
		Memo:              stringToSakeMemoPtr(detail.Memo),
		Tags:              tags,
		Images:            images,
		LikeCount:         generated.LikeCount(0),
		CreatedAt:         detail.CreatedAt,
		UpdatedAt:         detail.UpdatedAt,
	}
}

// toStockDetailResponse はドメインの詳細モデルを在庫詳細レスポンスに変換する。
func toStockDetailResponse(detail *entity.SakeDetail) generated.StockDetail {
	tags := make([]generated.SakeTag, 0, len(detail.Tags))
	for _, tag := range detail.Tags {
		tags = append(tags, generated.SakeTag{
			TagId:   generated.SakeTagID(tag.ID),
			TagName: generated.SakeTagName(tag.Name),
		})
	}
	images := make([]generated.SakeImage, 0, len(detail.Images))
	for _, image := range detail.Images {
		images = append(images, generated.SakeImage{
			ImageId:   image.ID,
			ImageKey:  generated.SakeImageKey(image.ImageKey),
			SortOrder: generated.SakeSortOrder(image.SortOrder),
		})
	}
	return generated.StockDetail{
		SakeId:            generated.SakeID(detail.ID),
		Name:              generated.SakeName(detail.Name.Name),
		Phonetic:          emptyToSakePhoneticPtr(detail.Name.Phonetic),
		Category:          toGeneratedCategory(detail.Category),
		AlcoholPercentage: float32ToAlcoholPercentagePtr(detail.Abv),
		VolumeMax:         int32ToVolumeMaxPtr(detail.PurchaseVolume),
		VolumeRemain:      int32ToVolumeRemainPtr(detail.RemainingVolume),
		Region:            stringToSakeRegionPtr(detail.Brewery.OriginRegion),
		Price:             int32ToPricePtr(detail.Price),
		Memo:              stringToSakeMemoPtr(detail.Memo),
		Tags:              tags,
		Images:            images,
		LikeCount:         generated.LikeCount(0),
		CreatedAt:         detail.CreatedAt,
		UpdatedAt:         detail.UpdatedAt,
	}
}

// toCreateStockInput は在庫登録リクエストをユースケース入力に変換する。
func toCreateStockInput(req generated.CreateStockRequest) usecase.CreateStockInput {
	return usecase.CreateStockInput{
		Category: entity.SakeCategory(req.Category),
		Kind:     entity.SakeKind{},
		Brewery: entity.Brewery{
			OriginRegion: sakeRegionToStringPtr(req.Region),
		},
		Name: entity.SakeName{
			Name:     string(req.Name),
			Phonetic: sakePhoneticToString(req.Phonetic),
		},
		Abv:             alcoholPercentageToFloat32Ptr(req.AlcoholPercentage),
		PurchaseVolume:  volumeMaxToInt32Ptr(req.VolumeMax),
		RemainingVolume: volumeRemainToInt32Ptr(req.VolumeRemain),
		Memo:            sakeMemoToStringPtr(req.Memo),
		DrinkStyles:     []entity.DrinkStyle{},
		TagNames:        toStringSlice(req.TagNames),
		Price:           priceToInt32Ptr(req.Price),
		ImageUrl:        firstImageKeyToStringPtr(req.ImageKeys),
		ImageKeys:       toImageKeyStrings(req.ImageKeys),
	}
}

// toUpdateStockInput は在庫更新リクエストをユースケース入力に変換する。
func toUpdateStockInput(id uuid.UUID, req generated.CreateStockRequest) usecase.UpdateStockInput {
	return usecase.UpdateStockInput{
		ID:       id,
		Category: entity.SakeCategory(req.Category),
		Kind:     entity.SakeKind{},
		Brewery: entity.Brewery{
			OriginRegion: sakeRegionToStringPtr(req.Region),
		},
		Name: entity.SakeName{
			Name:     string(req.Name),
			Phonetic: sakePhoneticToString(req.Phonetic),
		},
		Abv:             alcoholPercentageToFloat32Ptr(req.AlcoholPercentage),
		PurchaseVolume:  volumeMaxToInt32Ptr(req.VolumeMax),
		RemainingVolume: volumeRemainToInt32Ptr(req.VolumeRemain),
		Memo:            sakeMemoToStringPtr(req.Memo),
		DrinkStyles:     []entity.DrinkStyle{},
		TagNames:        toStringSlice(req.TagNames),
		Price:           priceToInt32Ptr(req.Price),
		ImageUrl:        firstImageKeyToStringPtr(req.ImageKeys),
		ImageKeys:       toImageKeyStrings(req.ImageKeys),
	}
}

func toGeneratedCategory(c entity.SakeCategory) generated.SakeCategory {
	switch c {
	case entity.SakeCategoryJapaneseSake:
		return generated.JAPANESESAKE
	case entity.SakeCategoryShochu:
		return generated.SHOCHU
	case entity.SakeCategoryAwamori:
		return generated.AWAMORI
	case entity.SakeCategoryBeer:
		return generated.BEER
	case entity.SakeCategoryWine:
		return generated.WINE
	default:
		return generated.OTHER
	}
}

func emptyToStringPtr(v string) *string {
	if v == "" {
		return nil
	}
	return &v
}

// emptyToSakePhoneticPtr は空文字を nil に正規化する。
func emptyToSakePhoneticPtr(v string) *generated.SakePhonetic {
	if v == "" {
		return nil
	}
	result := generated.SakePhonetic(v)
	return &result
}

// sakePhoneticToString は OpenAPI の nullable phonetic を文字列に変換する。
func sakePhoneticToString(v *generated.SakePhonetic) string {
	if v == nil {
		return ""
	}
	return string(*v)
}

// float32ToAlcoholPercentagePtr は *float32 を OpenAPI 型に変換する。
func float32ToAlcoholPercentagePtr(v *float32) *generated.SakeAlcoholPercentage {
	if v == nil {
		return nil
	}
	result := generated.SakeAlcoholPercentage(*v)
	return &result
}

// stringToSakeRegionPtr は *string を OpenAPI 型に変換する。
func stringToSakeRegionPtr(v *string) *generated.SakeRegion {
	if v == nil {
		return nil
	}
	result := generated.SakeRegion(*v)
	return &result
}

// stringToSakeMemoPtr は *string を OpenAPI 型に変換する。
func stringToSakeMemoPtr(v *string) *generated.SakeMemo {
	if v == nil {
		return nil
	}
	result := generated.SakeMemo(*v)
	return &result
}

// sakeMemoToStringPtr は OpenAPI の nullable memo を *string に変換する。
func sakeMemoToStringPtr(v *generated.SakeMemo) *string {
	if v == nil {
		return nil
	}
	result := string(*v)
	return &result
}

// int32ToVolumeMaxPtr は *int32 を OpenAPI 型に変換する。
func int32ToVolumeMaxPtr(v *int32) *generated.SakeVolumeMax {
	if v == nil {
		return nil
	}
	result := generated.SakeVolumeMax(*v)
	return &result
}

// int32ToVolumeRemainPtr は *int32 を OpenAPI 型に変換する。
func int32ToVolumeRemainPtr(v *int32) *generated.SakeVolumeRemain {
	if v == nil {
		return nil
	}
	result := generated.SakeVolumeRemain(*v)
	return &result
}

// int32ToPricePtr は *int32 を OpenAPI 型に変換する。
func int32ToPricePtr(v *int32) *generated.SakePrice {
	if v == nil {
		return nil
	}
	result := generated.SakePrice(*v)
	return &result
}

// alcoholPercentageToFloat32Ptr は OpenAPI 型を *float32 に変換する。
func alcoholPercentageToFloat32Ptr(v *generated.SakeAlcoholPercentage) *float32 {
	if v == nil {
		return nil
	}
	result := float32(*v)
	return &result
}

// volumeMaxToInt32Ptr は OpenAPI 型を *int32 に変換する。
func volumeMaxToInt32Ptr(v *generated.SakeVolumeMax) *int32 {
	if v == nil {
		return nil
	}
	result, err := utils.IntToInt32(int(*v))
	if err != nil {
		return nil
	}
	return &result
}

// volumeRemainToInt32Ptr は OpenAPI 型を *int32 に変換する。
func volumeRemainToInt32Ptr(v *generated.SakeVolumeRemain) *int32 {
	if v == nil {
		return nil
	}
	result, err := utils.IntToInt32(int(*v))
	if err != nil {
		return nil
	}
	return &result
}

// priceToInt32Ptr は OpenAPI 型を *int32 に変換する。
func priceToInt32Ptr(v *generated.SakePrice) *int32 {
	if v == nil {
		return nil
	}
	result, err := utils.IntToInt32(int(*v))
	if err != nil {
		return nil
	}
	return &result
}

// sakeRegionToStringPtr は OpenAPI の nullable region を *string に変換する。
func sakeRegionToStringPtr(v *generated.SakeRegion) *string {
	if v == nil {
		return nil
	}
	result := string(*v)
	return &result
}

// firstImageKeyToStringPtr は先頭画像キーのみを保存用に取り出す。
func firstImageKeyToStringPtr(v []generated.SakeImageKey) *string {
	if len(v) == 0 {
		return nil
	}
	result := string(v[0])
	return &result
}

// toStringSlice は generated の文字列配列を通常の文字列配列に変換する。
func toStringSlice[T ~string](values []T) []string {
	result := make([]string, 0, len(values))
	for _, value := range values {
		result = append(result, string(value))
	}
	return result
}

// toImageKeyStrings は image key 配列を保存用の文字列配列に変換する。
func toImageKeyStrings(values []generated.SakeImageKey) []string {
	return toStringSlice(values)
}
