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

// --- Helper ---

func newTestServer() (*SakeServerImpl, *MockListSakesUsecase, *MockGetSakeDetailUsecase, *MockCreateStockUsecase, *MockUpdateStockUsecase, *MockDeleteStockUsecase) {
	listUC := new(MockListSakesUsecase)
	detailUC := new(MockGetSakeDetailUsecase)
	createUC := new(MockCreateStockUsecase)
	updateUC := new(MockUpdateStockUsecase)
	deleteUC := new(MockDeleteStockUsecase)
	server := NewSakeServerImpl(listUC, detailUC, createUC, updateUC, deleteUC)
	return server, listUC, detailUC, createUC, updateUC, deleteUC
}

func newRouter(server *SakeServerImpl) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	generated.RegisterHandlers(router, server)
	return router
}

// --- ListSakes Tests ---

func TestListSakes_Success(t *testing.T) {
	server, listUC, _, _, _, _ := newTestServer()

	expectedOutput := &usecase.ListSakesOutput{
		Sakes: []entity.SakeListItem{
			{
				ID:           1,
				Category:     entity.SakeCategoryJapaneseSake,
				Name:         "獺祭 純米大吟醸50",
				ImagePreview: "https://example.com/image.jpg",
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
	assert.Equal(t, "https://example.com/image.jpg", (*response.Data)[0].ImagePreview)
	assert.NotNil(t, response.Meta)
	assert.Equal(t, int64(100), response.Meta.Total)
	listUC.AssertExpectations(t)
}

func TestListSakes_ValidationError_OffsetLessThan0(t *testing.T) {
	server, listUC, _, _, _, _ := newTestServer()
	router := newRouter(server)

	req := httptest.NewRequest(http.MethodGet, "/sakes?offset=-1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	listUC.AssertNotCalled(t, "Execute")
}

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
			server, listUC, _, _, _, _ := newTestServer()
			router := newRouter(server)

			req := httptest.NewRequest(http.MethodGet, tt.query, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusBadRequest, w.Code)
			listUC.AssertNotCalled(t, "Execute")
		})
	}
}

func TestListSakes_ValidationError_TypeIdLessThan1(t *testing.T) {
	server, listUC, _, _, _, _ := newTestServer()
	router := newRouter(server)

	req := httptest.NewRequest(http.MethodGet, "/sakes?typeId=0", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	listUC.AssertNotCalled(t, "Execute")
}

func TestListSakes_ValidationError_BreweryIdLessThan1(t *testing.T) {
	server, listUC, _, _, _, _ := newTestServer()
	router := newRouter(server)

	req := httptest.NewRequest(http.MethodGet, "/sakes?breweryId=0", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	listUC.AssertNotCalled(t, "Execute")
}

func TestListSakes_UsecaseError(t *testing.T) {
	server, listUC, _, _, _, _ := newTestServer()

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

func TestListSakes_DefaultValues(t *testing.T) {
	server, listUC, _, _, _, _ := newTestServer()

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

func TestGetSakeDetail_Success(t *testing.T) {
	server, _, detailUC, _, _, _ := newTestServer()
	now := time.Now().Truncate(time.Second)
	description := "5-10度に冷やして"

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
	detailUC.AssertExpectations(t)
}

func TestGetSakeDetail_NotFound(t *testing.T) {
	server, _, detailUC, _, _, _ := newTestServer()

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

func TestListStocks_Success(t *testing.T) {
	server, listUC, _, _, _, _ := newTestServer()

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

func TestCreateStock_Success(t *testing.T) {
	server, _, _, createUC, _, _ := newTestServer()

	createUC.On("Execute", mock.Anything, mock.Anything).Return(&usecase.CreateStockOutput{
		Sake: &entity.SakeListItem{
			ID:           1,
			Category:     entity.SakeCategoryJapaneseSake,
			Name:         "獺祭",
			ImagePreview: "",
		},
	}, nil)

	body := generated.CreateSakeRequest{
		Category: generated.SakeCategoryJAPANESESAKE,
		Kind:     generated.SakeKind{Id: 0, Name: "純米大吟醸"},
		Brewery:  generated.Brewery{Id: 0, Name: "旭酒造", OriginCountry: "日本"},
		Name:     generated.SakeName{Name: "獺祭", Phonetic: "だっさい"},
		Abv:      16.0,
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

func TestCreateStock_ValidationError_MissingName(t *testing.T) {
	server, _, _, createUC, _, _ := newTestServer()

	body := generated.CreateSakeRequest{
		Category: generated.SakeCategoryJAPANESESAKE,
		Kind:     generated.SakeKind{Name: "純米"},
		Brewery:  generated.Brewery{Name: "旭酒造", OriginCountry: "日本"},
		Name:     generated.SakeName{Name: "", Phonetic: ""},
		Abv:      16.0,
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

func TestCreateStock_InvalidJSON(t *testing.T) {
	server, _, _, createUC, _, _ := newTestServer()

	router := newRouter(server)
	req := httptest.NewRequest(http.MethodPost, "/stocks", bytes.NewReader([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	createUC.AssertNotCalled(t, "Execute")
}

// --- GetStockDetail Tests ---

func TestGetStockDetail_Success(t *testing.T) {
	server, _, detailUC, _, _, _ := newTestServer()
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

func TestUpdateStock_Success(t *testing.T) {
	server, _, detailUC, _, updateUC, _ := newTestServer()
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

func TestUpdateStock_NotFound(t *testing.T) {
	server, _, _, _, updateUC, _ := newTestServer()

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

func TestDeleteStock_Success(t *testing.T) {
	server, _, _, _, _, deleteUC := newTestServer()

	deleteUC.On("Execute", mock.Anything, int32(1)).Return(nil)

	router := newRouter(server)
	req := httptest.NewRequest(http.MethodDelete, "/stocks/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
	deleteUC.AssertExpectations(t)
}

func TestDeleteStock_NotFound(t *testing.T) {
	server, _, _, _, _, deleteUC := newTestServer()

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

// --- Error Handler Tests ---

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

func TestListSakes_BoundaryValue_OffsetMinimum(t *testing.T) {
	server, listUC, _, _, _, _ := newTestServer()

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

func TestListSakes_BoundaryValue_LimitMinimum(t *testing.T) {
	server, listUC, _, _, _, _ := newTestServer()

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

func TestListSakes_BoundaryValue_LimitMaximum(t *testing.T) {
	server, listUC, _, _, _, _ := newTestServer()

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
