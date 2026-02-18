---
name: golang-testing
description: testify/suite + testify/mock によるGoテストパターン。TDD手法に基づくクリーンアーキテクチャ各層のテスト戦略。
---

# Goテストパターン

testify/suite + testify/mock を使用したGoテストパターン集。TDD手法に基づくクリーンアーキテクチャ各層のテスト戦略。

## 活用タイミング

- 新規Go関数・メソッドの作成
- 既存コードへのテストカバレッジ追加
- パフォーマンスクリティカルなコードのベンチマーク
- TDDワークフローの実施

## TDDワークフロー

### RED-GREEN-REFACTORサイクル

```
RED     → テストスイートで失敗するテストを先に書く
GREEN   → テストをパスする最小限のコードを書く
REFACTOR → テストが緑のまま改善する
REPEAT  → 次の要件へ
```

## testify/suite テストスイートパターン

### 基本構造

```go
package usecase_test

import (
    "context"
    "testing"

    "github.com/stretchr/testify/mock"
    "github.com/stretchr/testify/suite"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

// テストスイート定義
type CreateFieldUseCaseTestSuite struct {
    suite.Suite
    useCase  *usecase.CreateFieldUseCase
    mockRepo *MockFieldRepository
    ctx      context.Context
}

// SetupTest は各テストの前に実行される
func (s *CreateFieldUseCaseTestSuite) SetupTest() {
    s.mockRepo = new(MockFieldRepository)
    s.useCase = usecase.NewCreateFieldUseCase(s.mockRepo, slog.Default())
    s.ctx = context.Background()
}

// TearDownTest は各テストの後に実行される
func (s *CreateFieldUseCaseTestSuite) TearDownTest() {
    s.mockRepo.AssertExpectations(s.T())
}

// テストスイートのエントリーポイント
func TestCreateFieldUseCaseTestSuite(t *testing.T) {
    suite.Run(t, new(CreateFieldUseCaseTestSuite))
}
```

### テストメソッド命名規則

```go
// 正常系: 有効なデータで圃場を作成できることを確認
func (s *CreateFieldUseCaseTestSuite) TestExecute_Success_ValidInput() {
    // テスト内容
}

// 異常系: 必須フィールド欠如でバリデーションエラー
func (s *CreateFieldUseCaseTestSuite) TestExecute_ValidationError_MissingName() {
    // テスト内容
}

// 境界値: 名前の最大文字数での動作確認
func (s *CreateFieldUseCaseTestSuite) TestExecute_BoundaryValue_MaxNameLength() {
    // テスト内容
}

// ビジネスロジックエラー: 重複する圃場名
func (s *CreateFieldUseCaseTestSuite) TestExecute_BusinessLogicError_DuplicateName() {
    // テスト内容
}
```

## testify/mock モックパターン

### モック定義

```go
// mock_field_repository_test.go
package usecase_test

import (
    "context"

    "github.com/stretchr/testify/mock"
)

// MockFieldRepository はFieldRepositoryのモック実装
type MockFieldRepository struct {
    mock.Mock
}

func (m *MockFieldRepository) Create(ctx context.Context, field *entity.Field) error {
    args := m.Called(ctx, field)
    return args.Error(0)
}

func (m *MockFieldRepository) FindByID(ctx context.Context, id string) (*entity.Field, error) {
    args := m.Called(ctx, id)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*entity.Field), args.Error(1)
}

func (m *MockFieldRepository) FindByName(ctx context.Context, name string) (*entity.Field, error) {
    args := m.Called(ctx, name)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*entity.Field), args.Error(1)
}
```

### モックの使用

