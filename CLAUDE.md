# CLAUDE.md

このファイルは、Claude Code (claude.ai/code) がこのリポジトリで作業する際のガイダンスを提供します。

## 前提

ユーザーは Go のスペシャリストですが、作業の効率化を図るためにあなたの作業を依頼しています。
Claude も Go のスペシャリストです。修正は Go のベストプラクティスに沿って実行してください。

## 絶対に守るルール

### 対話ルール

Claude がユーザーと対話する場合は必ず日本語で行ってください。思考プロセスにおいて英語の方が都合が良ければ英語で思考して構いません。
全ての出力に対して，全角の（）の使用を禁じます。半角の()を使用してください。

## プロジェクト概要

Go + クリーンアーキテクチャ + Package by Feature + OpenAPI-First 開発
**技術スタック**: Go 1.25+ / Gin / sqlc v1.29.0 / OpenAPI 3.0 / testify / PostgreSQL(5432) / Valkey(6379) / JWT 認証(将来 OAuth 2.0 対応予定)
**規模**: 初期開発段階 / 機能は開発に応じて柔軟に追加 / メイン DB(PostgreSQL) + Valkey(セッション/キャッシュ)

## 開発コマンド

```bash
# ビルド・実行
make build / make run / make dev / make clean

# テスト・品質
make test / make cover / make lint / make gosec-scan

# 依存関係
make deps  # go mod tidy

# API開発(OpenAPI仕様から自動生成)
make api-validate / make api-generate / make api-bundle / make api-gendoc

# データベース開発
make sqlc-generate                # sqlcコード生成
make db-migrate-up                # マイグレーション実行
make db-migrate-down              # マイグレーションロールバック
make db-migrate-create NAME=xxx   # 新規マイグレーション作成
```

## アーキテクチャ

### Package by Feature + クリーンアーキテクチャ

**機能構成**: 現在は基盤のみ実装済み / internal/features/ディレクトリに機能を柔軟に追加予定

```
internal/
├── features/                # Package by Feature(ドメイン機能)
│   └── <feature_name>/
│       ├── application/
│       │   ├── port/        # 外部サービスIF定義
│       │   ├── query/       # Query IF / Read Model型定義(CQRS読取)
│       │   └── usecase/     # ビジネスロジック ※Request/Response型はgenerated/から
│       ├── domain/
│       │   ├── entity/      # Domain Entity(純粋BL)
│       │   └── repository/  # Repository IF定義
│       ├── infrastructure/
│       │   ├── external/    # 外部サービス実装
│       │   ├── query/       # Query実装(sqlc使用・読取専用)
│       │   └── repository/  # Repository実装(sqlc使用・書込)
│       └── presentation/    # server_impl.go(oapi-codegen) / *_test.go
├── database/                # sqlc生成コード / postgres.go / valkey.go
│   └── sqlc/                # models.go / querier.go / *.sql.go(自動生成)
├── middleware/              # cors.go / request_id.go
├── apperror/                # エラー定義
├── logger/                  # ロギング
├── server/                  # サーバー設定
└── utils/                   # ユーティリティ(safe_conv.go等)
api/
└── generated/               # OpenAPI自動生成(手動編集禁止)
db/
├── migrations/              # マイグレーションファイル
└── queries/                 # sqlcクエリ定義
```

**アーキテクチャ原則**:

1. 依存性逆転(Domain が IF 定義)
2. API-First(OpenAPI 駆動)
3. CQRS(読取=Query / 書込=Repository+Entity)
4. DTO 分離(Request/Response vs DomainEntity)
5. 機能別 ServerInterface
6. 型安全なデータアクセス(sqlc による SQL → Go コード生成)
7. マイグレーション管理(golang-migrate)
8. JWT 認証(将来 OAuth 2.0 対応予定)

## 開発ワークフロー(新機能追加)

**手順**:

1. **API 仕様**: `api/paths/*.yaml`, `api/components/schemas/*.yaml`定義 → `make api-validate` → `make api-generate`
2. **パッケージ作成**: `mkdir -p internal/features/<name>/{application/{query,usecase,port},domain/{entity,repository},infrastructure/{query,repository,external},presentation}`
3. **DB 設計**:
   - マイグレーション作成: `make db-migrate-create NAME=create_xxx_table`
   - `db/migrations/XXXXXX_create_xxx_table.up.sql`にスキーマ定義
   - `db/migrations/XXXXXX_create_xxx_table.down.sql`にロールバック定義
   - 実行: `make db-migrate-up`
