package presentation

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/sake-kasu/sake-hack-backend/api/generated"
	"github.com/sake-kasu/sake-hack-backend/internal/apperror"
	"github.com/sake-kasu/sake-hack-backend/internal/features/like/application/usecase"
)

// --- Mock Usecases ---

// MockToggleLikeUsecase いいね追加のモック
type MockToggleLikeUsecase struct {
	mock.Mock
}

func (m *MockToggleLikeUsecase) Execute(ctx context.Context, input usecase.ToggleLikeInput) (*usecase.ToggleLikeOutput, error) {
	args := m.Called(ctx, input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*usecase.ToggleLikeOutput), args.Error(1)
}

// MockRemoveLikeUsecase いいね削除のモック
type MockRemoveLikeUsecase struct {
	mock.Mock
}

func (m *MockRemoveLikeUsecase) Execute(ctx context.Context, input usecase.RemoveLikeInput) (*usecase.RemoveLikeOutput, error) {
	args := m.Called(ctx, input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*usecase.RemoveLikeOutput), args.Error(1)
}

// --- Helper ---

// stubSakeServer はテスト用のスタブ(Sake Featureのメソッドを空実装)
type stubSakeServer struct{}

func (s *stubSakeServer) ListBreweries(_ *gin.Context, _ generated.ListBreweriesParams) {}
func (s *stubSakeServer) ListDrinkStyles(_ *gin.Context)                                {}
func (s *stubSakeServer) ListKinds(_ *gin.Context)                                      {}
func (s *stubSakeServer) ListSakes(_ *gin.Context, _ generated.ListSakesParams)         {}
func (s *stubSakeServer) GetSakeDetail(_ *gin.Context, _ int32, _ generated.GetSakeDetailParams) {
}
func (s *stubSakeServer) ListStocks(_ *gin.Context, _ generated.ListStocksParams) {}
func (s *stubSakeServer) CreateStock(_ *gin.Context)                              {}
func (s *stubSakeServer) DeleteStock(_ *gin.Context, _ int32)                     {}
func (s *stubSakeServer) GetStockDetail(_ *gin.Context, _ int32)                  {}
func (s *stubSakeServer) PatchStock(_ *gin.Context, _ int32)                      {}
func (s *stubSakeServer) UpdateStock(_ *gin.Context, _ int32)                     {}
func (s *stubSakeServer) CreateStockUploadUrl(_ *gin.Context, _ int32)            {}

// testCompositeServer テスト用のcompositeServer
type testCompositeServer struct {
	*stubSakeServer
	*LikeServerImpl
}

func newTestServer() (*LikeServerImpl, *MockToggleLikeUsecase, *MockRemoveLikeUsecase) {
	toggleUC := new(MockToggleLikeUsecase)
	removeUC := new(MockRemoveLikeUsecase)
	server := NewLikeServerImpl(toggleUC, removeUC)
	return server, toggleUC, removeUC
}

func newRouter(server *LikeServerImpl) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	composite := &testCompositeServer{
		stubSakeServer: &stubSakeServer{},
		LikeServerImpl: server,
	}
	generated.RegisterHandlers(router, composite)
	return router
}

// --- CreateSakeLike Tests ---

// テスト: いいね追加の正常処理
func TestCreateSakeLike_Success(t *testing.T) {
	server, toggleUC, _ := newTestServer()

	toggleUC.On("Execute", mock.Anything, usecase.ToggleLikeInput{
		SakeID: 1,
		Token:  "test-token-uuid",
	}).Return(&usecase.ToggleLikeOutput{
		LikeCount: 42,
		IsLiked:   true,
	}, nil)

	router := newRouter(server)
	req := httptest.NewRequest(http.MethodPost, "/sakes/1/likes", nil)
	req.Header.Set("X-Like-Token", "test-token-uuid")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response generated.LikeResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, int64(42), response.LikeCount)
	assert.True(t, response.IsLiked)
	toggleUC.AssertExpectations(t)
}

// テスト: X-Like-Tokenヘッダーが空の場合のバリデーションエラー
func TestCreateSakeLike_ValidationError_EmptyToken(t *testing.T) {
	server, toggleUC, _ := newTestServer()

	router := newRouter(server)
	req := httptest.NewRequest(http.MethodPost, "/sakes/1/likes", nil)
	req.Header.Set("X-Like-Token", "")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	toggleUC.AssertNotCalled(t, "Execute")
}

// テスト: X-Like-Tokenヘッダーが長すぎる場合のバリデーションエラー
func TestCreateSakeLike_ValidationError_TokenTooLong(t *testing.T) {
	server, toggleUC, _ := newTestServer()

	// 256文字のトークン
	longToken := ""
	for i := 0; i < 256; i++ {
		longToken += "a"
	}

	router := newRouter(server)
	req := httptest.NewRequest(http.MethodPost, "/sakes/1/likes", nil)
	req.Header.Set("X-Like-Token", longToken)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	toggleUC.AssertNotCalled(t, "Execute")
}