```go
// 正常系: 有効なデータで圃場を作成できることを確認
func (s *CreateFieldUseCaseTestSuite) TestExecute_Success_ValidInput() {
    input := &usecase.CreateFieldInput{
        Name:     "テスト圃場",
        Geometry: validGeometry,
    }

    // モックの期待値設定
    s.mockRepo.On("FindByName", s.ctx, "テスト圃場").Return(nil, apperror.NotFound("見つかりません"))
    s.mockRepo.On("Create", s.ctx, mock.AnythingOfType("*entity.Field")).Return(nil)

    // 実行
    result, err := s.useCase.Execute(s.ctx, input)

    // 検証
    require.NoError(s.T(), err, "圃場作成でエラーが発生しないこと")
    assert.Equal(s.T(), "テスト圃場", result.Name, "圃場名が一致すること")
}

// 異常系: 重複する圃場名の場合エラーを返すことを確認
func (s *CreateFieldUseCaseTestSuite) TestExecute_BusinessLogicError_DuplicateName() {
    input := &usecase.CreateFieldInput{
        Name: "既存の圃場",
    }

    existingField := &entity.Field{Name: "既存の圃場"}
    s.mockRepo.On("FindByName", s.ctx, "既存の圃場").Return(existingField, nil)

    result, err := s.useCase.Execute(s.ctx, input)

    require.Error(s.T(), err, "重複名でエラーが発生すること")
    assert.Nil(s.T(), result, "結果がnilであること")

    var appErr *apperror.AppError
    require.ErrorAs(s.T(), err, &appErr, "AppError型であること")
    assert.Equal(s.T(), apperror.CodeConflict, appErr.Code, "Conflictエラーコードであること")
}
```

## Presentation層テスト(Gin + httptest)

```go
package presentation_test

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/gin-gonic/gin"
    "github.com/stretchr/testify/suite"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

type FieldHandlerTestSuite struct {
    suite.Suite
    router      *gin.Engine
    mockUseCase *MockCreateFieldUseCase
}

func (s *FieldHandlerTestSuite) SetupTest() {
    gin.SetMode(gin.TestMode)
    s.router = gin.New()
    s.mockUseCase = new(MockCreateFieldUseCase)

    handler := presentation.NewFieldHandler(s.mockUseCase, slog.Default())
    s.router.POST("/fields", handler.CreateField)
}

func TestFieldHandlerTestSuite(t *testing.T) {
    suite.Run(t, new(FieldHandlerTestSuite))
}

// 正常系: 有効なリクエストで201を返すことを確認
func (s *FieldHandlerTestSuite) TestCreateField_Success_ValidRequest() {
    body := `{"name": "テスト圃場", "geometry": {"type": "Polygon", "coordinates": []}}`
    req := httptest.NewRequest(http.MethodPost, "/fields", bytes.NewBufferString(body))
    req.Header.Set("Content-Type", "application/json")
    w := httptest.NewRecorder()

    s.mockUseCase.On("Execute", mock.Anything, mock.Anything).Return(&entity.Field{
        Name: "テスト圃場",
    }, nil)

    s.router.ServeHTTP(w, req)

    assert.Equal(s.T(), http.StatusCreated, w.Code, "ステータスコードが201であること")
}

// 異常系: 不正なJSONで400を返すことを確認
func (s *FieldHandlerTestSuite) TestCreateField_ValidationError_InvalidJSON() {
    body := `{invalid json}`
    req := httptest.NewRequest(http.MethodPost, "/fields", bytes.NewBufferString(body))
    req.Header.Set("Content-Type", "application/json")
    w := httptest.NewRecorder()

    s.router.ServeHTTP(w, req)

    assert.Equal(s.T(), http.StatusBadRequest, w.Code, "ステータスコードが400であること")
}

// 異常系: 必須フィールド欠如で400を返すことを確認
func (s *FieldHandlerTestSuite) TestCreateField_ValidationError_MissingName() {
    body := `{"geometry": {"type": "Polygon", "coordinates": []}}`
    req := httptest.NewRequest(http.MethodPost, "/fields", bytes.NewBufferString(body))
    req.Header.Set("Content-Type", "application/json")
    w := httptest.NewRecorder()

    s.router.ServeHTTP(w, req)

    assert.Equal(s.T(), http.StatusBadRequest, w.Code, "ステータスコードが400であること")
}
```

## 統合テストパターン