4. **SQL クエリ定義**:
   - `db/queries/xxx.sql`にクエリ定義(sqlc アノテーション付き)
   - コード生成: `make sqlc-generate`
   - 生成コード: `internal/database/sqlc/`
5. **実装**: Domain(Entity+RepoIF) → Application(QueryIF+Usecase) → Infrastructure(Query/Repo 実装) → Presentation(ServerIF)
6. **テスト**: 各レイヤーで`*_test.go`(単体)/`*_integration_test.go`(統合)作成
7. **DI 登録**: `cmd/server/main.go`または`internal/server/server.go`に追加
8. **検証**: `make test` → `make cover` → `make lint` → `make gosec-scan` → `make build`

## 実装ガイドライン

**基本ルール**:

- OpenAPI 変更後は`make api-generate`必須(`api/generated/`は手動編集禁止)
- Request/Response 型は`api/generated/server.gen.go`から使用
- CQRS 適用(読取=Query / 書込=Repository)
- データアクセスは sqlc 使用(SQL → Go コード生成)
- マイグレーションは golang-migrate で管理
- 外部サービスは`application/port/`(IF) + `infrastructure/external/`(実装)
- エラー処理は`internal/apperror`使用
- ロギングは`internal/logger`使用(メソッドトレース推奨)
- API レスポンスは`null`でなく`[]`返却
- 実装と同時にテスト作成必須

**DB・Valkey 選択**:

- メイン DB(PostgreSQL): アプリ固有データ / sqlc 生成コード(`internal/database/sqlc/`) / 接続管理(`internal/database/postgres.go`)
- Valkey: セッション/認証/キャッシュ / `internal/database/valkey.go`

**設定管理** (`.env`):

- **設定方式**: 環境変数 + `.env`ファイル (caarlos0/env + godotenv)
- **設定構造体**: 機能別に分割 (server.go, database.go, cache.go, jwt.go, cors.go, logger.go)
- Server: SERVER_PORT(8080) / SERVER_MODE(debug) / SERVER_GRACEFUL_SHUTDOWN_TIMEOUT(30s)
- JWT: JWT_SECRET(必須) / JWT_EXPIRATION(86400) / JWT_COOKIE_*設定 ※将来実装予定
- OAuth: 将来実装予定(CIOS または汎用 OAuth 2.0)
- Cache(Valkey): CACHE_HOST/PORT/PASSWORD/DATABASE/POOL_SIZE/タイムアウト設定
- Database(PostgreSQL): DB_HOST/PORT/NAME/USER/PASSWORD(必須) / DB_MAX_OPEN_CONNS / DB_CONN_MAX_LIFETIME
- CORS: CORS_ALLOWED_ORIGINS(カンマ区切り) / ALLOWED_METHODS / ALLOWED_HEADERS / ALLOW_CREDENTIALS
- Logging: LOG_LEVEL(debug) / LOG_FORMAT(console)

**注意**: `.env`ファイルは機密情報を含むため.gitignoreに含まれる。開発時は`.env.sample`をコピーして使用。本番環境では環境変数で直接設定。

## コーディング規約

**命名規則**:

- DTO: OpenAPI 生成型(`generated/types.gen.go`) / Application 内部型(`~Param`, `~Input`)
- ファイル: 責務別分離(例: `user.go`=ユーザー関連 / `project.go`=プロジェクト関連)
- SQL ファイル: 機能別分離(`db/queries/users.sql`, `db/queries/projects.sql`)

**責務分離**:

- データアクセス: sqlc 生成コード(`internal/database/sqlc/`) + Repository 層でラップ
- 外部サービス: `application/port/`(IF 定義) + `infrastructure/external/`(実装)
- Usecase: 1 ファイル=1 責務(単一責任の原則)
- API 設計一貫性: API パスと実装の責務を一致

**実装注意点**:

1. OpenAPI 仕様 →`make api-generate`→ 実装の順序厳守
2. `api/generated/server.gen.go`に型が生成されていることを確認
3. SQL 作成 →`make sqlc-generate`→ 生成コード確認
4. マイグレーション作成 →`make db-migrate-up`→ 適用確認
5. 新機能作成時は`cmd/server/main.go`または`internal/server/server.go`に登録
6. 実装と同時にテスト作成(TDD 推奨)
7. 変更後は`make test`、`make lint`、`make gosec-scan`実行
8. 安全な型変換実施(後述)

### 型変換のセキュリティガイドライン(CWE-190 Integer Overflow 対策)

**必須ルール**:

