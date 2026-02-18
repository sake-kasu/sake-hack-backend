# Goテスト規約

## フレームワーク

**testify/suite + testify/mock** を使用する(標準テストのtable-drivenではない)。

## テスト命名規則

```
Test{メソッド}_{種類}_{詳細}
```

- `Success` / `ValidationError_{原因}` / `BoundaryValue_{詳細}` / `BusinessLogicError_{原因}`
- **全てのテストに日本語で「何をテストしているのか」のコメントを記載すること**

## ファイル構成

- `*_test.go` — 単体テスト
- `*_integration_test.go` — 統合テスト
- `mock_*_test.go` — モック定義

## カバレッジ目標

| 層 | 目標 |
|---|------|
| Presentation層 | 90% |
| Usecase層 | 95% |
| 重要ビジネスロジック | 100% |
| Domain層 | 100% |
| 生成コード | 除外 |

## 必須テストケース

- **正常系**: 有効データで成功
- **異常系**: 必須欠如/文字数違反/範囲外/フォーマット違反/不正JSON
- **境界値**: 最小値/最大値/最小値-1/最大値+1

## バリデーション責務

- **Presentation層**: フィールドバリデーション(形式/文字数/必須)
- **Usecase層**: ビジネスルール(重複/存在/権限)

## テストコード品質ルール

- `require.NoError(t, err)` パターンを使用
- `_` でのエラー無視は禁止
- TestMain内のログ出力: `log.Fatalf` を使用(`fmt.Printf` 禁止)
- エラーメッセージは日本語で記載
- `t.Skip()` でテストをパスさせることは禁止

## レースディテクション

`-race` フラグでのテスト実行を推奨:

```bash
go test -race ./...
```

## テスト実行コマンド

| コマンド | 説明 |
|---------|------|
| `make test` | 全テスト実行 |
| `make test-unit` | 単体テストのみ |
| `make test-integration` | 統合テストのみ |
| `make cover` | カバレッジレポート |

## 参照

- スキル: `golang-testing` — 詳細なGoテストパターンとヘルパー
