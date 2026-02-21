package presentation

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/sake-kasu/sake-hack-backend/api/generated"
	"github.com/sake-kasu/sake-hack-backend/internal/apperror"
	"github.com/sake-kasu/sake-hack-backend/internal/features/sake/application/usecase"
	"github.com/sake-kasu/sake-hack-backend/internal/features/sake/domain/entity"
)

// --- Mock Usecases ---

type MockListSakesUsecase struct {
	mock.Mock
}

func (m *MockListSakesUsecase) Execute(ctx context.Context, input usecase.ListSakesInput) (*usecase.ListSakesOutput, error) {
	args := m.Called(ctx, input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*usecase.ListSakesOutput), args.Error(1)
}

type MockGetSakeDetailUsecase struct {
	mock.Mock
}

func (m *MockGetSakeDetailUsecase) Execute(ctx context.Context, id int32) (*usecase.GetSakeDetailOutput, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*usecase.GetSakeDetailOutput), args.Error(1)
}

type MockCreateStockUsecase struct {
	mock.Mock
}

func (m *MockCreateStockUsecase) Execute(ctx context.Context, input usecase.CreateStockInput) (*usecase.CreateStockOutput, error) {
	args := m.Called(ctx, input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*usecase.CreateStockOutput), args.Error(1)
}

type MockUpdateStockUsecase struct {
	mock.Mock
}

func (m *MockUpdateStockUsecase) Execute(ctx context.Context, input usecase.UpdateStockInput) (*usecase.UpdateStockOutput, error) {
	args := m.Called(ctx, input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*usecase.UpdateStockOutput), args.Error(1)
}

type MockDeleteStockUsecase struct {
	mock.Mock
}

func (m *MockDeleteStockUsecase) Execute(ctx context.Context, id int32) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// MockPatchStockUsecase は在庫部分更新のモック
type MockPatchStockUsecase struct {
	mock.Mock
}

func (m *MockPatchStockUsecase) Execute(ctx context.Context, input usecase.PatchStockInput) error {
	args := m.Called(ctx, input)
	return args.Error(0)
}

// MockCreateUploadUrlUsecase は画像アップロードURL発行のモック
type MockCreateUploadUrlUsecase struct {
	mock.Mock
}

func (m *MockCreateUploadUrlUsecase) Execute(ctx context.Context, input usecase.CreateUploadUrlInput) (*usecase.CreateUploadUrlOutput, error) {
	args := m.Called(ctx, input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*usecase.CreateUploadUrlOutput), args.Error(1)
}

// MockURLResolver はURLResolverのモック
type MockURLResolver struct {
	mock.Mock
}

func (m *MockURLResolver) ResolveURL(ctx context.Context, objectKey string) (string, error) {
	args := m.Called(ctx, objectKey)
	return args.String(0), args.Error(1)
}

// MockListKindsUsecase は酒の種類一覧取得のモック
type MockListKindsUsecase struct {
	mock.Mock
}

func (m *MockListKindsUsecase) Execute(ctx context.Context) (*usecase.ListKindsOutput, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*usecase.ListKindsOutput), args.Error(1)
}

// MockListBreweriesUsecase は酒造一覧取得のモック
type MockListBreweriesUsecase struct {
	mock.Mock
}

func (m *MockListBreweriesUsecase) Execute(ctx context.Context, input usecase.ListBreweriesInput) (*usecase.ListBreweriesOutput, error) {
	args := m.Called(ctx, input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*usecase.ListBreweriesOutput), args.Error(1)
}

// MockListDrinkStylesUsecase は飲み方一覧取得のモック
type MockListDrinkStylesUsecase struct {
	mock.Mock
}

func (m *MockListDrinkStylesUsecase) Execute(ctx context.Context) (*usecase.ListDrinkStylesOutput, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*usecase.ListDrinkStylesOutput), args.Error(1)
}

// --- Helper ---

func newTestServer() (*SakeServerImpl, *MockListSakesUsecase, *MockGetSakeDetailUsecase, *MockCreateStockUsecase, *MockUpdateStockUsecase, *MockDeleteStockUsecase, *MockPatchStockUsecase, *MockCreateUploadUrlUsecase, *MockListKindsUsecase, *MockListBreweriesUsecase, *MockListDrinkStylesUsecase, *MockURLResolver) {
	listUC := new(MockListSakesUsecase)
	detailUC := new(MockGetSakeDetailUsecase)
	createUC := new(MockCreateStockUsecase)
	updateUC := new(MockUpdateStockUsecase)
	deleteUC := new(MockDeleteStockUsecase)
	patchUC := new(MockPatchStockUsecase)
	uploadUrlUC := new(MockCreateUploadUrlUsecase)
	listKindsUC := new(MockListKindsUsecase)
	listBreweriesUC := new(MockListBreweriesUsecase)
	listDrinkStylesUC := new(MockListDrinkStylesUsecase)
	urlResolver := new(MockURLResolver)
	server := NewSakeServerImpl(listUC, detailUC, createUC, updateUC, deleteUC, patchUC, uploadUrlUC, listKindsUC, listBreweriesUC, listDrinkStylesUC, urlResolver)
	return server, listUC, detailUC, createUC, updateUC, deleteUC, patchUC, uploadUrlUC, listKindsUC, listBreweriesUC, listDrinkStylesUC, urlResolver
}

