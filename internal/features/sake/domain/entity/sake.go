package entity

import (
	"time"

	"github.com/google/uuid"
)

// SakeCategory 酒のカテゴリ(Value Object)
type SakeCategory string

const (
	SakeCategoryJapaneseSake SakeCategory = "JAPANESE_SAKE"
	SakeCategoryWhisky       SakeCategory = "WHISKY"
	SakeCategoryWine         SakeCategory = "WINE"
	SakeCategoryBeer         SakeCategory = "BEER"
	SakeCategoryShochu       SakeCategory = "SHOCHU"
	SakeCategoryAwamori      SakeCategory = "AWAMORI"
	SakeCategoryRiqueur      SakeCategory = "RIQUEUR"
	SakeCategorySpirits      SakeCategory = "SPIRITS"
	SakeCategoryOther        SakeCategory = "OTHER"
)

// IsValid はカテゴリが有効か検証する
func (c SakeCategory) IsValid() bool {
	switch c {
	case SakeCategoryJapaneseSake, SakeCategoryWhisky, SakeCategoryWine,
		SakeCategoryBeer, SakeCategoryShochu, SakeCategoryAwamori,
		SakeCategoryRiqueur, SakeCategorySpirits, SakeCategoryOther:
		return true
	}
	return false
}

// SakeKind 酒の種類(小分類)
type SakeKind struct {
	ID   int32
	Name string
}

// Brewery 酒造
type Brewery struct {
	ID            int32
	Name          string
	OriginCountry string
	OriginRegion  *string
	Latitude      *float64
	Longitude     *float64
}

// DrinkStyle 飲み方
type DrinkStyle struct {
	ID          int32
	Name        string
	Description *string
}

// SakeName 酒名(名前 + 読み)
type SakeName struct {
	Name     string
	Phonetic string
}

// SakeListItem リスト表示用の酒情報(軽量)
type SakeListItem struct {
	ID           uuid.UUID
	Category     SakeCategory
	Name         string
	ImagePreview string
}

// SakeDetail 酒の詳細情報
type SakeDetail struct {
	ID              int32
	Category        SakeCategory
	Kind            SakeKind
	Brewery         Brewery
	Name            SakeName
	Abv             float32
	PurchaseVolume  float32
	RemainingVolume float32
	Memo            *string
	DrinkStyles     []DrinkStyle
	Price           int32
	ImageUrl        *string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// Pagination ページネーション情報
type Pagination struct {
	Total  int64
	Offset int32
	Limit  int32
}
