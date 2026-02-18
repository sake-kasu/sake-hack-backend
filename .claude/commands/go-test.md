---
description: TDDワークフローを強制する。testify/suite + testify/mockでテストを先に書き、実装し、カバレッジを検証する。
---

# Go TDDコマンド

testify/suite + testify/mock を使用したTDD(テスト駆動開発)手法をGoコードに適用する。

## 実行内容

1. **型/インターフェースの定義**: 関数シグネチャを先にスキャフォールド
2. **テストスイートの作成**: testify/suiteでテストケースを作成(RED)
3. **テストの実行**: テストが正しい理由で失敗することを確認
4. **実装**: テストをパスする最小限のコードを作成(GREEN)
5. **リファクタリング**: テストが緑のまま改善
6. **カバレッジ確認**: 目標カバレッジを達成

## 使用タイミング

- 新しいGo関数・メソッドの実装
- 既存コードへのテストカバレッジ追加
- バグ修正(失敗するテストを先に書く)
- 重要なビジネスロジックの構築

## TDDサイクル

```
RED     → testify/suiteで失敗するテストを書く
GREEN   → テストをパスする最小限のコードを実装
REFACTOR → テストが緑のまま改善
REPEAT  → 次のテストケースへ
```

## テストスイートの構成例

```go
package usecase_test

import (
    "context"
    "testing"

    "github.com/stretchr/testify/suite"
    "github.com/stretchr/testify/mock"
    "github.com/stretchr/testify/require"
    "github.com/stretchr/testify/assert"
)

type CreateFieldUseCaseTestSuite struct {
    suite.Suite
    useCase  *usecase.CreateFieldUseCase
    mockRepo *MockFieldRepository
    ctx      context.Context
}

func (s *CreateFieldUseCaseTestSuite) SetupTest() {
    s.mockRepo = new(MockFieldRepository)
    s.useCase = usecase.NewCreateFieldUseCase(s.mockRepo, slog.Default())
    s.ctx = context.Background()
}

func (s *CreateFieldUseCaseTestSuite) TearDownTest() {
    s.mockRepo.AssertExpectations(s.T())
}

func TestCreateFieldUseCaseTestSuite(t *testing.T) {
    suite.Run(t, new(CreateFieldUseCaseTestSuite))
}

// 正常系: 有効なデータで圃場を作成できることを確認
func (s *CreateFieldUseCaseTestSuite) TestExecute_Success_ValidInput() {
    input := &usecase.CreateFieldInput{Name: "テスト圃場"}
    s.mockRepo.On("Create", s.ctx, mock.Anything).Return(nil)

    result, err := s.useCase.Execute(s.ctx, input)

    require.NoError(s.T(), err, "エラーが発生しないこと")
    assert.Equal(s.T(), "テスト圃場", result.Name)
}

// 異常系: 必須フィールド欠如でエラーを返すことを確認
func (s *CreateFieldUseCaseTestSuite) TestExecute_ValidationError_MissingName() {
    input := &usecase.CreateFieldInput{Name: ""}

    result, err := s.useCase.Execute(s.ctx, input)

    require.Error(s.T(), err, "バリデーションエラーが発生すること")
    assert.Nil(s.T(), result)
}
```

## テスト命名規則

```
Test{メソッド}_{種類}_{詳細}
```

| 種類 | 例 |
|------|---|
| Success | `TestExecute_Success_ValidInput` |
| ValidationError | `TestExecute_ValidationError_MissingName` |
| BoundaryValue | `TestExecute_BoundaryValue_MaxNameLength` |
| BusinessLogicError | `TestExecute_BusinessLogicError_DuplicateName` |

## 必須テストケース

- **正常系**: 有効データで成功
- **異常系**: 必須欠如/文字数違反/範囲外/フォーマット違反/不正JSON
- **境界値**: 最小値/最大値/最小値-1/最大値+1

## テストコード品質ルール

- `require.NoError(t, err)` パターンを使用
- `_` でのエラー無視は禁止
- TestMain内: `log.Fatalf` を使用(`fmt.Printf` 禁止)
- 全テストに日本語コメントを記載
- `t.Skip()` でテストをパスさせることは禁止

## カバレッジ目標

| 層 | 目標 |
|---|------|
| Presentation層 | 90% |
| Usecase層 | 95% |
| 重要ビジネスロジック | 100% |
| Domain層 | 100% |
| 生成コード | 除外 |

## テスト実行コマンド

```bash
# 全テスト
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

## 関連コマンド

- `/go-build` — ビルドエラーの修正
- `/go-review` — 実装後のコードレビュー

## 関連

- スキル: `skills/golang-testing/`
- スキル: `skills/golang-patterns/`
