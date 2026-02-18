package presentation

import (
	"github.com/sake-kasu/sake-hack-backend/api/generated"
	"github.com/sake-kasu/sake-hack-backend/internal/features/sake/domain/entity"
	"github.com/sake-kasu/sake-hack-backend/internal/features/sake/application/usecase"
)

// toListSakesResponse ListSakesOutputをAPIレスポンスに変換
func toListSakesResponse(output *usecase.ListSakesOutput) generated.ListSakesResponse {
	sakes := make([]generated.Sake, 0, len(output.Sakes))
	for _, item := range output.Sakes {
		sakes = append(sakes, toSakeResponse(item))
	}

	return generated.ListSakesResponse{
		Data: &sakes,
		Meta: &generated.SakeListMeta{
			Total:  output.Pagination.Total,
			Offset: output.Pagination.Offset,
			Limit:  output.Pagination.Limit,
		},
	}
}

// toSakeResponse SakeListItemをAPI Sakeレスポンスに変換
func toSakeResponse(item entity.SakeListItem) generated.Sake {
	return generated.Sake{
		Id:           item.ID,
		Category:     generated.SakeCategory(item.Category),
		Name:         item.Name,
		ImagePreview: item.ImagePreview,
	}
}

// toSakeDetailResponse SakeDetailをAPI SakeDetailレスポンスに変換
func toSakeDetailResponse(detail *entity.SakeDetail) generated.SakeDetail {
	drinkStyles := make([]generated.DrinkStyle, 0, len(detail.DrinkStyles))
	for _, ds := range detail.DrinkStyles {
		drinkStyles = append(drinkStyles, generated.DrinkStyle{
			Id:          ds.ID,
			Name:        ds.Name,
			Description: ds.Description,
		})
	}

	return generated.SakeDetail{
		Id:       detail.ID,
		Category: generated.SakeCategory(detail.Category),
		Kind: generated.SakeKind{
			Id:   detail.Kind.ID,
			Name: detail.Kind.Name,
		},
		Brewery: generated.Brewery{
			Id:            detail.Brewery.ID,
			Name:          detail.Brewery.Name,
			OriginCountry: detail.Brewery.OriginCountry,
			OriginRegion:  detail.Brewery.OriginRegion,
			Latitude:      detail.Brewery.Latitude,
			Longitude:     detail.Brewery.Longitude,
		},
		Name: generated.SakeName{
			Name:     detail.Name.Name,
			Phonetic: detail.Name.Phonetic,
		},
		Abv:             detail.Abv,
		PurchaseVolume:  detail.PurchaseVolume,
		RemainingVolume: detail.RemainingVolume,
		Memo:            detail.Memo,
		DrinkStyles:     drinkStyles,
		Price:           detail.Price,
		ImageUrl:        detail.ImageUrl,
		CreatedAt:       detail.CreatedAt,
		UpdatedAt:       detail.UpdatedAt,
	}
}

// toCreateSakeResponse SakeListItemをCreateSakeResponseに変換
func toCreateSakeResponse(item *entity.SakeListItem) generated.CreateSakeResponse {
	sake := toSakeResponse(*item)
	return generated.CreateSakeResponse{
		Data: &sake,
	}
}

// toCreateStockInput CreateSakeRequestをCreateStockInputに変換
func toCreateStockInput(req generated.CreateSakeRequest) usecase.CreateStockInput {
	drinkStyles := make([]entity.DrinkStyle, 0, len(req.DrinkStyles))
	for _, ds := range req.DrinkStyles {
		drinkStyles = append(drinkStyles, entity.DrinkStyle{
			ID:          ds.Id,
			Name:        ds.Name,
			Description: ds.Description,
		})
	}

	return usecase.CreateStockInput{
		Category: entity.SakeCategory(req.Category),
		Kind: entity.SakeKind{
			ID:   req.Kind.Id,
			Name: req.Kind.Name,
		},
		Brewery: entity.Brewery{
			ID:            req.Brewery.Id,
			Name:          req.Brewery.Name,
			OriginCountry: req.Brewery.OriginCountry,
			OriginRegion:  req.Brewery.OriginRegion,
			Latitude:      req.Brewery.Latitude,
			Longitude:     req.Brewery.Longitude,
		},
		Name: entity.SakeName{
			Name:     req.Name.Name,
			Phonetic: req.Name.Phonetic,
		},
		Abv:             req.Abv,
		PurchaseVolume:  req.PurchaseVolume,
		RemainingVolume: req.RemainingVolume,
		Memo:            req.Memo,
		DrinkStyles:     drinkStyles,
		Price:           req.Price,
		ImageUrl:        req.ImageUrl,
	}
}

// toUpdateStockInput UpdateSakeRequestをUpdateStockInputに変換
func toUpdateStockInput(id int32, req generated.UpdateSakeRequest) usecase.UpdateStockInput {
	drinkStyles := make([]entity.DrinkStyle, 0, len(req.DrinkStyles))
	for _, ds := range req.DrinkStyles {
		drinkStyles = append(drinkStyles, entity.DrinkStyle{
			ID:          ds.Id,
			Name:        ds.Name,
			Description: ds.Description,
		})
	}

	return usecase.UpdateStockInput{
		ID:       id,
		Category: entity.SakeCategory(req.Category),
		Kind: entity.SakeKind{
			ID:   req.Kind.Id,
			Name: req.Kind.Name,
		},
		Brewery: entity.Brewery{
			ID:            req.Brewery.Id,
			Name:          req.Brewery.Name,
			OriginCountry: req.Brewery.OriginCountry,
			OriginRegion:  req.Brewery.OriginRegion,
			Latitude:      req.Brewery.Latitude,
			Longitude:     req.Brewery.Longitude,
		},
		Name: entity.SakeName{
			Name:     req.Name.Name,
			Phonetic: req.Name.Phonetic,
		},
		Abv:             req.Abv,
		PurchaseVolume:  req.PurchaseVolume,
		RemainingVolume: req.RemainingVolume,
		Memo:            req.Memo,
		DrinkStyles:     drinkStyles,
		Price:           req.Price,
		ImageUrl:        req.ImageUrl,
	}
}
