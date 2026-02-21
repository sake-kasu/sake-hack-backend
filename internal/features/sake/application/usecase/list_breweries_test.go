package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/sake-kasu/sake-hack-backend/internal/features/sake/domain/entity"
)

// テスト: 酒造一覧の正常取得(デフォルトパラメータ)
func TestListBreweriesUsecase_Execute_Success(t *testing.T) {
	mockQuery := new(MockSakeQuery)
	uc := NewListBreweriesUsecase(mockQuery)

	region := "山口県"
	lat := 34.1234
	lng := 131.5678
	expectedBreweries := []entity.Brewery{
		{ID: 1, Name: "旭酒造", OriginCountry: "日本", OriginRegion: &region, Latitude: &lat, Longitude: &lng},
	}

	mockQuery.On("ListBreweries", mock.Anything, (*string)(nil), int32(50)).Return(expectedBreweries, nil)

	input := ListBreweriesInput{
		Keyword: nil,
		Limit:   50,
	}

	output, err := uc.Execute(context.Background(), input)

	assert.NoError(t, err)
	assert.NotNil(t, output)
	assert.Len(t, output.Breweries, 1)
	assert.Equal(t, "旭酒造", output.Breweries[0].Name)
	assert.Equal(t, "日本", output.Breweries[0].OriginCountry)
	assert.Equal(t, &region, output.Breweries[0].OriginRegion)
	mockQuery.AssertExpectations(t)
}

// テスト: キーワード付き検索
func TestListBreweriesUsecase_Execute_Success_WithKeyword(t *testing.T) {
	mockQuery := new(MockSakeQuery)
	uc := NewListBreweriesUsecase(mockQuery)

	keyword := "旭"
	expectedBreweries := []entity.Brewery{
		{ID: 1, Name: "旭酒造", OriginCountry: "日本"},
		{ID: 2, Name: "旭日酒造", OriginCountry: "日本"},
	}

	mockQuery.On("ListBreweries", mock.Anything, &keyword, int32(50)).Return(expectedBreweries, nil)

	input := ListBreweriesInput{
		Keyword: &keyword,
		Limit:   50,
	}

	output, err := uc.Execute(context.Background(), input)

	assert.NoError(t, err)
	assert.NotNil(t, output)
	assert.Len(t, output.Breweries, 2)
	mockQuery.AssertExpectations(t)
}

// テスト: limitが範囲外の場合はデフォルト値(50)に補正
func TestListBreweriesUsecase_Execute_LimitValidation_LessThan1(t *testing.T) {
	mockQuery := new(MockSakeQuery)
	uc := NewListBreweriesUsecase(mockQuery)

	mockQuery.On("ListBreweries", mock.Anything, (*string)(nil), int32(50)).Return([]entity.Brewery{}, nil)

	input := ListBreweriesInput{
		Limit: 0,
	}

	output, err := uc.Execute(context.Background(), input)

	assert.NoError(t, err)
	assert.NotNil(t, output)
	mockQuery.AssertExpectations(t)
}

// テスト: limitが100超の場合はデフォルト値(50)に補正
func TestListBreweriesUsecase_Execute_LimitValidation_GreaterThan100(t *testing.T) {
	mockQuery := new(MockSakeQuery)
	uc := NewListBreweriesUsecase(mockQuery)

	mockQuery.On("ListBreweries", mock.Anything, (*string)(nil), int32(50)).Return([]entity.Brewery{}, nil)

	input := ListBreweriesInput{
		Limit: 101,
	}

	output, err := uc.Execute(context.Background(), input)

	assert.NoError(t, err)
	assert.NotNil(t, output)
	mockQuery.AssertExpectations(t)
}

// テスト: クエリエラー発生時
func TestListBreweriesUsecase_Execute_QueryError(t *testing.T) {
	mockQuery := new(MockSakeQuery)
	uc := NewListBreweriesUsecase(mockQuery)

	expectedErr := errors.New("database error")
	mockQuery.On("ListBreweries", mock.Anything, mock.Anything, mock.Anything).Return(nil, expectedErr)

	input := ListBreweriesInput{
		Limit: 50,
	}

	output, err := uc.Execute(context.Background(), input)

	assert.Error(t, err)
	assert.Nil(t, output)
	assert.Equal(t, expectedErr, err)
	mockQuery.AssertExpectations(t)
}

// テスト: 空の結果が返る場合
func TestListBreweriesUsecase_Execute_EmptyResult(t *testing.T) {
	mockQuery := new(MockSakeQuery)
	uc := NewListBreweriesUsecase(mockQuery)

	mockQuery.On("ListBreweries", mock.Anything, (*string)(nil), int32(50)).Return([]entity.Brewery{}, nil)

	input := ListBreweriesInput{
		Limit: 50,
	}

	output, err := uc.Execute(context.Background(), input)

	assert.NoError(t, err)
	assert.NotNil(t, output)
	assert.Len(t, output.Breweries, 0)
	mockQuery.AssertExpectations(t)
}
