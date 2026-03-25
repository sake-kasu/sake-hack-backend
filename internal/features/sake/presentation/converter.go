package presentation

import (
	"github.com/sake-kasu/sake-hack-backend/api/generated"
	"github.com/sake-kasu/sake-hack-backend/internal/features/sake/application/usecase"
	"github.com/sake-kasu/sake-hack-backend/internal/features/sake/domain/entity"
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
	return generated.SakeDetail{
		SakeId:            generated.SakeID(detail.ID),
		Name:              generated.SakeName(detail.Name.Name),
		Phonetic:          emptyToSakePhoneticPtr(detail.Name.Phonetic),
		Category:          toGeneratedCategory(detail.Category),
		AlcoholPercentage: float32ToAlcoholPercentagePtr(detail.Abv),
		Region:            stringToSakeRegionPtr(detail.Brewery.OriginRegion),
		Memo:              stringToSakeMemoPtr(detail.Memo),
		Tags:              []generated.SakeTag{},
		Images:            []generated.SakeImage{},
		LikeCount:         generated.LikeCount(0),
		CreatedAt:         detail.CreatedAt,
		UpdatedAt:         detail.UpdatedAt,
	}
}

// toStockDetailResponse はドメインの詳細モデルを在庫詳細レスポンスに変換する。
func toStockDetailResponse(detail *entity.SakeDetail) generated.StockDetail {
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
		Tags:              []generated.SakeTag{},
		Images:            []generated.SakeImage{},
		LikeCount:         generated.LikeCount(0),
		CreatedAt:         detail.CreatedAt,
		UpdatedAt:         detail.UpdatedAt,
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
