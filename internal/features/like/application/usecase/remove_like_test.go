package usecase

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/sake-kasu/sake-hack-backend/internal/apperror"
)

// テスト: いいね削除の正常処理
func TestRemoveLikeUsecase_Execute_Success(t *testing.T) {
	mockRepo := new(MockLikeRepository)
	uc := NewRemoveLikeUsecase(mockRepo)

	input := RemoveLikeInput{SakeID: 1, Token: "test-token"}

	mockRepo.On("Delete", mock.Anything, int32(1), "test-token").Return(int64(1), nil)
	mockRepo.On("GetCount", mock.Anything, int32(1)).Return(int64(41), nil)

	output, err := uc.Execute(context.Background(), input)

	assert.NoError(t, err)
	assert.NotNil(t, output)
	assert.Equal(t, int64(41), output.LikeCount)
	assert.False(t, output.IsLiked)
	mockRepo.AssertExpectations(t)
}

// テスト: いいねが存在しない場合のNotFoundエラー
func TestRemoveLikeUsecase_Execute_NotFound(t *testing.T) {
	mockRepo := new(MockLikeRepository)
	uc := NewRemoveLikeUsecase(mockRepo)

	input := RemoveLikeInput{SakeID: 1, Token: "test-token"}

	// rowsAffected == 0
	mockRepo.On("Delete", mock.Anything, int32(1), "test-token").Return(int64(0), nil)

	output, err := uc.Execute(context.Background(), input)

	assert.Error(t, err)
	assert.Nil(t, output)

	var appErr *apperror.AppError
	assert.ErrorAs(t, err, &appErr)
	assert.Equal(t, apperror.ErrCodeNotFound, appErr.Code)

	mockRepo.AssertExpectations(t)
	mockRepo.AssertNotCalled(t, "GetCount")
}

// テスト: Delete時にDBエラーが発生した場合
func TestRemoveLikeUsecase_Execute_DeleteError(t *testing.T) {
	mockRepo := new(MockLikeRepository)
	uc := NewRemoveLikeUsecase(mockRepo)

	input := RemoveLikeInput{SakeID: 1, Token: "test-token"}

	mockRepo.On("Delete", mock.Anything, int32(1), "test-token").Return(int64(0),
		apperror.DatabaseError("削除エラー", fmt.Errorf("delete error")),
	)

	output, err := uc.Execute(context.Background(), input)

	assert.Error(t, err)
	assert.Nil(t, output)
	mockRepo.AssertExpectations(t)
	mockRepo.AssertNotCalled(t, "GetCount")
}

// テスト: GetCount時にDBエラーが発生した場合
func TestRemoveLikeUsecase_Execute_GetCountError(t *testing.T) {
	mockRepo := new(MockLikeRepository)
	uc := NewRemoveLikeUsecase(mockRepo)

	input := RemoveLikeInput{SakeID: 1, Token: "test-token"}

	mockRepo.On("Delete", mock.Anything, int32(1), "test-token").Return(int64(1), nil)
	mockRepo.On("GetCount", mock.Anything, int32(1)).Return(int64(0),
		apperror.DatabaseError("カウント取得エラー", fmt.Errorf("count error")),
	)

	output, err := uc.Execute(context.Background(), input)

	assert.Error(t, err)
	assert.Nil(t, output)
	mockRepo.AssertExpectations(t)
}

// テスト: いいね数が0になるまで削除した場合
func TestRemoveLikeUsecase_Execute_Success_ZeroCount(t *testing.T) {
	mockRepo := new(MockLikeRepository)
	uc := NewRemoveLikeUsecase(mockRepo)

	input := RemoveLikeInput{SakeID: 1, Token: "test-token"}

	mockRepo.On("Delete", mock.Anything, int32(1), "test-token").Return(int64(1), nil)
	mockRepo.On("GetCount", mock.Anything, int32(1)).Return(int64(0), nil)

	output, err := uc.Execute(context.Background(), input)

	assert.NoError(t, err)
	assert.NotNil(t, output)
	assert.Equal(t, int64(0), output.LikeCount)
	assert.False(t, output.IsLiked)
	mockRepo.AssertExpectations(t)
}
