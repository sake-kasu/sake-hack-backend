---
name: golang-patterns
description: Go + Gin + oapi-codegen + クリーンアーキテクチャにおけるイディオマティックなGoパターン、設計原則、並行処理パターン
---

# Go開発パターン

Go + Gin + oapi-codegen + クリーンアーキテクチャ構成のAPIサーバーにおける設計パターン集。

## 活用タイミング

- 新規Goコードの作成
- Goコードのレビュー
- 既存Goコードのリファクタリング
- Featureパッケージの設計

## コア原則

### 1. シンプルさと明確さ

Goは巧妙さよりもシンプルさを重視する。コードは一目で理解できること。

```go
// 良い: 明確で直接的
func (uc *GetUserUseCase) Execute(ctx context.Context, id string) (*entity.User, error) {
    user, err := uc.repo.FindByID(ctx, id)
    if err != nil {
        return nil, fmt.Errorf("ユーザーの取得に失敗(ID=%s): %w", id, err)
    }
    return user, nil
}
```

### 2. ゼロ値を有用にする

型のゼロ値が初期化なしで使えるように設計する。

```go
// 良い: ゼロ値が有用
type Counter struct {
    mu    sync.Mutex
    count int // ゼロ値は0、すぐに使用可能
}

func (c *Counter) Inc() {
    c.mu.Lock()
    c.count++
    c.mu.Unlock()
}
```

### 3. インターフェースを受け取り、構造体を返す

```go
// 良い: インターフェースを受け取り、具象型を返す
func ProcessData(r io.Reader) (*Result, error) {
    data, err := io.ReadAll(r)
    if err != nil {
        return nil, err
    }
    return &Result{Data: data}, nil
}
```

## Featureパッケージ構成テンプレート

```
internal/features/<feature>/
├── domain/
│   ├── entity/         # エンティティ、Value Object
│   ├── repository/     # リポジトリインターフェース
│   └── dto/            # 機能間データ受け渡し用DTO
├── application/
│   ├── usecase/        # ユースケース実装
│   ├── port/           # 外部サービスインターフェース
│   └── query/          # CQRS読み取り系インターフェース
├── infrastructure/
│   ├── repository/     # リポジトリ実装(DB書き込み)
│   ├── query/          # クエリ実装(DB読み取り)
│   └── external/       # 外部サービス実装
└── presentation/
    └── handler.go      # HTTPハンドラー(ServerInterface実装)
```

## Gin + oapi-codegen ServerInterfaceパターン

```go
// presentation/handler.go
package presentation

// ハンドラー構造体
type FieldHandler struct {
    createUseCase *usecase.CreateFieldUseCase
    getUseCase    *usecase.GetFieldUseCase
    logger        *slog.Logger
}

// コンストラクタ
func NewFieldHandler(
    createUC *usecase.CreateFieldUseCase,
    getUC *usecase.GetFieldUseCase,
    logger *slog.Logger,
) *FieldHandler {
    return &FieldHandler{
        createUseCase: createUC,
        getUseCase:    getUC,
        logger:        logger,
    }
}

// ServerInterface メソッド実装
func (h *FieldHandler) CreateField(c *gin.Context) {
    var req openapi.CreateFieldRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        h.logger.Error("リクエストのバインドに失敗", slog.Any("error", err))
        c.JSON(http.StatusBadRequest, openapi.ErrorResponse{
            Message: "リクエストの形式が不正です",
        })
        return
    }

    result, err := h.createUseCase.Execute(c.Request.Context(), req)
    if err != nil {
        h.handleError(c, err)
        return
    }

    c.JSON(http.StatusCreated, result)
}
```

## Consumer-Providerインターフェースパターン

Feature間でデータを受け渡す場合のパターン:

```go
// === Consumer側: import機能 ===

// domain/dto/field_input.go
package dto

// FieldBatchInput は圃場の一括登録に必要なデータを定義
type FieldBatchInput struct {
    Name     string
    Geometry []byte
    H3Index  string
}

// application/usecase/process_import.go
package usecase

// FieldBatchRepository はConsumer側で定義するProviderインターフェース
type FieldBatchRepository interface {
    UpsertBatch(ctx context.Context, inputs []dto.FieldBatchInput) error
}

type ProcessImportUseCase struct {
    fieldRepo FieldBatchRepository
    logger    *slog.Logger
}

// === Provider側: field機能 ===

// infrastructure/repository/field_repository.go
// Consumer側のDTO型をimport(Domain層同士なので許容)
import importdto "project/internal/features/import/domain/dto"

func (r *FieldRepository) UpsertBatch(ctx context.Context, inputs []importdto.FieldBatchInput) error {
    // 実装
}

// === DI層 ===

// server/router.go
fieldRepo := fieldrepository.NewFieldRepository(db)
processImportUC := usecase.NewProcessImportUseCase(fieldRepo, logger)
```

