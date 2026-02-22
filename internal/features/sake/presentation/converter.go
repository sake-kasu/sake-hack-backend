package presentation

import (
	"context"

	"go.uber.org/zap"

	"github.com/sake-kasu/sake-hack-backend/api/generated"
	"github.com/sake-kasu/sake-hack-backend/internal/features/sake/application/port"
	"github.com/sake-kasu/sake-hack-backend/internal/features/sake/application/usecase"
	"github.com/sake-kasu/sake-hack-backend/internal/features/sake/domain/entity"
	"github.com/sake-kasu/sake-hack-backend/internal/logger"
)

// sakeConverter はレスポンス変換を行う構造体
type sakeConverter struct {
	urlResolver port.URLResolver
}

// newSakeConverter はコンストラクタ
func newSakeConverter(urlResolver port.URLResolver) *sakeConverter {
	return &sakeConverter{urlResolver: urlResolver}
}

// resolveImageURL はオブジェクトキーから署名付きURLを生成する
func (conv *sakeConverter) resolveImageURL(ctx context.Context, objectKey *string) string {
	if objectKey == nil || *objectKey == "" {
		return ""
	}
	url, err := conv.urlResolver.ResolveURL(ctx, *objectKey)
	if err != nil {
		logger.Warn(ctx, "署名付きURLの生成に失敗しました", zap.String("object_key", *objectKey), zap.Error(err))
		return ""
	}
	return url
}

// toListSakesResponse ListSakesOutputをAPIレスポンスに変換
func (conv *sakeConverter) toListSakesResponse(ctx context.Context, output *usecase.ListSakesOutput) generated.ListSakesResponse {
	sakes := make([]generated.Sake, 0, len(output.Sakes))
	for _, item := range output.Sakes {
		sakes = append(sakes, conv.toSakeResponse(ctx, item))
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
func (conv *sakeConverter) toSakeResponse(ctx context.Context, item entity.SakeListItem) generated.Sake {
	return generated.Sake{
		Id:           item.ID,
		Category:     generated.SakeCategory(item.Category),
		Name:         item.Name,
		ImagePreview: conv.resolveImageURL(ctx, item.ObjectKey),
		LikeCount:    item.LikeCount,
		IsLiked:      item.IsLiked,
	}
}

// toSakeDetailResponse SakeDetailをAPI SakeDetailレスポンスに変換
func (conv *sakeConverter) toSakeDetailResponse(ctx context.Context, detail *entity.SakeDetail) generated.SakeDetail {
	drinkStyles := make([]generated.DrinkStyle, 0, len(detail.DrinkStyles))
	for _, ds := range detail.DrinkStyles {
		drinkStyles = append(drinkStyles, generated.DrinkStyle{
			Id:          ds.ID,
			Name:        ds.Name,
			Description: ds.Description,
		})
	}

	imageURL := conv.resolveImageURL(ctx, detail.ObjectKey)

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
		LikeCount:       detail.LikeCount,
		IsLiked:         detail.IsLiked,
		ObjectKey:       detail.ObjectKey,
		ImageUrl:        &imageURL,
		CreatedAt:       detail.CreatedAt,
		UpdatedAt:       detail.UpdatedAt,
	}
}

// toListKindsResponse ListKindsOutputをAPIレスポンスに変換
func (conv *sakeConverter) toListKindsResponse(output *usecase.ListKindsOutput) generated.ListKindsResponse {
	kinds := make([]generated.SakeKind, 0, len(output.Kinds))
	for _, k := range output.Kinds {
		kinds = append(kinds, generated.SakeKind{
			Id:   k.ID,
			Name: k.Name,
		})
	}

	return generated.ListKindsResponse{
		Data: &kinds,
	}
}

// toListBreweriesResponse ListBreweriesOutputをAPIレスポンスに変換
func (conv *sakeConverter) toListBreweriesResponse(output *usecase.ListBreweriesOutput) generated.ListBreweriesResponse {
	breweries := make([]generated.Brewery, 0, len(output.Breweries))
	for _, b := range output.Breweries {
		breweries = append(breweries, generated.Brewery{
			Id:            b.ID,
			Name:          b.Name,
			OriginCountry: b.OriginCountry,
			OriginRegion:  b.OriginRegion,
			Latitude:      b.Latitude,
			Longitude:     b.Longitude,
		})
	}

	return generated.ListBreweriesResponse{
		Data: &breweries,
	}
}

// toListDrinkStylesResponse ListDrinkStylesOutputをAPIレスポンスに変換
func (conv *sakeConverter) toListDrinkStylesResponse(output *usecase.ListDrinkStylesOutput) generated.ListDrinkStylesResponse {
	drinkStyles := make([]generated.DrinkStyle, 0, len(output.DrinkStyles))
	for _, ds := range output.DrinkStyles {
		drinkStyles = append(drinkStyles, generated.DrinkStyle{
			Id:          ds.ID,
			Name:        ds.Name,
			Description: ds.Description,
		})
	}

	return generated.ListDrinkStylesResponse{
		Data: &drinkStyles,
	}
}

// toCreateSakeResponse SakeListItemをCreateSakeResponseに変換
func (conv *sakeConverter) toCreateSakeResponse(ctx context.Context, item *entity.SakeListItem) generated.CreateSakeResponse {
	sake := conv.toSakeResponse(ctx, *item)
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
		ObjectKey:       req.ObjectKey,
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
		ObjectKey:       req.ObjectKey,
	}
}
