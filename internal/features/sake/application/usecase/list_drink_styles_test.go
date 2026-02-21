package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/sake-kasu/sake-hack-backend/internal/features/sake/domain/entity"
)

// テスト: 飲み方一覧の正常取得
func TestListDrinkStylesUsecase_Execute_Success(t *testing.T) {
	mockQuery := new(MockSakeQuery)
	uc := NewListDrinkStylesUsecase(mockQuery)

	desc1 := "5-10度に冷やして"
	desc2 := "50度前後に温めて"
	expectedStyles := []entity.DrinkStyle{
		{ID: 1, Name: "冷酒", Description: &desc1},
		{ID: 2, Name: "常温", Description: nil},
		{ID: 3, Name: "熱燗", Description: &desc2},
	}

	mockQuery.On("ListDrinkStyles", mock.Anything).Return(expectedStyles, nil)

	output, err := uc.Execute(context.Background())

	assert.NoError(t, err)
	assert.NotNil(t, output)
	assert.Len(t, output.DrinkStyles, 3)
	assert.Equal(t, int32(1), output.DrinkStyles[0].ID)
	assert.Equal(t, "冷酒", output.DrinkStyles[0].Name)
	assert.Equal(t, &desc1, output.DrinkStyles[0].Description)
	assert.Nil(t, output.DrinkStyles[1].Description)
	mockQuery.AssertExpectations(t)
}

// テスト: クエリエラー発生時
func TestListDrinkStylesUsecase_Execute_QueryError(t *testing.T) {
	mockQuery := new(MockSakeQuery)
	uc := NewListDrinkStylesUsecase(mockQuery)

	expectedErr := errors.New("database error")
	mockQuery.On("ListDrinkStyles", mock.Anything).Return(nil, expectedErr)

	output, err := uc.Execute(context.Background())

	assert.Error(t, err)
	assert.Nil(t, output)
	assert.Equal(t, expectedErr, err)
	mockQuery.AssertExpectations(t)
}

// テスト: 空の結果が返る場合
func TestListDrinkStylesUsecase_Execute_EmptyResult(t *testing.T) {
	mockQuery := new(MockSakeQuery)
	uc := NewListDrinkStylesUsecase(mockQuery)

	mockQuery.On("ListDrinkStyles", mock.Anything).Return([]entity.DrinkStyle{}, nil)

	output, err := uc.Execute(context.Background())

	assert.NoError(t, err)
	assert.NotNil(t, output)
	assert.Len(t, output.DrinkStyles, 0)
	mockQuery.AssertExpectations(t)
}
