package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/sake-kasu/sake-hack-backend/internal/features/sake/application/query"
	"github.com/sake-kasu/sake-hack-backend/internal/features/sake/domain/entity"
)

// MockSakeQuery はSakeQueryのモック
type MockSakeQuery struct {
	mock.Mock
}

func (m *MockSakeQuery) List(ctx context.Context, filter query.ListSakesFilter) ([]entity.SakeListItem, entity.Pagination, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, entity.Pagination{}, args.Error(2)
	}
	return args.Get(0).([]entity.SakeListItem), args.Get(1).(entity.Pagination), args.Error(2)
}

func (m *MockSakeQuery) GetDetail(ctx context.Context, id int32) (*entity.SakeDetail, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.SakeDetail), args.Error(1)
}

func (m *MockSakeQuery) ListKinds(ctx context.Context) ([]entity.SakeKind, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entity.SakeKind), args.Error(1)
}

func (m *MockSakeQuery) ListBreweries(ctx context.Context, keyword *string, limit int32) ([]entity.Brewery, error) {
	args := m.Called(ctx, keyword, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entity.Brewery), args.Error(1)
}

func (m *MockSakeQuery) ListDrinkStyles(ctx context.Context) ([]entity.DrinkStyle, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entity.DrinkStyle), args.Error(1)
}

func TestListSakesUsecase_Execute_Success(t *testing.T) {
	mockQuery := new(MockSakeQuery)
	uc := NewListSakesUsecase(mockQuery)

	kindID := int32(1)
	breweryID := int32(2)
	expectedSakes := []entity.SakeListItem{
		{
			ID:       1,
			Category: entity.SakeCategoryJapaneseSake,
			Name:     "獺祭",
		},
	}
	expectedPagination := entity.Pagination{
		Total:  100,
		Offset: 0,
		Limit:  20,
	}

	mockQuery.On("List", mock.Anything, mock.MatchedBy(func(filter query.ListSakesFilter) bool {
		return filter.Offset == 0 && filter.Limit == 20 &&
			filter.KindID != nil && *filter.KindID == kindID &&
			filter.BreweryID != nil && *filter.BreweryID == breweryID
	})).Return(expectedSakes, expectedPagination, nil)

	input := ListSakesInput{
		KindID:    &kindID,
		BreweryID: &breweryID,
		Offset:    0,
		Limit:     20,
	}

	output, err := uc.Execute(context.Background(), input)

	assert.NoError(t, err)
	assert.NotNil(t, output)
	assert.Len(t, output.Sakes, 1)
	assert.Equal(t, int32(1), output.Sakes[0].ID)
	assert.Equal(t, entity.SakeCategoryJapaneseSake, output.Sakes[0].Category)
	assert.Equal(t, int64(100), output.Pagination.Total)
	mockQuery.AssertExpectations(t)
}

func TestListSakesUsecase_Execute_OffsetValidation(t *testing.T) {
	mockQuery := new(MockSakeQuery)
	uc := NewListSakesUsecase(mockQuery)

	mockQuery.On("List", mock.Anything, mock.MatchedBy(func(filter query.ListSakesFilter) bool {
		return filter.Offset == 0
	})).Return([]entity.SakeListItem{}, entity.Pagination{}, nil)

	input := ListSakesInput{
		Offset: -1,
		Limit:  20,
	}

	output, err := uc.Execute(context.Background(), input)

	assert.NoError(t, err)
	assert.NotNil(t, output)
	mockQuery.AssertExpectations(t)
}

func TestListSakesUsecase_Execute_LimitValidation_LessThan1(t *testing.T) {
	mockQuery := new(MockSakeQuery)
	uc := NewListSakesUsecase(mockQuery)

	mockQuery.On("List", mock.Anything, mock.MatchedBy(func(filter query.ListSakesFilter) bool {
		return filter.Limit == 20
	})).Return([]entity.SakeListItem{}, entity.Pagination{}, nil)

	input := ListSakesInput{
		Offset: 0,
		Limit:  0,
	}

	output, err := uc.Execute(context.Background(), input)

	assert.NoError(t, err)
	assert.NotNil(t, output)
	mockQuery.AssertExpectations(t)
}

func TestListSakesUsecase_Execute_LimitValidation_GreaterThan100(t *testing.T) {
	mockQuery := new(MockSakeQuery)
	uc := NewListSakesUsecase(mockQuery)

	mockQuery.On("List", mock.Anything, mock.MatchedBy(func(filter query.ListSakesFilter) bool {
		return filter.Limit == 20
	})).Return([]entity.SakeListItem{}, entity.Pagination{}, nil)

	input := ListSakesInput{
		Offset: 0,
		Limit:  101,
	}

	output, err := uc.Execute(context.Background(), input)

	assert.NoError(t, err)
	assert.NotNil(t, output)
	mockQuery.AssertExpectations(t)
}

func TestListSakesUsecase_Execute_QueryError(t *testing.T) {
	mockQuery := new(MockSakeQuery)
	uc := NewListSakesUsecase(mockQuery)

	expectedErr := errors.New("database error")
	mockQuery.On("List", mock.Anything, mock.Anything).Return(nil, entity.Pagination{}, expectedErr)

	input := ListSakesInput{
		Offset: 0,
		Limit:  20,
	}

	output, err := uc.Execute(context.Background(), input)

	assert.Error(t, err)
	assert.Nil(t, output)
	assert.Equal(t, expectedErr, err)
	mockQuery.AssertExpectations(t)
}

func TestListSakesUsecase_Execute_EmptyResult(t *testing.T) {
	mockQuery := new(MockSakeQuery)
	uc := NewListSakesUsecase(mockQuery)

	mockQuery.On("List", mock.Anything, mock.Anything).Return([]entity.SakeListItem{}, entity.Pagination{
		Total:  0,
		Offset: 0,
		Limit:  20,
	}, nil)

	input := ListSakesInput{
		Offset: 0,
		Limit:  20,
	}

	output, err := uc.Execute(context.Background(), input)

	assert.NoError(t, err)
	assert.NotNil(t, output)
	assert.Len(t, output.Sakes, 0)
	assert.Equal(t, int64(0), output.Pagination.Total)
	mockQuery.AssertExpectations(t)
}
