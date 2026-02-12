# Goパターン規約

## Functional Options

```go
type Option func(*Server)

func WithPort(port int) Option {
    return func(s *Server) { s.port = port }
}

func NewServer(opts ...Option) *Server {
    s := &Server{port: 8080}
    for _, opt := range opts {
        opt(s)
    }
    return s
}
```

## 小さなインターフェース

インターフェースは**使用する側(Consumer)で定義**し、実装する側(Provider)では定義しない。

## 依存性注入

コンストラクタ関数で依存性を注入する:

```go
func NewUserUseCase(repo UserRepository, logger *slog.Logger) *UserUseCase {
    return &UserUseCase{repo: repo, logger: logger}
}
```

## クリーンアーキテクチャ依存関係ルール(必須)

```
Presentation → Application → Domain ← Infrastructure
```

### Feature間の直接import禁止

- `make arch-check` で検証される
- ある機能(Consumer)が別の機能(Provider)のロジックを必要とする場合:
  1. Consumer側の `domain/dto/` で入力型を定義
  2. Consumer側の `application/usecase/` でProviderインターフェースを定義
  3. Provider側の `infrastructure/` がConsumerのDTO型をimport
  4. `server/router.go` でDI結合

### 禁止パターン

- **sharedパッケージ禁止**: `features/shared/` のような共有パッケージは作成しない
- **features外からfeaturesへの参照禁止**: `internal/infrastructure/` 等からfeaturesパッケージをimportしない(`server/` は例外)
- **Infrastructure → Application の参照禁止**
- **型変換はApplication層で実施**

## 参照

- スキル: `golang-patterns` — 包括的なGoパターン(並行処理、インターフェース設計、メモリ最適化等)
- アーキテクチャ検証: `make arch-check`
