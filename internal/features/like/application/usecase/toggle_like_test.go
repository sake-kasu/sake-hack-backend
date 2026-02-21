package usecase

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/sake-kasu/sake-hack-backend/internal/apperror"
)

// MockLikeRepository いいねリポジトリのモック
type MockLikeRepository struct {
	mock.Mock
}

func (m *MockLikeRepository) Create(ctx context.Context, sakeID int32, token string) error {
	args := m.Called(ctx, sakeID, token)
	return args.Error(0)
}

func (m *MockLikeRepository) Delete(ctx context.Context, sakeID int32, token string) (int64, error) {
	args := m.Called(ctx, sakeID, token)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockLikeRepository) GetCount(ctx context.Context, sakeID int32) (int64, error) {
	args := m.Called(ctx, sakeID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockLikeRepository) Exists(ctx context.Context, sakeID int32, token string) (bool, error) {
	args := m.Called(ctx, sakeID, token)
	return args.Bool(0), args.Error(1)
}

// テスト: いいね追加の正常処理
func TestToggleLikeUsecase_Execute_Success(t *testing.T) {
	mockRepo := new(MockLikeRepository)
	uc := NewToggleLikeUsecase(mockRepo)

	input := ToggleLikeInput{SakeID: 1, Token: "test-token"}

	mockRepo.On("Create", mock.Anything, int32(1), "test-token").Return(nil)
	mockRepo.On("GetCount", mock.Anything, int32(1)).Return(int64(42), nil)
	mockRepo.On("Exists", mock.Anything, int32(1), "test-token").Return(true, nil)

	output, err := uc.Execute(context.Background(), input)

	assert.NoError(t, err)
	assert.NotNil(t, output)
	assert.Equal(t, int64(42), output.LikeCount)
	assert.True(t, output.IsLiked)
	mockRepo.AssertExpectations(t)
}

// テスト: 既にいいね済み(ON CONFLICT DO NOTHING)でも正常に返る
func TestToggleLikeUsecase_Execute_Success_AlreadyLiked(t *testing.T) {
	mockRepo := new(MockLikeRepository)
	uc := NewToggleLikeUsecase(mockRepo)

	input := ToggleLikeInput{SakeID: 1, Token: "test-token"}

	// ON CONFLICT DO NOTHINGでもエラーにならない
	mockRepo.On("Create", mock.Anything, int32(1), "test-token").Return(nil)
	mockRepo.On("GetCount", mock.Anything, int32(1)).Return(int64(10), nil)
	mockRepo.On("Exists", mock.Anything, int32(1), "test-token").Return(true, nil)

	output, err := uc.Execute(context.Background(), input)

	assert.NoError(t, err)
	assert.NotNil(t, output)
	assert.Equal(t, int64(10), output.LikeCount)
	assert.True(t, output.IsLiked)
	mockRepo.AssertExpectations(t)
}

// テスト: Create時にDBエラーが発生した場合
func TestToggleLikeUsecase_Execute_CreateError(t *testing.T) {
	mockRepo := new(MockLikeRepository)
	uc := NewToggleLikeUsecase(mockRepo)

	input := ToggleLikeInput{SakeID: 1, Token: "test-token"}

	mockRepo.On("Create", mock.Anything, int32(1), "test-token").Return(
		apperror.DatabaseError("データベースエラー", fmt.Errorf("connection error")),
	)

	output, err := uc.Execute(context.Background(), input)

	assert.Error(t, err)
	assert.Nil(t, output)
	mockRepo.AssertExpectations(t)
	mockRepo.AssertNotCalled(t, "GetCount")
	mockRepo.AssertNotCalled(t, "Exists")
}

// テスト: GetCount時にDBエラーが発生した場合
func TestToggleLikeUsecase_Execute_GetCountError(t *testing.T) {
	mockRepo := new(MockLikeRepository)
	uc := NewToggleLikeUsecase(mockRepo)

	input := ToggleLikeInput{SakeID: 1, Token: "test-token"}

	mockRepo.On("Create", mock.Anything, int32(1), "test-token").Return(nil)
	mockRepo.On("GetCount", mock.Anything, int32(1)).Return(int64(0),
		apperror.DatabaseError("カウント取得エラー", fmt.Errorf("count error")),
	)

	output, err := uc.Execute(context.Background(), input)

	assert.Error(t, err)
	assert.Nil(t, output)
	mockRepo.AssertExpectations(t)
	mockRepo.AssertNotCalled(t, "Exists")
}

// テスト: Exists時にDBエラーが発生した場合
func TestToggleLikeUsecase_Execute_ExistsError(t *testing.T) {
	mockRepo := new(MockLikeRepository)
	uc := NewToggleLikeUsecase(mockRepo)

	input := ToggleLikeInput{SakeID: 1, Token: "test-token"}

	mockRepo.On("Create", mock.Anything, int32(1), "test-token").Return(nil)
	mockRepo.On("GetCount", mock.Anything, int32(1)).Return(int64(42), nil)
	mockRepo.On("Exists", mock.Anything, int32(1), "test-token").Return(false,
		apperror.DatabaseError("存在確認エラー", fmt.Errorf("exists error")),
	)

	output, err := uc.Execute(context.Background(), input)

	assert.Error(t, err)
	assert.Nil(t, output)
	mockRepo.AssertExpectations(t)
}
