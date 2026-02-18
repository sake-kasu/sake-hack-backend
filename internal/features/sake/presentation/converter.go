package presentation

import (
	"github.com/sake-kasu/sake-hack-backend/api/generated"
	"github.com/sake-kasu/sake-hack-backend/internal/features/sake/application/usecase"
	"github.com/sake-kasu/sake-hack-backend/internal/features/sake/domain/entity"
	"github.com/sake-kasu/sake-hack-backend/internal/utils"
)

// toListSakesResponse ListSakesOutputをAPIレスポンスに変換
func toListSakesResponse(output *usecase.ListSakesOutput) generated.ListSakesResponse {
	sakes := make([]generated.SakeDetail, 0, len(output.Sakes))
	for _, item := range output.Sakes {
		sakes = append(sakes, toSakeDetail(item))
	}

	return generated.ListSakesResponse{
		Total: int(output.Total),
		Sakes: sakes,
	}
}

// toSakeDetail entity.SakeDetailをgenerated.SakeDetailに変換（公開用）
func toSakeDetail(detail entity.SakeDetail) generated.SakeDetail {
	// int32 IDを決定的なUUIDに変換
	id := utils.Int32ToUUID(detail.ID)

	return generated.SakeDetail{
		Id:                id,
		Name:              generated.SakeName(detail.Name.Name),
		ImageUrl:          detail.ImageUrl,
		Category:          generated.SakeCategory(detail.Category),
		Description:       &detail.Kind.Name,
		AlcoholPercentage: &detail.Abv,
		Region:            detail.Brewery.OriginRegion,
		Memo:              detail.Memo,
		CreatedAt:         detail.CreatedAt,
		UpdatedAt:         detail.UpdatedAt,
	}
}