1. **直接キャスト禁止**: `int32(value)` は使用禁止
2. **必ず`utils.IntToInt32()`/`utils.Int64ToInt32()`使用**(`internal/utils/safe_conv.go`)
3. **エラーハンドリング必須**: レイヤーごとに適切な処理
   - Presentation 層: HTTP エラーレスポンス(`apperror.BadRequestError`)
   - Repository 層: ログ出力 + フォールバック値
   - Parser 層: エラーメッセージ返却

**使用例**:

```go
// ✅ 正しい変換
difficulty, err := utils.IntToInt32(int(req.Difficulty))
if err != nil {
    return nil, apperror.BadRequestError("難易度の値が不正です")
}

// ✅ strconv使用時
id64, err := strconv.ParseInt(idStr, 10, 32)  // 第3引数で範囲指定
return int32(id64), nil

// ❌ 絶対禁止
return int32(value)  // gosec G115/G109エラー
```

**検証**: `make gosec-scan` / `gosec -include=G109,G115 ./...`

**修正優先度**: 1.外部入力変換 > 2.DB 値変換 > 3.内部固定値

### Value Object の作成方針

**特性**: 不変性 / 値による等価性 / 自己検証

**必須実装**(型エイリアスだけでは不十分):

1. **型定義**: `type TransactionType string` + const 定義
2. **コンストラクタ**(バリデーション付き): `func NewTransactionType(value string) (*TransactionType, error)`
3. **バリデーションメソッド**: `func (t TransactionType) IsValid() bool`

**レイヤー別エラーハンドリング**:

- Presentation 層: `apperror.BadRequestError(err.Error())` または `apperror.NewValidationError("入力値が不正です").AddField("field", err.Error())`
- Repository 層: `apperror.DatabaseError("操作に失敗しました", err)` でエラーをラップ
- CSV Parser 層: `apperror.NewValidationError("CSV解析エラー").AddField(fmt.Sprintf("row_%d", rowNum), err.Error())`

**必須ルール**:

1. 直接キャスト禁止: `entity.TransactionType(*s)` 使用不可
2. 必ずコンストラクタ使用: `entity.NewTransactionType(*s)`
3. すべての入力経路(API/CSV/DB)で統一バリデーション
4. エラーハンドリング必須
5. 正常値・不正値両方のテスト作成

## テスト実装ガイドライン

**必須方針**: 実装と同時にテスト作成

**フレームワーク**: testify/mock + httptest

**ファイル構成**: `*_test.go`(単体) / `mock_*_test.go`(モック)

**カバレッジ要件**: Presentation 層=90% / Usecase 層=95% / 重要 BL=100% / Domain 層=100%

**必須テストケース**:

- 正常系: 有効データで成功
- 異常系: 必須欠如/文字数違反/範囲外/フォーマット違反/不正 JSON
- 境界値: 最小値/最大値/最小値-1/最大値+1

**命名規則**: `Test{メソッド}_{種類}_{詳細}` (Success / ValidationError*{原因} / BoundaryValue*{詳細} / BusinessLogicError\_{原因})

**バリデーション責務**:

- Presentation 層: フィールドバリデーション(形式/文字数/必須)
- Usecase 層: ビジネスルール(重複/存在/権限)

**テスト実行コマンド**:

```bash
make test              # 全テスト実行
make test-unit         # 単体テストのみ実行(-short / 高速)
make cover             # カバレッジ測定
```

**開発フロー**: OpenAPI 定義 → テストケース設計 → テスト雛形作成(TDD) → 実装 → テスト全通 → カバレッジ確認 → リンター

## セキュリティ

**機密ファイル読取禁止**: `.env`ファイル(機密情報含む)は.gitignoreに含まれる

**開発時**: `.env.sample`をコピー → `.env`を作成

**本番**: 環境変数で直接設定 / 暗号化値使用 / `JWT_COOKIE_SECURE=true`必須

## Git コミット

**フォーマット**: `{type}:{emoji}{対象の説明}(#チケット番号)`

**Type**: add(新機能) / fix(バグ) / update(改善) / refactor / docs / test / style / chore / remove

**Emoji**: 📝(ドキュメント) / 🐛(バグ) / ⚡(改善) / ♻️(リファクタ) / 📚(ドキュメント) / 🧪(テスト) / 🎨(UI) / 🔧(設定) / 🗑️(削除) / ✨(新機能) / 🔒(セキュリティ)

**重要ルール**:

- "Generated with Claude Code" / "Co-Authored-By: Claude" 含めない
- 変更内容のみ記述(理由・影響明記)
- チケット番号=ブランチ名`feature/xxxxx`の`xxxxx`部分

