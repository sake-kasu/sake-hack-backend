---
description: Goビルドエラー、go vet警告、Linterの問題を段階的に修正する。go-build-resolverエージェントを呼び出す。
---

# Goビルドと修正

**go-build-resolver** エージェントを呼び出してGoビルドエラーを段階的に修正する。

## 実行内容

1. **診断の実行**: `make build`, `make lint`
2. **エラーの解析**: ファイルごとにグループ化し重要度順にソート
3. **段階的修正**: 一つずつエラーを修正
4. **各修正の検証**: 変更のたびにビルドを再実行
5. **サマリーレポート**: 修正内容と残りの問題を報告

## 使用タイミング

- `make build` がエラーで失敗する場合
- `make lint` が問題を報告する場合
- モジュール依存関係が壊れた場合
- ビルドを壊す変更をpullした後

## 実行される診断コマンド

```bash
# ビルドチェック
make build

# Lint(フォーマット + 静的解析)
make lint

# 依存関係修復
make deps

# 生成コードの再生成(必要な場合)
make generate
```

## 重要な注意事項

- **生成コードは編集禁止**
  - `api/generated/` — oapi-codegen生成コード
  - `internal/database/sqlc/` — SQLC生成コード
  - 生成コードに問題がある場合は `make generate` で再生成
- **`//nolint` や `#nosec` の追加禁止**
- **最小限の変更**: リファクタリングせず、エラーのみ修正

## よく修正されるエラー

| エラー | 典型的な修正 |
|-------|------------|
| `undefined: X` | importの追加またはタイプミス修正 |
| `cannot use X as Y` | 型変換または代入の修正 |
| `missing return` | return文の追加 |
| `X does not implement Y` | 不足メソッドの追加 |
| `import cycle` | パッケージの再構成(Consumer側インターフェース定義) |
| `declared but not used` | 変数の削除または使用 |
| `cannot find package` | `make deps` または `go get` |

## 修正戦略

1. **ビルドエラーを先に** — コードがコンパイルされること
2. **Vet警告を次に** — 疑わしい構造を修正
3. **Lint警告を最後に** — スタイルとベストプラクティス
4. **一つずつ修正** — 各変更を検証
5. **最小限の変更** — リファクタリングではなく修正のみ

## 停止条件

以下の場合はエージェントが停止して報告する:
- 同じエラーが3回の試行後も残る
- 修正がより多くのエラーを導入する
- アーキテクチャ変更が必要
- 外部依存関係の手動インストールが必要

## 関連コマンド

- `/go-test` — ビルド成功後にテストを実行
- `/go-review` — コード品質をレビュー

## 関連

- エージェント: `agents/go-build-resolver.md`
- スキル: `skills/golang-patterns/`