// テスト: いいね追加時のユースケースエラー
func TestCreateSakeLike_UsecaseError(t *testing.T) {
	server, toggleUC, _ := newTestServer()

	toggleUC.On("Execute", mock.Anything, mock.Anything).Return(
		nil,
		apperror.DatabaseError("データベースエラー", fmt.Errorf("connection error")),
	)

	router := newRouter(server)
	req := httptest.NewRequest(http.MethodPost, "/sakes/1/likes", nil)
	req.Header.Set("X-Like-Token", "test-token-uuid")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var response generated.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.NotNil(t, response.Errors)
	assert.Equal(t, apperror.ErrCodeDatabaseError, (*response.Errors)[0].Code)
	toggleUC.AssertExpectations(t)
}

// --- DeleteSakeLike Tests ---

// テスト: いいね削除の正常処理
func TestDeleteSakeLike_Success(t *testing.T) {
	server, _, removeUC := newTestServer()

	removeUC.On("Execute", mock.Anything, usecase.RemoveLikeInput{
		SakeID: 1,
		Token:  "test-token-uuid",
	}).Return(&usecase.RemoveLikeOutput{
		LikeCount: 41,
		IsLiked:   false,
	}, nil)

	router := newRouter(server)
	req := httptest.NewRequest(http.MethodDelete, "/sakes/1/likes", nil)
	req.Header.Set("X-Like-Token", "test-token-uuid")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response generated.LikeResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, int64(41), response.LikeCount)
	assert.False(t, response.IsLiked)
	removeUC.AssertExpectations(t)
}

// テスト: いいね削除時にいいねが存在しない場合
func TestDeleteSakeLike_NotFound(t *testing.T) {
	server, _, removeUC := newTestServer()

	removeUC.On("Execute", mock.Anything, mock.Anything).Return(
		nil,
		apperror.NotFoundError("いいねが見つかりません"),
	)

	router := newRouter(server)
	req := httptest.NewRequest(http.MethodDelete, "/sakes/1/likes", nil)
	req.Header.Set("X-Like-Token", "test-token-uuid")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	removeUC.AssertExpectations(t)
}

// テスト: いいね削除時のX-Like-Tokenヘッダー空
func TestDeleteSakeLike_ValidationError_EmptyToken(t *testing.T) {
	server, _, removeUC := newTestServer()

	router := newRouter(server)
	req := httptest.NewRequest(http.MethodDelete, "/sakes/1/likes", nil)
	req.Header.Set("X-Like-Token", "")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	removeUC.AssertNotCalled(t, "Execute")
}

// --- Error Handler Tests ---

// テスト: ValidationErrorのHTTPレスポンス変換
func TestHandleError_ValidationError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	verr := apperror.NewValidationError("バリデーションエラー").
		AddField("field1", "エラーメッセージ")

	handleError(c, verr)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response generated.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.NotNil(t, response.Errors)
	assert.Equal(t, apperror.ErrCodeValidationError, (*response.Errors)[0].Code)
}

// テスト: AppErrorのHTTPレスポンス変換
func TestHandleError_AppError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	handleError(c, apperror.NotFoundError("見つかりません"))

	assert.Equal(t, http.StatusNotFound, w.Code)
}

// テスト: 予期しないエラーのHTTPレスポンス変換
func TestHandleError_UnexpectedError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	handleError(c, fmt.Errorf("unexpected error"))

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

// --- Validator Tests ---

// テスト: 有効なトークンでバリデーション成功
func TestValidateLikeToken_Success(t *testing.T) {
	err := validateLikeToken("valid-token-uuid-v4")
	assert.NoError(t, err)
}

// テスト: 空トークンでバリデーション失敗
func TestValidateLikeToken_ValidationError_Empty(t *testing.T) {
	err := validateLikeToken("")
	assert.Error(t, err)
}

// テスト: 最大長ちょうどのトークンでバリデーション成功
func TestValidateLikeToken_BoundaryValue_MaxLength(t *testing.T) {
	token := ""
	for i := 0; i < 255; i++ {
		token += "a"
	}
	err := validateLikeToken(token)
	assert.NoError(t, err)
}

// テスト: 最大長+1のトークンでバリデーション失敗
func TestValidateLikeToken_BoundaryValue_MaxLengthPlus1(t *testing.T) {
	token := ""
	for i := 0; i < 256; i++ {
		token += "a"
	}
	err := validateLikeToken(token)
	assert.Error(t, err)
}
