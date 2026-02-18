# Goセキュリティ規約

## シークレット管理

```go
apiKey := os.Getenv("API_KEY")
if apiKey == "" {
    log.Fatal("API_KEYが設定されていません")
}
```

- ハードコードされたシークレット(APIキー、パスワード、トークン)は絶対に禁止
- `.env`、`credentials.json` 等はコミットしない

## セキュリティスキャン

**gosec** による静的セキュリティ解析を必須とする:

```bash
make gosec-scan
```

- 生成コード(`api/generated/`, `internal/database/sqlc/`)は自動除外
- severity=high のみを検出

## `#nosec` 使用禁止(必須)

- `#nosec` アノテーションによる gosec 警告の抑制は**絶対に禁止**
- 警告が出た場合はコードを修正して対応すること
- `//nolint` コメントも同様に禁止

## Context & Timeout

全ての外部呼び出し(DB、HTTP、gRPC等)に `context.Context` を使用:

```go
ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
defer cancel()
```

## SQLインジェクション防止

パラメータ化クエリを必須とする:

```go
// 禁止: 文字列結合によるクエリ
db.Query("SELECT * FROM users WHERE id = " + userID)

// 必須: パラメータ化クエリ
db.Query("SELECT * FROM users WHERE id = $1", userID)
```

## コマンドインジェクション防止

```go
// 禁止: シェル経由のコマンド実行
exec.Command("sh", "-c", "echo " + userInput)

// 必須: 直接引数渡し
exec.Command("echo", userInput)
```

## エラーメッセージのセキュリティ

- エラーメッセージに内部実装の詳細を含めない
- スタックトレースやDB構造をユーザーに露出しない
- 日本語のユーザー向けエラーメッセージを使用

## TLS/暗号

- `InsecureSkipVerify: true` は禁止
- MD5/SHA1 をセキュリティ目的で使用しない
- `unsafe` パッケージは正当な理由がない限り使用しない

## 参照

- セキュリティスキャン: `make gosec-scan`
- エージェント: `go-reviewer` — セキュリティチェック含むコードレビュー
