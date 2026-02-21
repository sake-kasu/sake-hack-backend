package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/sake-kasu/sake-hack-backend/internal/features/sake/domain/entity"
)

// テスト: 酒の種類一覧の正常取得
func TestListKindsUsecase_Execute_Success(t *testing.T) {
	mockQuery := new(MockSakeQuery)
	uc := NewListKindsUsecase(mockQuery)

	expectedKinds := []entity.SakeKind{
		{ID: 1, Name: "純米大吟醸"},
		{ID: 2, Name: "大吟醸"},
		{ID: 3, Name: "純米吟醸"},
	}

	mockQuery.On("ListKinds", mock.Anything).Return(expectedKinds, nil)

	output, err := uc.Execute(context.Background())

	assert.NoError(t, err)
	assert.NotNil(t, output)
	assert.Len(t, output.Kinds, 3)
	assert.Equal(t, int32(1), output.Kinds[0].ID)
	assert.Equal(t, "純米大吟醸", output.Kinds[0].Name)
	mockQuery.AssertExpectations(t)
}

// テスト: クエリエラー発生時
func TestListKindsUsecase_Execute_QueryError(t *testing.T) {
	mockQuery := new(MockSakeQuery)
	uc := NewListKindsUsecase(mockQuery)

	expectedErr := errors.New("database error")
	mockQuery.On("ListKinds", mock.Anything).Return(nil, expectedErr)

	output, err := uc.Execute(context.Background())

	assert.Error(t, err)
	assert.Nil(t, output)
	assert.Equal(t, expectedErr, err)
	mockQuery.AssertExpectations(t)
}

// テスト: 空の結果が返る場合
func TestListKindsUsecase_Execute_EmptyResult(t *testing.T) {
	mockQuery := new(MockSakeQuery)
	uc := NewListKindsUsecase(mockQuery)

	mockQuery.On("ListKinds", mock.Anything).Return([]entity.SakeKind{}, nil)

	output, err := uc.Execute(context.Background())

	assert.NoError(t, err)
	assert.NotNil(t, output)
	assert.Len(t, output.Kinds, 0)
	mockQuery.AssertExpectations(t)
}
