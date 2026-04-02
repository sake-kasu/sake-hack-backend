package presentation

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/sake-kasu/sake-hack-backend/api/generated"
	"github.com/sake-kasu/sake-hack-backend/internal/features/sake/application/usecase"
	"github.com/sake-kasu/sake-hack-backend/internal/features/sake/domain/entity"
)

type mockListSakesUsecase struct {
	mock.Mock
}

func (m *mockListSakesUsecase) Execute(ctx context.Context, input usecase.ListSakesInput) (*usecase.ListSakesOutput, error) {
	args := m.Called(ctx, input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*usecase.ListSakesOutput), args.Error(1)
}

type mockGetSakeDetailUsecase struct{}

func (m *mockGetSakeDetailUsecase) Execute(ctx context.Context, id uuid.UUID) (*usecase.GetSakeDetailOutput, error) {
	return nil, nil
}

type mockCreateStockUsecase struct{}

func (m *mockCreateStockUsecase) Execute(ctx context.Context, input usecase.CreateStockInput) (*usecase.CreateStockOutput, error) {
	return nil, nil
}

type mockCreateStockUsecaseWithResult struct {
	output *usecase.CreateStockOutput
}

func (m *mockCreateStockUsecaseWithResult) Execute(ctx context.Context, input usecase.CreateStockInput) (*usecase.CreateStockOutput, error) {
	return m.output, nil
}

type mockUpdateStockUsecase struct{}

func (m *mockUpdateStockUsecase) Execute(ctx context.Context, input usecase.UpdateStockInput) (*usecase.UpdateStockOutput, error) {
	return nil, nil
}

type mockUpdateStockUsecaseWithResult struct {
	output *usecase.UpdateStockOutput
}

func (m *mockUpdateStockUsecaseWithResult) Execute(ctx context.Context, input usecase.UpdateStockInput) (*usecase.UpdateStockOutput, error) {
	return m.output, nil
}

type mockDeleteStockUsecase struct{}

func (m *mockDeleteStockUsecase) Execute(ctx context.Context, id uuid.UUID) error {
	return nil
}

type mockDeleteStockUsecaseWithResult struct {
	err error
}

func (m *mockDeleteStockUsecaseWithResult) Execute(ctx context.Context, id uuid.UUID) error {
	return m.err
}

func newTestServer(listUC *mockListSakesUsecase) *SakeServerImpl {
	return NewSakeServerImpl(
		listUC,
		&mockGetSakeDetailUsecase{},
		&mockCreateStockUsecase{},
		&mockUpdateStockUsecase{},
		&mockDeleteStockUsecase{},
	)
}

func newRouter(server *SakeServerImpl) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	generated.RegisterHandlers(router, server)
	return router
}

func TestListSakes_Success(t *testing.T) {
	listUC := new(mockListSakesUsecase)
	server := newTestServer(listUC)

	expectedOutput := &usecase.ListSakesOutput{
		Sakes: []entity.SakeListItem{
			{
				ID:           uuid.MustParse("11111111-1111-1111-1111-111111111111"),
				Category:     entity.SakeCategoryJapaneseSake,
				Name:         "獺祭 純米大吟醸50",
				ImagePreview: "preview.jpg",
			},
		},
		Pagination: entity.Pagination{Total: 1, Offset: 0, Limit: 20},
	}

	listUC.On("Execute", mock.Anything, usecase.ListSakesInput{
		KindID: nil, BreweryID: nil, Offset: 0, Limit: 20,
	}).Return(expectedOutput, nil)

	req := httptest.NewRequest(http.MethodGet, "/sakes?offset=0&limit=20", nil)
	w := httptest.NewRecorder()
	newRouter(server).ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response generated.ListSakesResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Len(t, response.Data, 1)
	assert.Equal(t, "獺祭 純米大吟醸50", string(response.Data[0].Name))
	assert.Equal(t, generated.JAPANESESAKE, response.Data[0].Category)
}

func TestListStocks_Success(t *testing.T) {
	listUC := new(mockListSakesUsecase)
	server := newTestServer(listUC)

	expectedOutput := &usecase.ListSakesOutput{
		Sakes: []entity.SakeListItem{
			{
				ID:           uuid.MustParse("22222222-2222-2222-2222-222222222222"),
				Category:     entity.SakeCategoryBeer,
				Name:         "Yona Yona Ale",
				ImagePreview: "preview.jpg",
			},
		},
		Pagination: entity.Pagination{Total: 1, Offset: 0, Limit: 20},
	}

	listUC.On("Execute", mock.Anything, usecase.ListSakesInput{
		KindID: nil, BreweryID: nil, Offset: 0, Limit: 20,
	}).Return(expectedOutput, nil)

	req := httptest.NewRequest(http.MethodGet, "/stocks?offset=0&limit=20", nil)
	w := httptest.NewRecorder()
	newRouter(server).ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response generated.ListStockResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Len(t, response.Data, 1)
	assert.Equal(t, "Yona Yona Ale", string(response.Data[0].Name))
	assert.Equal(t, generated.BEER, response.Data[0].Category)
}

