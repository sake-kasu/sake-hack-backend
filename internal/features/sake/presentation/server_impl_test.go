package presentation

import (
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

type mockUpdateStockUsecase struct{}

func (m *mockUpdateStockUsecase) Execute(ctx context.Context, input usecase.UpdateStockInput) (*usecase.UpdateStockOutput, error) {
	return nil, nil
}

type mockDeleteStockUsecase struct{}

func (m *mockDeleteStockUsecase) Execute(ctx context.Context, id int32) error {
	return nil
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

func TestCreateStock_NotImplemented(t *testing.T) {
	listUC := new(mockListSakesUsecase)
	server := newTestServer(listUC)

	req := httptest.NewRequest(http.MethodPost, "/stocks", nil)
	w := httptest.NewRecorder()
	newRouter(server).ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotImplemented, w.Code)
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