## エラーハンドリングパターン

### apperrorカスタムエラー型

```go
// internal/apperror/error.go
package apperror

// AppError はアプリケーション固有のエラー型
type AppError struct {
    Code    ErrorCode
    Message string
    Err     error
}

func (e *AppError) Error() string {
    if e.Err != nil {
        return fmt.Sprintf("%s: %v", e.Message, e.Err)
    }
    return e.Message
}

func (e *AppError) Unwrap() error {
    return e.Err
}

// エラーコード定義
type ErrorCode int

const (
    CodeNotFound      ErrorCode = iota
    CodeValidation
    CodeUnauthorized
    CodeConflict
    CodeInternal
)

// エラー生成ヘルパー
func NotFound(msg string) *AppError {
    return &AppError{Code: CodeNotFound, Message: msg}
}

func Validation(msg string) *AppError {
    return &AppError{Code: CodeValidation, Message: msg}
}

func Wrap(err error, msg string) *AppError {
    return &AppError{Code: CodeInternal, Message: msg, Err: err}
}
```

### Presentation層でのエラーハンドリング

```go
func (h *FieldHandler) handleError(c *gin.Context, err error) {
    var appErr *apperror.AppError
    if errors.As(err, &appErr) {
        switch appErr.Code {
        case apperror.CodeNotFound:
            c.JSON(http.StatusNotFound, openapi.ErrorResponse{Message: appErr.Message})
        case apperror.CodeValidation:
            c.JSON(http.StatusBadRequest, openapi.ErrorResponse{Message: appErr.Message})
        case apperror.CodeUnauthorized:
            c.JSON(http.StatusUnauthorized, openapi.ErrorResponse{Message: appErr.Message})
        case apperror.CodeConflict:
            c.JSON(http.StatusConflict, openapi.ErrorResponse{Message: appErr.Message})
        default:
            h.logger.Error("内部エラー", slog.Any("error", err))
            c.JSON(http.StatusInternalServerError, openapi.ErrorResponse{Message: "内部エラーが発生しました"})
        }
        return
    }

    h.logger.Error("予期しないエラー", slog.Any("error", err))
    c.JSON(http.StatusInternalServerError, openapi.ErrorResponse{Message: "内部エラーが発生しました"})
}
```

## slogベースのロギングパターン

```go
// ロガー初期化
logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
    Level: slog.LevelInfo,
}))

// 構造化ログ出力
logger.Info("圃場を作成しました",
    slog.String("field_id", field.ID),
    slog.String("name", field.Name),
)

logger.Error("圃場の作成に失敗",
    slog.Any("error", err),
    slog.String("user_id", userID),
)
```

## エラーラップとコンテキスト

```go
// 良い: 日本語でコンテキスト付きエラーラップ
func (r *FieldRepository) FindByID(ctx context.Context, id string) (*entity.Field, error) {
    field, err := r.queries.GetField(ctx, id)
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return nil, apperror.NotFound("指定された圃場が見つかりません")
        }
        return nil, fmt.Errorf("圃場の取得に失敗(ID=%s): %w", id, err)
    }
    return mapToEntity(field), nil
}
```

## 並行処理パターン

### ワーカープール

```go
func WorkerPool(jobs <-chan Job, results chan<- Result, numWorkers int) {
    var wg sync.WaitGroup

    for i := 0; i < numWorkers; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            for job := range jobs {
                results <- process(job)
            }
        }()
    }

    wg.Wait()
    close(results)
}
```

### Context によるキャンセルとタイムアウト

```go
func FetchWithTimeout(ctx context.Context, url string) ([]byte, error) {
    ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
    defer cancel()

    req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
    if err != nil {
        return nil, fmt.Errorf("リクエストの作成に失敗: %w", err)
    }

    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        return nil, fmt.Errorf("データの取得に失敗(%s): %w", url, err)
    }
    defer resp.Body.Close()

    return io.ReadAll(resp.Body)
}
```

### errgroupによる協調ゴルーチン