func TestCreateStock_InvalidRequest(t *testing.T) {
	listUC := new(mockListSakesUsecase)
	server := newTestServer(listUC)

	req := httptest.NewRequest(http.MethodPost, "/stocks", nil)
	w := httptest.NewRecorder()
	newRouter(server).ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreateStock_Success(t *testing.T) {
	id := uuid.MustParse("44444444-4444-4444-4444-444444444444")
	server := NewSakeServerImpl(
		new(mockListSakesUsecase),
		&mockGetSakeDetailUsecaseWithResult{
			output: &usecase.GetSakeDetailOutput{
				Detail: &entity.SakeDetail{
					ID:       id,
					Category: entity.SakeCategoryJapaneseSake,
					Name:     entity.SakeName{Name: "獺祭", Phonetic: "だっさい"},
					Brewery: entity.Brewery{
						OriginRegion: strPtr("Yamaguchi"),
					},
					Abv:             float32Ptr(16),
					PurchaseVolume:  int32Ptr(720),
					RemainingVolume: int32Ptr(100),
					Price:           int32Ptr(3000),
					Memo:            strPtr("memo"),
				},
			},
		},
		&mockCreateStockUsecaseWithResult{
			output: &usecase.CreateStockOutput{
				Sake: &entity.SakeListItem{
					ID:       id,
					Category: entity.SakeCategoryJapaneseSake,
					Name:     "獺祭",
				},
			},
		},
		&mockUpdateStockUsecase{},
		&mockDeleteStockUsecase{},
	)

	req := httptest.NewRequest(http.MethodPost, "/stocks", bytes.NewBufferString(`{
		"name":"獺祭",
		"category":"JAPANESE_SAKE",
		"phonetic":"だっさい",
		"alcoholPercentage":16,
		"volumeMax":720,
		"volumeRemain":100,
		"region":"Yamaguchi",
		"price":3000,
		"memo":"memo",
		"tagNames":["fruity"],
		"imageKeys":["img-1"]
	}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	newRouter(server).ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestGetStockDetail_Success(t *testing.T) {
	id := uuid.MustParse("33333333-3333-3333-3333-333333333333")
	server := NewSakeServerImpl(
		new(mockListSakesUsecase),
		&mockGetSakeDetailUsecaseWithResult{
			output: &usecase.GetSakeDetailOutput{
				Detail: &entity.SakeDetail{
					ID:       id,
					Category: entity.SakeCategoryBeer,
					Name: entity.SakeName{
						Name:     "Yona Yona Ale",
						Phonetic: "",
					},
					Brewery: entity.Brewery{
						OriginRegion: strPtr("Nagano"),
					},
					Abv:             float32Ptr(5.5),
					PurchaseVolume:  int32Ptr(350),
					RemainingVolume: int32Ptr(80),
					Price:           int32Ptr(320),
					Memo:            strPtr("test"),
				},
			},
		},
		&mockCreateStockUsecase{},
		&mockUpdateStockUsecase{},
		&mockDeleteStockUsecase{},
	)

	req := httptest.NewRequest(http.MethodGet, "/stocks/"+id.String(), nil)
	w := httptest.NewRecorder()
	newRouter(server).ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestUpdateStock_Success(t *testing.T) {
	id := uuid.MustParse("55555555-5555-5555-5555-555555555555")
	server := NewSakeServerImpl(
		new(mockListSakesUsecase),
		&mockGetSakeDetailUsecaseWithResult{
			output: &usecase.GetSakeDetailOutput{
				Detail: &entity.SakeDetail{
					ID:       id,
					Category: entity.SakeCategoryJapaneseSake,
					Name:     entity.SakeName{Name: "新しい獺祭", Phonetic: "あたらしいだっさい"},
					Brewery: entity.Brewery{
						OriginRegion: strPtr("Yamaguchi"),
					},
					Abv:             float32Ptr(15),
					PurchaseVolume:  int32Ptr(720),
					RemainingVolume: int32Ptr(90),
					Price:           int32Ptr(3200),
					Memo:            strPtr("updated"),
				},
			},
		},
		&mockCreateStockUsecase{},
		&mockUpdateStockUsecaseWithResult{
			output: &usecase.UpdateStockOutput{
				Sake: &entity.SakeListItem{
					ID:       id,
					Category: entity.SakeCategoryJapaneseSake,
					Name:     "新しい獺祭",
				},
			},
		},
		&mockDeleteStockUsecase{},
	)

	req := httptest.NewRequest(http.MethodPut, "/stocks/"+id.String(), bytes.NewBufferString(`{
		"name":"新しい獺祭",
		"category":"JAPANESE_SAKE",
		"phonetic":"あたらしいだっさい",
		"alcoholPercentage":15,
		"volumeMax":720,
		"volumeRemain":90,
		"region":"Yamaguchi",
		"price":3200,
		"memo":"updated",
		"tagNames":["fruity"],
		"imageKeys":["img-1"]
	}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	newRouter(server).ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestDeleteStock_Success(t *testing.T) {
	id := uuid.MustParse("66666666-6666-6666-6666-666666666666")
	server := NewSakeServerImpl(
		new(mockListSakesUsecase),
		&mockGetSakeDetailUsecase{},
		&mockCreateStockUsecase{},
		&mockUpdateStockUsecase{},
		&mockDeleteStockUsecaseWithResult{},
	)

	req := httptest.NewRequest(http.MethodDelete, "/stocks/"+id.String(), nil)
	w := httptest.NewRecorder()
	newRouter(server).ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
}

type mockGetSakeDetailUsecaseWithResult struct {
	output *usecase.GetSakeDetailOutput
}

func (m *mockGetSakeDetailUsecaseWithResult) Execute(ctx context.Context, id uuid.UUID) (*usecase.GetSakeDetailOutput, error) {
	return m.output, nil
}

func strPtr(v string) *string {
	return &v
}

func int32Ptr(v int32) *int32 {
	return &v
}

func float32Ptr(v float32) *float32 {
	return &v
}