```go
// field_repository_integration_test.go
package repository_test

import (
    "context"
    "testing"

    "github.com/stretchr/testify/suite"
    "github.com/stretchr/testify/require"
    "github.com/stretchr/testify/assert"
)

type FieldRepositoryIntegrationTestSuite struct {
    suite.Suite
    repo *repository.FieldRepository
    db   *sql.DB
    ctx  context.Context
}

func (s *FieldRepositoryIntegrationTestSuite) SetupSuite() {
    // TestContainersでPostgreSQLを起動
    // DBコネクションを確立
    // マイグレーションを実行
}

func (s *FieldRepositoryIntegrationTestSuite) TearDownSuite() {
    if s.db != nil {
        err := s.db.Close()
        if err != nil {
            log.Fatalf("DB接続のクローズに失敗: %v", err)
        }
    }
}

func (s *FieldRepositoryIntegrationTestSuite) SetupTest() {
    s.ctx = context.Background()
    // 各テスト前にデータをクリーンアップ
}

func TestFieldRepositoryIntegrationTestSuite(t *testing.T) {
    if testing.Short() {
        t.Skip("統合テストをスキップ")
    }
    suite.Run(t, new(FieldRepositoryIntegrationTestSuite))
}

// 正常系: 圃場をDBに保存して取得できることを確認
func (s *FieldRepositoryIntegrationTestSuite) TestCreate_Success_ValidField() {
    field := &entity.Field{
        Name:     "統合テスト圃場",
        Geometry: validGeometry,
    }

    err := s.repo.Create(s.ctx, field)
    require.NoError(s.T(), err, "圃場の作成に失敗しないこと")
    assert.NotEmpty(s.T(), field.ID, "IDが設定されること")

    // 取得して検証
    found, err := s.repo.FindByID(s.ctx, field.ID)
    require.NoError(s.T(), err, "圃場の取得に失敗しないこと")
    assert.Equal(s.T(), field.Name, found.Name, "圃場名が一致すること")
}
```

## カバレッジ

### カバレッジ目標

| 層 | 目標 |
|---|------|
| Presentation層 | 90% |
| Usecase層 | 95% |
| 重要ビジネスロジック | 100% |
| Domain層 | 100% |
| 生成コード | 除外 |

### テスト実行コマンド

```bash
# 全テスト実行
make test

# 単体テストのみ
make test-unit

# 統合テストのみ
make test-integration

# カバレッジレポート(生成コード除外)
make cover

# レースディテクション付き
go test -race ./...
```

## ベンチマーク

### 基本ベンチマーク

```go
func BenchmarkProcess(b *testing.B) {
    data := generateTestData(1000)
    b.ResetTimer() // セットアップ時間を除外

    for i := 0; i < b.N; i++ {
        Process(data)
    }
}

// 実行: go test -bench=BenchmarkProcess -benchmem
```

### サイズ別ベンチマーク

```go
func BenchmarkSort(b *testing.B) {
    sizes := []int{100, 1000, 10000, 100000}

    for _, size := range sizes {
        b.Run(fmt.Sprintf("size=%d", size), func(b *testing.B) {
            data := generateRandomSlice(size)
            b.ResetTimer()

            for i := 0; i < b.N; i++ {
                tmp := make([]int, len(data))
                copy(tmp, data)
                sort.Ints(tmp)
            }
        })
    }
}
```

## ファジング(Go 1.18+)

```go
func FuzzParseJSON(f *testing.F) {
    // シードコーパスを追加
    f.Add(`{"name": "テスト"}`)
    f.Add(`{"count": 123}`)
    f.Add(`[]`)
    f.Add(`""`)

    f.Fuzz(func(t *testing.T, input string) {
        var result map[string]interface{}
        err := json.Unmarshal([]byte(input), &result)

        if err != nil {
            return // 不正なJSONは期待通り
        }

        // パース成功した場合、再エンコードが動作すること
        _, err = json.Marshal(result)
        if err != nil {
            t.Errorf("Unmarshal成功後のMarshalに失敗: %v", err)
        }
    })
}

// 実行: go test -fuzz=FuzzParseJSON -fuzztime=30s
```

## ベストプラクティス

**すべきこと:**
- テストを先に書く(TDD)
- testify/suite でテストスイートを構成する
- `require.NoError(t, err)` でエラーを検証する
- 日本語のコメントで何をテストしているか記載する
- 独立したテストには `t.Parallel()` を使用する
- リソースのクリーンアップには `t.Cleanup()` を使用する
- テスト名は `Test{メソッド}_{種類}_{詳細}` の形式にする

**してはいけないこと:**
- `_` でのエラー無視(禁止)
- `t.Skip()` でテストをパスさせる(禁止)
- テスト内で `fmt.Printf` を使用する(`log.Fatalf` を使用)
- `time.Sleep()` をテストで使用する
- テストを無視して放置する

**心得**: テストはドキュメント。コードの使い方を示すもの。明確に書き、常に最新に保つこと。