func newRouter(server *SakeServerImpl) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	generated.RegisterHandlers(router, server)
	return router
}

// --- ListSakes Tests ---

// テスト: 酒一覧の正常取得(デフォルトパラメータ)
func TestListSakes_Success(t *testing.T) {
	server, listUC, _, _, _, _, _, _, _, _, _, urlResolver := newTestServer()
	objectKey := "sakes/abc123.jpg"

	// URLResolverのモック設定: object_keyから署名付きURLに変換
	urlResolver.On("ResolveURL", mock.Anything, objectKey).Return("https://signed-url.example.com/sakes/abc123.jpg", nil)

	expectedOutput := &usecase.ListSakesOutput{
		Sakes: []entity.SakeListItem{
			{
				ID:        1,
				Category:  entity.SakeCategoryJapaneseSake,
				Name:      "獺祭 純米大吟醸50",
				ObjectKey: &objectKey,
			},
		},
		Pagination: entity.Pagination{Total: 100, Offset: 0, Limit: 20},
	}

	listUC.On("Execute", mock.Anything, mock.MatchedBy(func(input usecase.ListSakesInput) bool {
		return input.Offset == 0 && input.Limit == 20
	})).Return(expectedOutput, nil)

	router := newRouter(server)
	req := httptest.NewRequest(http.MethodGet, "/sakes", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response generated.ListSakesResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.NotNil(t, response.Data)
	assert.Len(t, *response.Data, 1)
	assert.Equal(t, int32(1), (*response.Data)[0].Id)
	assert.Equal(t, "獺祭 純米大吟醸50", (*response.Data)[0].Name)
	assert.Equal(t, generated.SakeCategory("JAPANESE_SAKE"), (*response.Data)[0].Category)
	// imagePreviewには署名付きURLが返される
	assert.Equal(t, "https://signed-url.example.com/sakes/abc123.jpg", (*response.Data)[0].ImagePreview)
	assert.NotNil(t, response.Meta)
	assert.Equal(t, int64(100), response.Meta.Total)
	listUC.AssertExpectations(t)
	urlResolver.AssertExpectations(t)
}

// テスト: offsetが負数の場合のバリデーションエラー
func TestListSakes_ValidationError_OffsetLessThan0(t *testing.T) {
	server, listUC, _, _, _, _, _, _, _, _, _, _ := newTestServer()
	router := newRouter(server)

	req := httptest.NewRequest(http.MethodGet, "/sakes?offset=-1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	listUC.AssertNotCalled(t, "Execute")
}

// テスト: limitの範囲外(0以下、100超)
func TestListSakes_ValidationError_LimitOutOfRange(t *testing.T) {
	tests := []struct {
		name  string
		query string
	}{
		{"LimitLessThan1", "/sakes?limit=0"},
		{"LimitGreaterThan100", "/sakes?limit=101"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server, listUC, _, _, _, _, _, _, _, _, _, _ := newTestServer()
			router := newRouter(server)

			req := httptest.NewRequest(http.MethodGet, tt.query, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusBadRequest, w.Code)
			listUC.AssertNotCalled(t, "Execute")
		})
	}
}

// テスト: typeIdが0以下の場合のバリデーションエラー
func TestListSakes_ValidationError_TypeIdLessThan1(t *testing.T) {
	server, listUC, _, _, _, _, _, _, _, _, _, _ := newTestServer()
	router := newRouter(server)

	req := httptest.NewRequest(http.MethodGet, "/sakes?typeId=0", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	listUC.AssertNotCalled(t, "Execute")
}

// テスト: breweryIdが0以下の場合のバリデーションエラー
func TestListSakes_ValidationError_BreweryIdLessThan1(t *testing.T) {
	server, listUC, _, _, _, _, _, _, _, _, _, _ := newTestServer()
	router := newRouter(server)

	req := httptest.NewRequest(http.MethodGet, "/sakes?breweryId=0", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	listUC.AssertNotCalled(t, "Execute")
}

// テスト: ユースケースエラー発生時のレスポンス
func TestListSakes_UsecaseError(t *testing.T) {
	server, listUC, _, _, _, _, _, _, _, _, _, _ := newTestServer()

	listUC.On("Execute", mock.Anything, mock.Anything).Return(
		nil,
		apperror.DatabaseError("データベースエラー", fmt.Errorf("connection error")),
	)

	router := newRouter(server)
	req := httptest.NewRequest(http.MethodGet, "/sakes", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var response generated.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.NotNil(t, response.Errors)
	assert.Equal(t, apperror.ErrCodeDatabaseError, (*response.Errors)[0].Code)
	listUC.AssertExpectations(t)
}

// テスト: デフォルトパラメータの確認(offset=0, limit=20)
func TestListSakes_DefaultValues(t *testing.T) {
	server, listUC, _, _, _, _, _, _, _, _, _, _ := newTestServer()

	listUC.On("Execute", mock.Anything, mock.MatchedBy(func(input usecase.ListSakesInput) bool {
		return input.Offset == 0 && input.Limit == 20 && input.KindID == nil && input.BreweryID == nil
	})).Return(&usecase.ListSakesOutput{
		Sakes:      []entity.SakeListItem{},
		Pagination: entity.Pagination{Total: 0, Offset: 0, Limit: 20},
	}, nil)

	router := newRouter(server)
	req := httptest.NewRequest(http.MethodGet, "/sakes", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	listUC.AssertExpectations(t)
}

// --- GetSakeDetail Tests ---

// テスト: 酒詳細の正常取得
func TestGetSakeDetail_Success(t *testing.T) {
	server, _, detailUC, _, _, _, _, _, _, _, _, urlResolver := newTestServer()
	now := time.Now().Truncate(time.Second)
	description := "5-10度に冷やして"
	objectKey := "sakes/detail123.jpg"

	// URLResolverのモック設定
	urlResolver.On("ResolveURL", mock.Anything, objectKey).Return("https://signed-url.example.com/sakes/detail123.jpg", nil)

	detailUC.On("Execute", mock.Anything, int32(1)).Return(&usecase.GetSakeDetailOutput{
		Detail: &entity.SakeDetail{
			ID:       1,
			Category: entity.SakeCategoryJapaneseSake,
			Kind:     entity.SakeKind{ID: 1, Name: "純米大吟醸"},
			Brewery: entity.Brewery{
				ID: 1, Name: "旭酒造", OriginCountry: "日本",
			},
			Name:            entity.SakeName{Name: "獺祭", Phonetic: "だっさい"},
			Abv:             16.0,
			PurchaseVolume:  720,
			RemainingVolume: 500,
			Price:           3000,
			DrinkStyles: []entity.DrinkStyle{
				{ID: 1, Name: "冷酒", Description: &description},
			},
			ObjectKey: &objectKey,
			CreatedAt: now,
			UpdatedAt: now,
		},
	}, nil)

	router := newRouter(server)
	req := httptest.NewRequest(http.MethodGet, "/sakes/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response generated.SakeDetail
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, int32(1), response.Id)
	assert.Equal(t, generated.SakeCategory("JAPANESE_SAKE"), response.Category)
	assert.Equal(t, "獺祭", response.Name.Name)
	assert.Equal(t, float32(16.0), response.Abv)
	assert.Len(t, response.DrinkStyles, 1)
	// imageUrlには署名付きURLが返される
	assert.NotNil(t, response.ImageUrl)
	assert.Equal(t, "https://signed-url.example.com/sakes/detail123.jpg", *response.ImageUrl)
	detailUC.AssertExpectations(t)
	urlResolver.AssertExpectations(t)
}

// テスト: 存在しない酒IDの場合
func TestGetSakeDetail_NotFound(t *testing.T) {
	server, _, detailUC, _, _, _, _, _, _, _, _, _ := newTestServer()

	detailUC.On("Execute", mock.Anything, int32(999)).Return(
		nil,
		apperror.NotFoundError("酒が見つかりません"),
	)

	router := newRouter(server)
	req := httptest.NewRequest(http.MethodGet, "/sakes/999", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	detailUC.AssertExpectations(t)
}

// --- ListStocks Tests ---

// テスト: 在庫一覧の正常取得
func TestListStocks_Success(t *testing.T) {
	server, listUC, _, _, _, _, _, _, _, _, _, _ := newTestServer()

	listUC.On("Execute", mock.Anything, mock.MatchedBy(func(input usecase.ListSakesInput) bool {
		return input.Offset == 0 && input.Limit == 20
	})).Return(&usecase.ListSakesOutput{
		Sakes:      []entity.SakeListItem{},
		Pagination: entity.Pagination{Total: 0, Offset: 0, Limit: 20},
	}, nil)

	router := newRouter(server)
	req := httptest.NewRequest(http.MethodGet, "/stocks", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	listUC.AssertExpectations(t)
}

// --- CreateStock Tests ---

// テスト: 在庫登録の正常処理
func TestCreateStock_Success(t *testing.T) {
	server, _, _, createUC, _, _, _, _, _, _, _, _ := newTestServer()

	createUC.On("Execute", mock.Anything, mock.Anything).Return(&usecase.CreateStockOutput{
		Sake: &entity.SakeListItem{
			ID:       1,
			Category: entity.SakeCategoryJapaneseSake,
			Name:     "獺祭",
		},
	}, nil)

	body := generated.CreateSakeRequest{
		Category:        generated.SakeCategoryJAPANESESAKE,
		Kind:            generated.SakeKind{Id: 0, Name: "純米大吟醸"},
		Brewery:         generated.Brewery{Id: 0, Name: "旭酒造", OriginCountry: "日本"},
		Name:            generated.SakeName{Name: "獺祭", Phonetic: "だっさい"},
		Abv:             16.0,
		PurchaseVolume:  720,
		RemainingVolume: 500,
		Price:           3000,
		DrinkStyles:     []generated.DrinkStyle{},
	}
	bodyBytes, _ := json.Marshal(body)

	router := newRouter(server)
	req := httptest.NewRequest(http.MethodPost, "/stocks", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response generated.CreateSakeResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.NotNil(t, response.Data)
	assert.Equal(t, int32(1), response.Data.Id)
	createUC.AssertExpectations(t)
}

// テスト: 酒名欠如時のバリデーションエラー
func TestCreateStock_ValidationError_MissingName(t *testing.T) {
	server, _, _, createUC, _, _, _, _, _, _, _, _ := newTestServer()

	body := generated.CreateSakeRequest{
		Category:        generated.SakeCategoryJAPANESESAKE,
		Kind:            generated.SakeKind{Name: "純米"},
		Brewery:         generated.Brewery{Name: "旭酒造", OriginCountry: "日本"},
		Name:            generated.SakeName{Name: "", Phonetic: ""},
		Abv:             16.0,
		PurchaseVolume:  720,
		RemainingVolume: 500,
		Price:           3000,
		DrinkStyles:     []generated.DrinkStyle{},
	}
	bodyBytes, _ := json.Marshal(body)

	router := newRouter(server)
	req := httptest.NewRequest(http.MethodPost, "/stocks", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	createUC.AssertNotCalled(t, "Execute")
}

// テスト: 不正JSONの場合
func TestCreateStock_InvalidJSON(t *testing.T) {
	server, _, _, createUC, _, _, _, _, _, _, _, _ := newTestServer()

	router := newRouter(server)
	req := httptest.NewRequest(http.MethodPost, "/stocks", bytes.NewReader([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	createUC.AssertNotCalled(t, "Execute")
}

// --- GetStockDetail Tests ---

// テスト: 在庫詳細の正常取得
func TestGetStockDetail_Success(t *testing.T) {
	server, _, detailUC, _, _, _, _, _, _, _, _, _ := newTestServer()
	now := time.Now().Truncate(time.Second)

	detailUC.On("Execute", mock.Anything, int32(1)).Return(&usecase.GetSakeDetailOutput{
		Detail: &entity.SakeDetail{
			ID:              1,
			Category:        entity.SakeCategoryWhisky,
			Kind:            entity.SakeKind{ID: 2, Name: "シングルモルト"},
			Brewery:         entity.Brewery{ID: 2, Name: "サントリー", OriginCountry: "日本"},
			Name:            entity.SakeName{Name: "山崎", Phonetic: "やまざき"},
			Abv:             43.0,
			PurchaseVolume:  700,
			RemainingVolume: 350,
			Price:           15000,
			DrinkStyles:     []entity.DrinkStyle{},
			CreatedAt:       now,
			UpdatedAt:       now,
		},
	}, nil)

	router := newRouter(server)
	req := httptest.NewRequest(http.MethodGet, "/stocks/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response generated.SakeDetail
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, int32(1), response.Id)
	assert.Equal(t, generated.SakeCategory("WHISKY"), response.Category)
	detailUC.AssertExpectations(t)
}

// --- UpdateStock Tests ---

// テスト: 在庫更新の正常処理
func TestUpdateStock_Success(t *testing.T) {
	server, _, detailUC, _, updateUC, _, _, _, _, _, _, _ := newTestServer()
	now := time.Now().Truncate(time.Second)

	updateUC.On("Execute", mock.Anything, mock.Anything).Return(&usecase.UpdateStockOutput{
		Sake: &entity.SakeListItem{ID: 1, Category: entity.SakeCategoryJapaneseSake, Name: "獺祭"},
	}, nil)

	detailUC.On("Execute", mock.Anything, int32(1)).Return(&usecase.GetSakeDetailOutput{
		Detail: &entity.SakeDetail{
			ID:              1,
			Category:        entity.SakeCategoryJapaneseSake,
			Kind:            entity.SakeKind{ID: 1, Name: "純米大吟醸"},
			Brewery:         entity.Brewery{ID: 1, Name: "旭酒造", OriginCountry: "日本"},
			Name:            entity.SakeName{Name: "獺祭", Phonetic: "だっさい"},
			Abv:             16.0,
			PurchaseVolume:  720,
			RemainingVolume: 300,
			Price:           3000,
			DrinkStyles:     []entity.DrinkStyle{},
			CreatedAt:       now,
			UpdatedAt:       now,
		},
	}, nil)

	body := generated.UpdateSakeRequest{
		Category:        generated.SakeCategoryJAPANESESAKE,
		Kind:            generated.SakeKind{Id: 1, Name: "純米大吟醸"},
		Brewery:         generated.Brewery{Id: 1, Name: "旭酒造", OriginCountry: "日本"},
		Name:            generated.SakeName{Name: "獺祭", Phonetic: "だっさい"},
		Abv:             16.0,
		PurchaseVolume:  720,
		RemainingVolume: 300,
		Price:           3000,
		DrinkStyles:     []generated.DrinkStyle{},
	}
	bodyBytes, _ := json.Marshal(body)

	router := newRouter(server)
	req := httptest.NewRequest(http.MethodPut, "/stocks/1", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response generated.SakeDetail
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, int32(1), response.Id)
	assert.Equal(t, float32(300), response.RemainingVolume)
	updateUC.AssertExpectations(t)
	detailUC.AssertExpectations(t)
}

// テスト: 存在しない酒IDの更新
func TestUpdateStock_NotFound(t *testing.T) {
	server, _, _, _, updateUC, _, _, _, _, _, _, _ := newTestServer()

	updateUC.On("Execute", mock.Anything, mock.Anything).Return(
		nil,
		apperror.NotFoundError("酒が見つかりません"),
	)

	body := generated.UpdateSakeRequest{
		Category:        generated.SakeCategoryJAPANESESAKE,
		Kind:            generated.SakeKind{Name: "純米"},
		Brewery:         generated.Brewery{Name: "旭酒造", OriginCountry: "日本"},
		Name:            generated.SakeName{Name: "獺祭", Phonetic: "だっさい"},
		Abv:             16.0,
		PurchaseVolume:  720,
		RemainingVolume: 500,
		Price:           3000,
		DrinkStyles:     []generated.DrinkStyle{},
	}
	bodyBytes, _ := json.Marshal(body)

	router := newRouter(server)
	req := httptest.NewRequest(http.MethodPut, "/stocks/999", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	updateUC.AssertExpectations(t)
}

// --- DeleteStock Tests ---

// テスト: 在庫削除の正常処理
func TestDeleteStock_Success(t *testing.T) {
	server, _, _, _, _, deleteUC, _, _, _, _, _, _ := newTestServer()

	deleteUC.On("Execute", mock.Anything, int32(1)).Return(nil)

	router := newRouter(server)
	req := httptest.NewRequest(http.MethodDelete, "/stocks/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
	deleteUC.AssertExpectations(t)
}

// テスト: 存在しない酒IDの削除
func TestDeleteStock_NotFound(t *testing.T) {
	server, _, _, _, _, deleteUC, _, _, _, _, _, _ := newTestServer()

	deleteUC.On("Execute", mock.Anything, int32(999)).Return(
		apperror.NotFoundError("酒が見つかりません"),
	)

	router := newRouter(server)
	req := httptest.NewRequest(http.MethodDelete, "/stocks/999", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	deleteUC.AssertExpectations(t)
}

// --- PatchStock Tests ---

// テスト: 在庫部分更新の正常処理(objectKey更新後にSakeDetailを返す)
func TestPatchStock_Success(t *testing.T) {
	server, _, detailUC, _, _, _, patchUC, _, _, _, _, urlResolver := newTestServer()
	now := time.Now().Truncate(time.Second)

	objectKey := "sakes/1/test-uuid.jpg"

	patchUC.On("Execute", mock.Anything, usecase.PatchStockInput{
		ID:        1,
		ObjectKey: objectKey,
	}).Return(nil)

	urlResolver.On("ResolveURL", mock.Anything, objectKey).Return("https://signed-url.example.com/sakes/1/test-uuid.jpg", nil)
	detailUC.On("Execute", mock.Anything, int32(1)).Return(&usecase.GetSakeDetailOutput{
		Detail: &entity.SakeDetail{
			ID:              1,
			Category:        entity.SakeCategoryJapaneseSake,
			Kind:            entity.SakeKind{ID: 1, Name: "純米大吟醸"},
			Brewery:         entity.Brewery{ID: 1, Name: "旭酒造", OriginCountry: "日本"},
			Name:            entity.SakeName{Name: "獺祭", Phonetic: "だっさい"},
			Abv:             16.0,
			PurchaseVolume:  720,
			RemainingVolume: 500,
			Price:           3000,
			DrinkStyles:     []entity.DrinkStyle{},
			ObjectKey:       &objectKey,
			CreatedAt:       now,
			UpdatedAt:       now,
		},
	}, nil)

	body := generated.PatchStockRequest{ObjectKey: "sakes/1/test-uuid.jpg"}
	bodyBytes, _ := json.Marshal(body)

	router := newRouter(server)
	req := httptest.NewRequest(http.MethodPatch, "/stocks/1", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response generated.SakeDetail
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, int32(1), response.Id)
	patchUC.AssertExpectations(t)
	detailUC.AssertExpectations(t)
}

// テスト: objectKey空文字のバリデーションエラー
func TestPatchStock_ValidationError_MissingObjectKey(t *testing.T) {
	server, _, _, _, _, _, patchUC, _, _, _, _, _ := newTestServer()

	body := generated.PatchStockRequest{ObjectKey: ""}
	bodyBytes, _ := json.Marshal(body)

	router := newRouter(server)
	req := httptest.NewRequest(http.MethodPatch, "/stocks/1", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	patchUC.AssertNotCalled(t, "Execute")
}

// テスト: objectKeyにパストラバーサルが含まれる場合のバリデーションエラー
func TestPatchStock_ValidationError_PathTraversal(t *testing.T) {
	server, _, _, _, _, _, patchUC, _, _, _, _, _ := newTestServer()

	body := generated.PatchStockRequest{ObjectKey: "sakes/1/../../../etc/passwd"}
	bodyBytes, _ := json.Marshal(body)

	router := newRouter(server)
	req := httptest.NewRequest(http.MethodPatch, "/stocks/1", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	patchUC.AssertNotCalled(t, "Execute")
}

// テスト: objectKeyのプレフィックスが不正な場合のバリデーションエラー
func TestPatchStock_ValidationError_InvalidPrefix(t *testing.T) {
	server, _, _, _, _, _, patchUC, _, _, _, _, _ := newTestServer()

	body := generated.PatchStockRequest{ObjectKey: "invalid/path/image.jpg"}
	bodyBytes, _ := json.Marshal(body)

	router := newRouter(server)
	req := httptest.NewRequest(http.MethodPatch, "/stocks/1", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	patchUC.AssertNotCalled(t, "Execute")
}

// テスト: 存在しない酒IDの部分更新
func TestPatchStock_NotFound(t *testing.T) {
	server, _, _, _, _, _, patchUC, _, _, _, _, _ := newTestServer()

	patchUC.On("Execute", mock.Anything, mock.Anything).Return(
		apperror.NotFoundError("酒が見つかりません"),
	)

	body := generated.PatchStockRequest{ObjectKey: "sakes/999/test-uuid.jpg"}
	bodyBytes, _ := json.Marshal(body)

	router := newRouter(server)
	req := httptest.NewRequest(http.MethodPatch, "/stocks/999", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	patchUC.AssertExpectations(t)
}

// テスト: 不正JSONの場合
func TestPatchStock_InvalidJSON(t *testing.T) {
	server, _, _, _, _, _, patchUC, _, _, _, _, _ := newTestServer()

	router := newRouter(server)
	req := httptest.NewRequest(http.MethodPatch, "/stocks/1", bytes.NewReader([]byte("invalid")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	patchUC.AssertNotCalled(t, "Execute")
}

// --- CreateStockUploadUrl Tests ---

// テスト: 画像アップロードURL発行の正常処理
func TestCreateStockUploadUrl_Success(t *testing.T) {
	server, _, _, _, _, _, _, uploadUrlUC, _, _, _, _ := newTestServer()

	uploadUrlUC.On("Execute", mock.Anything, usecase.CreateUploadUrlInput{
		SakeID:      1,
		ContentType: "image/jpeg",
		Filename:    "photo.jpg",
	}).Return(&usecase.CreateUploadUrlOutput{
		UploadURL: "https://s3.example.com/presigned-url",
		ObjectKey: "sakes/1/uuid-value.jpg",
	}, nil)

	body := generated.PresignedUrlRequest{
		ContentType: "image/jpeg",
		Filename:    "photo.jpg",
	}
	bodyBytes, _ := json.Marshal(body)

	router := newRouter(server)
	req := httptest.NewRequest(http.MethodPost, "/stocks/1/upload-url", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response generated.PresignedUrlResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "https://s3.example.com/presigned-url", response.UploadUrl)
	assert.Equal(t, "sakes/1/uuid-value.jpg", response.ObjectKey)
	uploadUrlUC.AssertExpectations(t)
}

// テスト: Content-Type未指定のバリデーションエラー
func TestCreateStockUploadUrl_ValidationError_MissingContentType(t *testing.T) {
	server, _, _, _, _, _, _, uploadUrlUC, _, _, _, _ := newTestServer()

	body := generated.PresignedUrlRequest{
		ContentType: "",
		Filename:    "photo.jpg",
	}
	bodyBytes, _ := json.Marshal(body)

	router := newRouter(server)
	req := httptest.NewRequest(http.MethodPost, "/stocks/1/upload-url", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	uploadUrlUC.AssertNotCalled(t, "Execute")
}

// テスト: 許可されていないContent-Typeのバリデーションエラー
func TestCreateStockUploadUrl_ValidationError_InvalidContentType(t *testing.T) {
	server, _, _, _, _, _, _, uploadUrlUC, _, _, _, _ := newTestServer()

	body := generated.PresignedUrlRequest{
		ContentType: "application/pdf",
		Filename:    "doc.pdf",
	}
	bodyBytes, _ := json.Marshal(body)

	router := newRouter(server)
	req := httptest.NewRequest(http.MethodPost, "/stocks/1/upload-url", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	uploadUrlUC.AssertNotCalled(t, "Execute")
}

// テスト: ファイル名未指定のバリデーションエラー
func TestCreateStockUploadUrl_ValidationError_MissingFilename(t *testing.T) {
	server, _, _, _, _, _, _, uploadUrlUC, _, _, _, _ := newTestServer()

	body := generated.PresignedUrlRequest{
		ContentType: "image/png",
		Filename:    "",
	}
	bodyBytes, _ := json.Marshal(body)

	router := newRouter(server)
	req := httptest.NewRequest(http.MethodPost, "/stocks/1/upload-url", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	uploadUrlUC.AssertNotCalled(t, "Execute")
}

// テスト: ファイル名にパストラバーサルが含まれる場合のバリデーションエラー
func TestCreateStockUploadUrl_ValidationError_FilenamePathTraversal(t *testing.T) {
	server, _, _, _, _, _, _, uploadUrlUC, _, _, _, _ := newTestServer()

	body := generated.PresignedUrlRequest{
		ContentType: "image/jpeg",
		Filename:    "../../../etc/passwd",
	}
	bodyBytes, _ := json.Marshal(body)

	router := newRouter(server)
	req := httptest.NewRequest(http.MethodPost, "/stocks/1/upload-url", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	uploadUrlUC.AssertNotCalled(t, "Execute")
}

// テスト: ユースケースエラー発生時のレスポンス
func TestCreateStockUploadUrl_UsecaseError(t *testing.T) {
	server, _, _, _, _, _, _, uploadUrlUC, _, _, _, _ := newTestServer()

	uploadUrlUC.On("Execute", mock.Anything, mock.Anything).Return(
		nil,
		apperror.InternalServerError("署名付きアップロードURLの生成に失敗しました"),
	)

	body := generated.PresignedUrlRequest{
		ContentType: "image/jpeg",
		Filename:    "photo.jpg",
	}
	bodyBytes, _ := json.Marshal(body)

	router := newRouter(server)
	req := httptest.NewRequest(http.MethodPost, "/stocks/1/upload-url", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	uploadUrlUC.AssertExpectations(t)
}

// --- Error Handler Tests ---

// テスト: バリデーションエラーのHTTPレスポンス変換
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

	var response generated.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.NotNil(t, response.Errors)
	assert.Equal(t, apperror.ErrCodeNotFound, (*response.Errors)[0].Code)
}

// テスト: 予期しないエラーのHTTPレスポンス変換
func TestHandleError_UnexpectedError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	handleError(c, fmt.Errorf("unexpected error"))

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var response generated.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.NotNil(t, response.Errors)
	assert.Equal(t, apperror.ErrCodeInternalError, (*response.Errors)[0].Code)
}

// --- Boundary Value Tests ---

// テスト: offset最小値(0)での正常取得
func TestListSakes_BoundaryValue_OffsetMinimum(t *testing.T) {
	server, listUC, _, _, _, _, _, _, _, _, _, _ := newTestServer()

	listUC.On("Execute", mock.Anything, mock.MatchedBy(func(input usecase.ListSakesInput) bool {
		return input.Offset == 0
	})).Return(&usecase.ListSakesOutput{
		Sakes:      []entity.SakeListItem{},
		Pagination: entity.Pagination{Total: 0, Offset: 0, Limit: 20},
	}, nil)

	router := newRouter(server)
	req := httptest.NewRequest(http.MethodGet, "/sakes?offset=0", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	listUC.AssertExpectations(t)
}

// テスト: limit最小値(1)での正常取得
func TestListSakes_BoundaryValue_LimitMinimum(t *testing.T) {
	server, listUC, _, _, _, _, _, _, _, _, _, _ := newTestServer()

	listUC.On("Execute", mock.Anything, mock.MatchedBy(func(input usecase.ListSakesInput) bool {
		return input.Limit == 1
	})).Return(&usecase.ListSakesOutput{
		Sakes:      []entity.SakeListItem{},
		Pagination: entity.Pagination{Total: 0, Offset: 0, Limit: 1},
	}, nil)

	router := newRouter(server)
	req := httptest.NewRequest(http.MethodGet, "/sakes?limit=1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	listUC.AssertExpectations(t)
}

// テスト: limit最大値(100)での正常取得
func TestListSakes_BoundaryValue_LimitMaximum(t *testing.T) {
	server, listUC, _, _, _, _, _, _, _, _, _, _ := newTestServer()

	listUC.On("Execute", mock.Anything, mock.MatchedBy(func(input usecase.ListSakesInput) bool {
		return input.Limit == 100
	})).Return(&usecase.ListSakesOutput{
		Sakes:      []entity.SakeListItem{},
		Pagination: entity.Pagination{Total: 0, Offset: 0, Limit: 100},
	}, nil)

	router := newRouter(server)
	req := httptest.NewRequest(http.MethodGet, "/sakes?limit=100", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	listUC.AssertExpectations(t)
}

// --- ListKinds Tests ---

// テスト: 酒の種類一覧の正常取得
func TestListKinds_Success(t *testing.T) {
	server, _, _, _, _, _, _, _, listKindsUC, _, _, _ := newTestServer()

	listKindsUC.On("Execute", mock.Anything).Return(&usecase.ListKindsOutput{
		Kinds: []entity.SakeKind{
			{ID: 1, Name: "純米大吟醸"},
			{ID: 2, Name: "大吟醸"},
		},
	}, nil)

	router := newRouter(server)
	req := httptest.NewRequest(http.MethodGet, "/kinds", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response generated.ListKindsResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.NotNil(t, response.Data)
	assert.Len(t, *response.Data, 2)
	assert.Equal(t, int32(1), (*response.Data)[0].Id)
	assert.Equal(t, "純米大吟醸", (*response.Data)[0].Name)
	listKindsUC.AssertExpectations(t)
}

// テスト: 酒の種類一覧取得時のユースケースエラー
func TestListKinds_UsecaseError(t *testing.T) {
	server, _, _, _, _, _, _, _, listKindsUC, _, _, _ := newTestServer()

	listKindsUC.On("Execute", mock.Anything).Return(
		nil,
		apperror.DatabaseError("データベースエラー", fmt.Errorf("connection error")),
	)

	router := newRouter(server)
	req := httptest.NewRequest(http.MethodGet, "/kinds", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	listKindsUC.AssertExpectations(t)
}

// --- ListBreweries Tests ---

// テスト: 酒造一覧の正常取得(デフォルトパラメータ)
func TestListBreweries_Success(t *testing.T) {
	server, _, _, _, _, _, _, _, _, listBreweriesUC, _, _ := newTestServer()

	region := "山口県"
	lat := 34.1234
	lng := 131.5678
	listBreweriesUC.On("Execute", mock.Anything, mock.MatchedBy(func(input usecase.ListBreweriesInput) bool {
		return input.Keyword == nil && input.Limit == 50
	})).Return(&usecase.ListBreweriesOutput{
		Breweries: []entity.Brewery{
			{ID: 1, Name: "旭酒造", OriginCountry: "日本", OriginRegion: &region, Latitude: &lat, Longitude: &lng},
		},
	}, nil)

	router := newRouter(server)
	req := httptest.NewRequest(http.MethodGet, "/breweries", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response generated.ListBreweriesResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.NotNil(t, response.Data)
	assert.Len(t, *response.Data, 1)
	assert.Equal(t, int32(1), (*response.Data)[0].Id)
	assert.Equal(t, "旭酒造", (*response.Data)[0].Name)
	assert.Equal(t, "日本", (*response.Data)[0].OriginCountry)
	assert.Equal(t, &region, (*response.Data)[0].OriginRegion)
	listBreweriesUC.AssertExpectations(t)
}

// テスト: 酒造一覧のキーワード検索付き正常取得
func TestListBreweries_Success_WithKeyword(t *testing.T) {
	server, _, _, _, _, _, _, _, _, listBreweriesUC, _, _ := newTestServer()

	listBreweriesUC.On("Execute", mock.Anything, mock.MatchedBy(func(input usecase.ListBreweriesInput) bool {
		return input.Keyword != nil && *input.Keyword == "旭" && input.Limit == 10
	})).Return(&usecase.ListBreweriesOutput{
		Breweries: []entity.Brewery{
			{ID: 1, Name: "旭酒造", OriginCountry: "日本"},
		},
	}, nil)

	router := newRouter(server)
	req := httptest.NewRequest(http.MethodGet, "/breweries?keyword=%E6%97%AD&limit=10", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response generated.ListBreweriesResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.NotNil(t, response.Data)
	assert.Len(t, *response.Data, 1)
	listBreweriesUC.AssertExpectations(t)
}

// テスト: 酒造一覧のlimitが範囲外の場合のバリデーションエラー
func TestListBreweries_ValidationError_LimitOutOfRange(t *testing.T) {
	tests := []struct {
		name  string
		query string
	}{
		{"LimitLessThan1", "/breweries?limit=0"},
		{"LimitGreaterThan100", "/breweries?limit=101"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server, _, _, _, _, _, _, _, _, listBreweriesUC, _, _ := newTestServer()
			router := newRouter(server)

			req := httptest.NewRequest(http.MethodGet, tt.query, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusBadRequest, w.Code)
			listBreweriesUC.AssertNotCalled(t, "Execute")
		})
	}
}

// テスト: 酒造一覧取得時のユースケースエラー
func TestListBreweries_UsecaseError(t *testing.T) {
	server, _, _, _, _, _, _, _, _, listBreweriesUC, _, _ := newTestServer()

	listBreweriesUC.On("Execute", mock.Anything, mock.Anything).Return(
		nil,
		apperror.DatabaseError("データベースエラー", fmt.Errorf("connection error")),
	)

	router := newRouter(server)
	req := httptest.NewRequest(http.MethodGet, "/breweries", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	listBreweriesUC.AssertExpectations(t)
}

// --- ListDrinkStyles Tests ---

// テスト: 飲み方一覧の正常取得
func TestListDrinkStyles_Success(t *testing.T) {
	server, _, _, _, _, _, _, _, _, _, listDrinkStylesUC, _ := newTestServer()

	desc := "5-10度に冷やして"
	listDrinkStylesUC.On("Execute", mock.Anything).Return(&usecase.ListDrinkStylesOutput{
		DrinkStyles: []entity.DrinkStyle{
			{ID: 1, Name: "冷酒", Description: &desc},
			{ID: 2, Name: "常温", Description: nil},
		},
	}, nil)

	router := newRouter(server)
	req := httptest.NewRequest(http.MethodGet, "/drink-styles", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response generated.ListDrinkStylesResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.NotNil(t, response.Data)
	assert.Len(t, *response.Data, 2)
	assert.Equal(t, int32(1), (*response.Data)[0].Id)
	assert.Equal(t, "冷酒", (*response.Data)[0].Name)
	assert.Equal(t, &desc, (*response.Data)[0].Description)
	assert.Nil(t, (*response.Data)[1].Description)
	listDrinkStylesUC.AssertExpectations(t)
}

// テスト: 飲み方一覧取得時のユースケースエラー
func TestListDrinkStyles_UsecaseError(t *testing.T) {
	server, _, _, _, _, _, _, _, _, _, listDrinkStylesUC, _ := newTestServer()

	listDrinkStylesUC.On("Execute", mock.Anything).Return(
		nil,
		apperror.DatabaseError("データベースエラー", fmt.Errorf("connection error")),
	)

	router := newRouter(server)
	req := httptest.NewRequest(http.MethodGet, "/drink-styles", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	listDrinkStylesUC.AssertExpectations(t)
}