**粒度**: 1 機能=1 コミット / ファイル種別別 / 影響範囲別 / WIP 禁止

**タイミング**: 機能完了時 / 設定変更完了時 / ドキュメント更新完了時 / バグ修正完了時 / リファクタ完了時

## ログ運用

**メソッドトレース**: `defer logger.TraceMethodAuto(ctx, param)()`(request_id/goroutine_id/method/phase/duration_ms 記録)

**エラーログ分類**:

- Database: `logger.LogDatabaseError(ctx, "CREATE", "projects", err, ...)`
- Business: `logger.LogBusinessError(ctx, "ルール名", err, ...)`
- Validation: `logger.LogValidationError(ctx, "field", value, "required")`

**ログレベル**: DEBUG(メソッドトレース) / INFO(HTTP/成功) / WARN(Business/Validation) / ERROR(System/DB)

**解析手順**:

1. HTTP ステータス ≥400 でエラー特定 → request_id 抽出
2. request_id で全ログフィルタ → 時系列ソート
3. メソッドトレースで失敗箇所特定
4. エラータイプ別詳細で原因分析

**jq 例**:

```bash
# エラー特定
jq 'select(.status >= 400) | {request_id, status, path, error}' app.log
# リクエスト追跡
jq 'select(.request_id == "xxx")' app.log | sort_by(.time)
# DBエラー集計
jq 'select(.error_type == "database") | {table, operation, db_error_code}' app.log | jq -s 'group_by(.db_error_code) | map({error_code: .[0].db_error_code, count: length})'
```

**保管**: ローカル=標準出力 / 本番=CloudWatch・ELK / 保存期間=エラー 30 日・アクセス 7 日 / ローテーション=日次・最大 10GB

## エラーハンドリング戦略

### 基本方針

**エラーは失われてはならない**: すべてのエラーは適切にラップし、コンテキスト情報を付与してログに記録する。

**レイヤー別責務**:

- **Presentation 層**: HTTP リクエストバリデーション、AppError → HTTP レスポンス変換
- **Application 層**: ビジネスロジックエラー、外部エラー → AppError 変換
- **Infrastructure 層**: データベース/外部 API エラー → AppError ラップ
- **Domain 層**: ビジネスルール違反を error で返却

### エラー種別と使い分け

| エラー種別            | 使用場面                   | HTTP ステータス | 例                                   |
| --------------------- | -------------------------- | --------------- | ------------------------------------ |
| `BadRequestError`     | 不正なリクエスト           | 400             | 必須パラメータ欠如、フォーマット違反 |
| `UnauthorizedError`   | 認証失敗                   | 401             | トークン無効、ログイン必須           |
| `ForbiddenError`      | 権限不足                   | 403             | リソースへのアクセス権なし           |
| `NotFoundError`       | リソース不存在             | 404             | ユーザー/プロジェクト未発見          |
| `ConflictError`       | 競合                       | 409             | 一意制約違反、楽観的ロック失敗       |
| `DatabaseError`       | DB エラー                  | 500             | クエリ失敗、接続エラー               |
| `InternalServerError` | その他システムエラー       | 500             | 予期しないエラー                     |
| `ValidationError`     | 単純なバリデーションエラー | 400             | 簡易チェック                         |
| `NewValidationError`  | フィールド別バリデーション | 400             | 複数フィールドのエラー               |

### エラーラップのベストプラクティス

**1. データベースエラー**:

```go
user, err := r.queries.CreateUser(ctx, params)
if err != nil {
    // エラーをラップして返す
    return nil, apperror.DatabaseError("ユーザー作成に失敗しました", err)
}
```

**2. 外部 API エラー**:

```go
resp, err := client.FetchUserInfo(token)
if err != nil {
    return nil, apperror.InternalServerError("外部API呼び出しに失敗しました").
        WithErr(err).
        WithDetails("api", "user_info")
}
```

**3. ビジネスロジックエラー**:

```go
// Domainでエラー定義
var ErrProjectNotFound = errors.New("プロジェクトが見つかりません")

// Usecaseで変換
project, err := r.repo.FindByID(ctx, id)
if errors.Is(err, domain.ErrProjectNotFound) {
    return nil, apperror.NotFoundError("プロジェクトが見つかりません").
        WithDetails("project_id", id)
}
```

**4. フィールド別バリデーション**:

```go
func validateCreateUser(req *CreateUserRequest) error {
    verr := apperror.NewValidationError("ユーザー作成リクエストが不正です")

    if req.Name == "" {
        verr.AddField("name", "名前は必須です")
    }
    if req.Email == "" {
        verr.AddField("email", "メールアドレスは必須です")
    }
    if req.Age < 0 || req.Age > 150 {
        verr.AddField("age", "年齢は0〜150の範囲で入力してください")
    }

    if verr.HasErrors() {
        return verr
    }
    return nil
}
```

### エラーログとの連携

```go
import "your-project/internal/logger"

// Repositoryでのエラー
user, err := r.queries.GetUser(ctx, id)
if err != nil {
    logger.LogDatabaseError(ctx, "SELECT", "users", err, map[string]interface{}{
        "user_id": id,
    })
    return nil, apperror.DatabaseError("ユーザー取得に失敗しました", err)
}

// Usecaseでのビジネスエラー
if project.Status != "active" {
    logger.LogBusinessError(ctx, "inactive_project_access", nil, map[string]interface{}{
        "project_id": project.ID,
        "status":     project.Status,
    })
    return apperror.BadRequestError("非アクティブなプロジェクトにはアクセスできません")
}
```

### エラーエスカレーション

**レイヤー別エラー型の原則**:

| レイヤー       | 返すエラー型        | 変換ルール                                      |
| -------------- | ------------------- | ----------------------------------------------- |
| Domain         | 標準 `error`        | `var ErrUserNotFound = errors.New("...")` 定義  |
| Infrastructure | `apperror.AppError` | DB/外部 API エラー → `apperror.*Error()` に変換 |
| Application    | `apperror.AppError` | Domain エラー → `apperror.*Error()` に変換      |
| Presentation   | HTTP レスポンス     | `errors.As()` で判定 → JSON 変換                |

**実装例**:

```go
// Domain層：標準errorのみ
var ErrUserNotFound = errors.New("user not found")

// Infrastructure層：DB/外部APIエラー → AppError変換
if errors.Is(err, pgx.ErrNoRows) {
    return nil, apperror.NotFoundError("ユーザーが見つかりません").WithDetails("user_id", id)
}

// Application層：Domainエラー → AppError変換
user, err := entity.NewUser(input.Name, input.Email)
if errors.Is(err, entity.ErrEmailRequired) {
    return nil, apperror.BadRequestError("メールアドレスは必須です")
}

// Presentation層：AppError → HTTPレスポンス
var appErr *apperror.AppError
if errors.As(err, &appErr) {
    c.JSON(appErr.Status, appErr)
    return
}
```

**重要ルール**:

- Domain 層で `apperror` を使わない(依存性逆転違反)
- すでに `AppError` の場合は再ラップ不要
- 各レイヤーで適切なログ記録必須

### エラーハンドリング

**errors.Is** でエラー種別を判定:

```go
if errors.Is(err, pgx.ErrNoRows) {
    return apperror.NotFoundError("レコードが見つかりません")
}
```

**errors.As** でエラー型を取得:

```go
var appErr *apperror.AppError
if errors.As(err, &appErr) {
    log.Printf("App Error: code=%s, status=%d", appErr.Code, appErr.Status)
}
```

**errors.Unwrap** で元のエラーを取得:

```go
unwrapped := errors.Unwrap(appErr)
// unwrapped は元の pgx エラーなど
```

## 開発チェックリスト

**新機能追加**:

- OpenAPI 定義 → `make api-validate` → `make api-generate` → マイグレーション作成(`make db-migrate-create`) → SQL クエリ定義 → `make sqlc-generate` → 機能パッケージ作成 → 4 層実装(Domain→Application→Infrastructure→Presentation) → テスト作成 → DI 登録(`cmd/server/main.go` or `internal/server/server.go`) → `make test/cover/lint/gosec-scan/build`

**コーディング**:

- `api/generated/server.gen.go`から型使用 / sqlc 生成コード(`internal/database/sqlc/`)使用 / 責務別ファイル分離 / 外部サービス=`application/port/` + `infrastructure/external/` / `logger.TraceMethodAuto(ctx, param)()` / エラーログ分類(Database/Business/Validation) / 空配列`[]`返却 / 安全な型変換(`utils.IntToInt32()`) / コミット=`{type}:{emoji}{説明}(#xxx)`

**セキュリティ**:

- `.env`ファイル読取禁止(.gitignore済み) / `.env.sample`参照 / 本番=`JWT_COOKIE_SECURE=true` / 暗号化値使用 / `make gosec-scan`で脆弱性チェック

**構造理解**: 初期開発段階 / 機能は features/に柔軟に追加 / メイン DB(PostgreSQL) + Valkey / sqlc + OpenAPI 駆動 / クリーンアーキテクチャ
