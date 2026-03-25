package presentation

import (
	"encoding/binary"

	"github.com/google/uuid"
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
		SakeId:        int32ToSakeID(item.ID),
		Category:      toGeneratedCategory(item.Category),
		Name:          generated.SakeName(item.Name),
		Phonetic:      nil,
		FirstImageKey: emptyToStringPtr(item.ImagePreview),
		LikeCount:     generated.LikeCount(0),
	}
}

func toStockResponse(item entity.SakeListItem) generated.Stock {
	return generated.Stock{
		SakeId:        int32ToSakeID(item.ID),
		Category:      toGeneratedCategory(item.Category),
		Name:          generated.SakeName(item.Name),
		Phonetic:      nil,
		FirstImageKey: emptyToStringPtr(item.ImagePreview),
		LikeCount:     generated.LikeCount(0),
	}
}

func int32ToSakeID(id int32) generated.SakeID {
	var raw uuid.UUID
	binary.BigEndian.PutUint32(raw[12:], uint32(id))
	return generated.SakeID(raw)
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