```go
import "golang.org/x/sync/errgroup"

func FetchAll(ctx context.Context, urls []string) ([][]byte, error) {
    g, ctx := errgroup.WithContext(ctx)
    results := make([][]byte, len(urls))

    for i, url := range urls {
        i, url := i, url
        g.Go(func() error {
            data, err := FetchWithTimeout(ctx, url)
            if err != nil {
                return err
            }
            results[i] = data
            return nil
        })
    }

    if err := g.Wait(); err != nil {
        return nil, err
    }
    return results, nil
}
```

### ゴルーチンリーク防止

```go
// 悪い: コンテキストキャンセル時にゴルーチンリーク
func leakyFetch(ctx context.Context, url string) <-chan []byte {
    ch := make(chan []byte)
    go func() {
        data, _ := fetch(url)
        ch <- data // レシーバーがなければ永久にブロック
    }()
    return ch
}

// 良い: キャンセルを適切にハンドリング
func safeFetch(ctx context.Context, url string) <-chan []byte {
    ch := make(chan []byte, 1) // バッファ付きチャネル
    go func() {
        data, err := fetch(url)
        if err != nil {
            return
        }
        select {
        case ch <- data:
        case <-ctx.Done():
        }
    }()
    return ch
}
```

## Graceful Shutdown

```go
func GracefulShutdown(server *http.Server) {
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

    <-quit
    slog.Info("サーバーをシャットダウンしています...")

    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    if err := server.Shutdown(ctx); err != nil {
        slog.Error("サーバーの強制シャットダウン", slog.Any("error", err))
    }

    slog.Info("サーバーが正常終了しました")
}
```

## Functional Optionsパターン

```go
type Server struct {
    addr    string
    timeout time.Duration
    logger  *slog.Logger
}

type Option func(*Server)

func WithTimeout(d time.Duration) Option {
    return func(s *Server) {
        s.timeout = d
    }
}

func WithLogger(l *slog.Logger) Option {
    return func(s *Server) {
        s.logger = l
    }
}

func NewServer(addr string, opts ...Option) *Server {
    s := &Server{
        addr:    addr,
        timeout: 30 * time.Second,
        logger:  slog.Default(),
    }
    for _, opt := range opts {
        opt(s)
    }
    return s
}
```

## メモリとパフォーマンス

### スライスの事前割り当て

```go
// 悪い: スライスが複数回拡張される
func processItems(items []Item) []Result {
    var results []Result
    for _, item := range items {
        results = append(results, process(item))
    }
    return results
}

// 良い: 単一割り当て
func processItems(items []Item) []Result {
    results := make([]Result, 0, len(items))
    for _, item := range items {
        results = append(results, process(item))
    }
    return results
}
```

### ループ内の文字列結合回避

```go
// 悪い: 多数の文字列割り当て
func join(parts []string) string {
    var result string
    for _, p := range parts {
        result += p + ","
    }
    return result
}

// 良い: strings.Builder で単一割り当て
func join(parts []string) string {
    var sb strings.Builder
    for i, p := range parts {
        if i > 0 {
            sb.WriteString(",")
        }
        sb.WriteString(p)
    }
    return sb.String()
}

// 最良: 標準ライブラリを使用
func join(parts []string) string {
    return strings.Join(parts, ",")
}
```

## ツール統合

### 必須コマンド

```bash
# ビルド
make build

# テスト
make test
make test-unit
make test-integration
make cover

# 静的解析・Lint
make lint

# セキュリティ
make gosec-scan

# アーキテクチャ検証
make arch-check

# コード生成
make generate
```

## Goイディオム クイックリファレンス

| イディオム | 説明 |
|-----------|------|
| インターフェースを受け取り、構造体を返す | 関数はインターフェースを受け取り具象型を返す |
| エラーは値 | エラーを例外ではなくファーストクラスの値として扱う |
| メモリ共有ではなく通信で共有 | ゴルーチン間の調整にチャネルを使用 |
| ゼロ値を有用に | 型は明示的な初期化なしで動作すべき |
| 少しのコピーは少しの依存より良い | 不要な外部依存を避ける |
| 明確さは巧妙さに勝る | 読みやすさをトリック的なコードより優先 |
| 早期リターン | エラーを先に処理し、ハッピーパスのインデントを浅く |

**心得**: Goコードは予測可能で一貫性があり、理解しやすい「退屈な」コードが最良。迷ったらシンプルに保つこと。
