package presentation

import (
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

// MockListSakesUsecase はListSakesUsecaseのモック
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

func TestListSakes_Success(t *testing.T) {
	// テストデータ準備
	now := time.Now()
	originRegion := "山口県"
	description := "5-10度に冷やして"
	expectedOutput := &usecase.ListSakesOutput{
		Sakes: []entity.Sake{
			{
				ID: 1,
				Type: entity.SakeType{
					ID:   1,
					Name: "純米大吟醸",
				},
				Brewery: entity.Brewery{
					ID:            1,
					Name:          "旭酒造",
					OriginCountry: "日本",
					OriginRegion:  &originRegion,
				},
				Name:       "獺祭 純米大吟醸50",
				ABV:        16.0,
				TasteNotes: "フルーティーな香りと上品な甘み",
				DrinkStyles: []entity.DrinkStyle{
					{
						ID:          1,
						Name:        "冷酒",
						Description: &description,
					},
				},
				CreatedAt: now,
				UpdatedAt: now,
			},
		},
		Pagination: entity.Pagination{
			Total:  100,
			Offset: 0,
			Limit:  20,
		},
	}

	// モック設定
	mockUsecase := new(MockListSakesUsecase)
	mockUsecase.On("Execute", mock.Anything, mock.MatchedBy(func(input usecase.ListSakesInput) bool {
		return input.Offset == 0 && input.Limit == 20
	})).Return(expectedOutput, nil)

	// テスト対象
	server := NewSakeServerImpl(mockUsecase)

	// Ginテストモード
	gin.SetMode(gin.TestMode)
	router := gin.New()
	generated.RegisterHandlers(router, server)

	// リクエスト作成
	req := httptest.NewRequest(http.MethodGet, "/sakes?page=1&per_page=20", nil)
	w := httptest.NewRecorder()

	// リクエスト実行
	router.ServeHTTP(w, req)

	// レスポンス検証
	assert.Equal(t, http.StatusOK, w.Code)

	var response generated.ListSakesResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.NotNil(t, response.Data)
	assert.Len(t, *response.Data, 1)
	assert.Equal(t, int32(1), (*response.Data)[0].Id)
	assert.Equal(t, "獺祭 純米大吟醸50", (*response.Data)[0].Name)
	assert.Equal(t, generated.SakeCategory("純米大吟醸"), (*response.Data)[0].Category)
	assert.Equal(t, "", (*response.Data)[0].ImagePreview)
	assert.NotNil(t, response.Meta)
	assert.Equal(t, int64(100), response.Meta.Total)

	mockUsecase.AssertExpectations(t)
}

func TestListSakes_ValidationError_OffsetLessThan0(t *testing.T) {
	mockUsecase := new(MockListSakesUsecase)
	server := NewSakeServerImpl(mockUsecase)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	generated.RegisterHandlers(router, server)

	req := httptest.NewRequest(http.MethodGet, "/sakes?offset=-1", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response generated.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.NotNil(t, response.Errors)
	assert.Greater(t, len(*response.Errors), 0)
	assert.Equal(t, apperror.ErrCodeValidationError, (*response.Errors)[0].Code)

	mockUsecase.AssertNotCalled(t, "Execute")
}

func TestListSakes_ValidationError_LimitLessThan1(t *testing.T) {
	mockUsecase := new(MockListSakesUsecase)
	server := NewSakeServerImpl(mockUsecase)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	generated.RegisterHandlers(router, server)

	req := httptest.NewRequest(http.MethodGet, "/sakes?limit=0", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockUsecase.AssertNotCalled(t, "Execute")
}

func TestListSakes_ValidationError_LimitGreaterThan100(t *testing.T) {
	mockUsecase := new(MockListSakesUsecase)
	server := NewSakeServerImpl(mockUsecase)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	generated.RegisterHandlers(router, server)

	req := httptest.NewRequest(http.MethodGet, "/sakes?limit=101", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockUsecase.AssertNotCalled(t, "Execute")
}

func TestListSakes_ValidationError_TypeIdLessThan1(t *testing.T) {
	mockUsecase := new(MockListSakesUsecase)
	server := NewSakeServerImpl(mockUsecase)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	generated.RegisterHandlers(router, server)

	req := httptest.NewRequest(http.MethodGet, "/sakes?typeId=0", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockUsecase.AssertNotCalled(t, "Execute")
}

func TestListSakes_ValidationError_BreweryIdLessThan1(t *testing.T) {
	mockUsecase := new(MockListSakesUsecase)
	server := NewSakeServerImpl(mockUsecase)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	generated.RegisterHandlers(router, server)

	req := httptest.NewRequest(http.MethodGet, "/sakes?breweryId=0", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockUsecase.AssertNotCalled(t, "Execute")
}

func TestListSakes_UsecaseError_Database(t *testing.T) {
	mockUsecase := new(MockListSakesUsecase)
	mockUsecase.On("Execute", mock.Anything, mock.Anything).Return(
		nil,
		apperror.DatabaseError("データベースエラーが発生しました", fmt.Errorf("connection error")),
	)

	server := NewSakeServerImpl(mockUsecase)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	generated.RegisterHandlers(router, server)

	req := httptest.NewRequest(http.MethodGet, "/sakes", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var response generated.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.NotNil(t, response.Errors)
	assert.Greater(t, len(*response.Errors), 0)
	assert.Equal(t, apperror.ErrCodeDatabaseError, (*response.Errors)[0].Code)

	mockUsecase.AssertExpectations(t)
}

func TestListSakes_DefaultValues(t *testing.T) {
	mockUsecase := new(MockListSakesUsecase)
	mockUsecase.On("Execute", mock.Anything, mock.MatchedBy(func(input usecase.ListSakesInput) bool {
		return input.Offset == 0 && input.Limit == 20 && input.TypeID == nil && input.BreweryID == nil
	})).Return(&usecase.ListSakesOutput{
		Sakes:      []entity.Sake{},
		Pagination: entity.Pagination{Total: 0, Offset: 0, Limit: 20},
	}, nil)

	server := NewSakeServerImpl(mockUsecase)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	generated.RegisterHandlers(router, server)

	req := httptest.NewRequest(http.MethodGet, "/sakes", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockUsecase.AssertExpectations(t)
}

func TestListSakes_BoundaryValue_OffsetMinimum(t *testing.T) {
	mockUsecase := new(MockListSakesUsecase)
	mockUsecase.On("Execute", mock.Anything, mock.MatchedBy(func(input usecase.ListSakesInput) bool {
		return input.Offset == 0 && input.Limit == 20
	})).Return(&usecase.ListSakesOutput{
		Sakes:      []entity.Sake{},
		Pagination: entity.Pagination{Total: 0, Offset: 0, Limit: 20},
	}, nil)

	server := NewSakeServerImpl(mockUsecase)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	generated.RegisterHandlers(router, server)

	req := httptest.NewRequest(http.MethodGet, "/sakes?offset=0", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockUsecase.AssertExpectations(t)
}

func TestListSakes_BoundaryValue_LimitMinimum(t *testing.T) {
	mockUsecase := new(MockListSakesUsecase)
	mockUsecase.On("Execute", mock.Anything, mock.MatchedBy(func(input usecase.ListSakesInput) bool {
		return input.Offset == 0 && input.Limit == 1
	})).Return(&usecase.ListSakesOutput{
		Sakes:      []entity.Sake{},
		Pagination: entity.Pagination{Total: 0, Offset: 0, Limit: 1},
	}, nil)

	server := NewSakeServerImpl(mockUsecase)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	generated.RegisterHandlers(router, server)

	req := httptest.NewRequest(http.MethodGet, "/sakes?limit=1", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockUsecase.AssertExpectations(t)
}

func TestListSakes_BoundaryValue_LimitMaximum(t *testing.T) {
	mockUsecase := new(MockListSakesUsecase)
	mockUsecase.On("Execute", mock.Anything, mock.MatchedBy(func(input usecase.ListSakesInput) bool {
		return input.Offset == 0 && input.Limit == 100
	})).Return(&usecase.ListSakesOutput{
		Sakes:      []entity.Sake{},
		Pagination: entity.Pagination{Total: 0, Offset: 0, Limit: 100},
	}, nil)

	server := NewSakeServerImpl(mockUsecase)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	generated.RegisterHandlers(router, server)

	req := httptest.NewRequest(http.MethodGet, "/sakes?limit=100", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockUsecase.AssertExpectations(t)
}

func TestListSakes_BoundaryValue_TypeIdMinimum(t *testing.T) {
	mockUsecase := new(MockListSakesUsecase)
	mockUsecase.On("Execute", mock.Anything, mock.MatchedBy(func(input usecase.ListSakesInput) bool {
		return input.Offset == 0 && input.Limit == 20 && input.TypeID != nil && *input.TypeID == 1
	})).Return(&usecase.ListSakesOutput{
		Sakes:      []entity.Sake{},
		Pagination: entity.Pagination{Total: 0, Offset: 0, Limit: 20},
	}, nil)

	server := NewSakeServerImpl(mockUsecase)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	generated.RegisterHandlers(router, server)

	req := httptest.NewRequest(http.MethodGet, "/sakes?typeId=1", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockUsecase.AssertExpectations(t)
}

func TestListSakes_BoundaryValue_BreweryIdMinimum(t *testing.T) {
	mockUsecase := new(MockListSakesUsecase)
	mockUsecase.On("Execute", mock.Anything, mock.MatchedBy(func(input usecase.ListSakesInput) bool {
		return input.Offset == 0 && input.Limit == 20 && input.BreweryID != nil && *input.BreweryID == 1
	})).Return(&usecase.ListSakesOutput{
		Sakes:      []entity.Sake{},
		Pagination: entity.Pagination{Total: 0, Offset: 0, Limit: 20},
	}, nil)

	server := NewSakeServerImpl(mockUsecase)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	generated.RegisterHandlers(router, server)

	req := httptest.NewRequest(http.MethodGet, "/sakes?breweryId=1", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockUsecase.AssertExpectations(t)
}

func TestToSakeResponse(t *testing.T) {
	now := time.Now()
	description := "冷やして"
	memo := "メモ"
	sake := entity.Sake{
		ID: 1,
		Type: entity.SakeType{
			ID:   1,
			Name: "純米大吟醸",
		},
		Brewery: entity.Brewery{
			ID:            1,
			Name:          "旭酒造",
			OriginCountry: "日本",
		},
		Name:       "獺祭",
		ABV:        16.0,
		TasteNotes: "フルーティー",
		Memo:       &memo,
		DrinkStyles: []entity.DrinkStyle{
			{
				ID:          1,
				Name:        "冷酒",
				Description: &description,
			},
		},
		CreatedAt: now,
		UpdatedAt: now,
	}

	response := toSakeResponse(sake)

	assert.Equal(t, int32(1), response.Id)
	assert.Equal(t, "獺祭", response.Name)
	assert.Equal(t, generated.SakeCategory("純米大吟醸"), response.Category)
	assert.Equal(t, "", response.ImagePreview)
}

func TestValidateListSakesParams_Success(t *testing.T) {
	offset := int32(0)
	limit := int32(20)
	params := generated.ListSakesParams{
		Offset: &offset,
		Limit:  &limit,
	}

	err := validateListSakesParams(params)
	assert.NoError(t, err)
}

func TestValidateListSakesParams_MultipleErrors(t *testing.T) {
	offset := int32(-1)
	limit := int32(101)
	params := generated.ListSakesParams{
		Offset: &offset,
		Limit:  &limit,
	}

	err := validateListSakesParams(params)
	assert.Error(t, err)

	valErr, ok := err.(*apperror.ValidationError)
	assert.True(t, ok)
	assert.True(t, valErr.HasErrors())
	assert.Contains(t, valErr.Fields, "offset")
	assert.Contains(t, valErr.Fields, "limit")
}

func TestHandleError_ValidationError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	verr := apperror.NewValidationError("バリデーションエラー")
	verr = verr.AddField("field1", "エラーメッセージ")

	handleError(c, verr)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response generated.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.NotNil(t, response.Errors)
	assert.Greater(t, len(*response.Errors), 0)
	assert.Equal(t, apperror.ErrCodeValidationError, (*response.Errors)[0].Code)
	assert.Contains(t, (*response.Errors)[0].Message, "field1")
}

func TestHandleError_AppError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	appErr := apperror.NotFoundError("見つかりません")

	handleError(c, appErr)

	assert.Equal(t, http.StatusNotFound, w.Code)

	var response generated.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.NotNil(t, response.Errors)
	assert.Greater(t, len(*response.Errors), 0)
	assert.Equal(t, apperror.ErrCodeNotFound, (*response.Errors)[0].Code)
}

func TestHandleError_UnexpectedError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// 予期しないエラー(AppErrorでもValidationErrorでもない)
	unexpectedErr := fmt.Errorf("unexpected error")

	handleError(c, unexpectedErr)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var response generated.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.NotNil(t, response.Errors)
	assert.Greater(t, len(*response.Errors), 0)
	assert.Equal(t, apperror.ErrCodeInternalError, (*response.Errors)[0].Code)
	assert.Equal(t, "内部エラーが発生しました", (*response.Errors)[0].Message)
}
