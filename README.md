# sake-hack-backend

Go + クリーンアーキテクチャ + Package by Feature + OpenAPI-First 開発

## 技術スタック

- **言語**: Go 1.25.4
- **Web フレームワーク**: Gin
- **データベース**:
  - PostgreSQL 18 + PostGIS 3.6
  - Valkey 9.0.0 (セッション/キャッシュ)
- **API 仕様**: OpenAPI 3.0
- **SQL**: sqlc (型安全な SQL クエリ生成)
- **マイグレーション**: golang-migrate
- **ロギング**: zap
- **テスト**: testify, TestContainers
- **開発ツール**: Air (ホットリロード)

## アーキテクチャ

### Package by Feature + クリーンアーキテクチャ

```
internal/
├── features/              # 機能別パッケージ
│   ├── sake/
│   │   ├── domain/        # Entity + Repository IF
│   │   │   ├── entity/    # ドメインエンティティ
│   │   │   └── repository/# リポジトリインターフェース
│   │   ├── application/   # Usecase + Query IF
│   │   │   ├── usecase/   # ユースケース実装
│   │   │   └── query/     # クエリインターフェース(CQRS)
│   │   ├── infrastructure/# Repository実装 + Query実装
│   │   │   ├── repository/# リポジトリ実装
│   │   │   └── query/     # クエリ実装
│   │   └── presentation/  # HTTPハンドラ
│   │       ├── server_impl.go  # Ginハンドラ
│   │       ├── converter.go    # DTO変換
│   │       └── validator.go    # バリデーション
│   └── (他の機能...)
├── middleware/            # HTTP middleware
├── logger/                # 構造化ログ
├── server/                # サーバー設定
├── apperror/              # カスタムエラー
├── utils/                 # ユーティリティ
├── api/                   # Openapi
├── database/              # DB接続管理
└── docs/                  # 要件資料（サブモジュール）
```

#### アーキテクチャの特徴

1. **クリーンアーキテクチャ**: 依存性逆転の原則に基づいた層分離
2. **CQRS (Command Query Responsibility Segregation)**: 読み取り(Query)と書き込み(Command)の分離
3. **DI (Dependency Injection)**: インターフェースによる疎結合
4. **Package by Feature**: 機能単位でのパッケージング

#### 実装例: GET /sakes API のデータフロー

```
HTTPリクエスト (GET /sakes?category=JAPANESE_SAKE&limit=20)
    ↓
┌─────────────────────────────────────────────────────────────┐
│ [1] Presentation Layer (プレゼンテーション層)                │
│ internal/features/sake/presentation/server_impl.go          │
│ - GetSakes(): HTTPリクエスト受信                             │
│ - パラメータ抽出・バリデーション                              │
│ - ListSakesInputに変換                                       │
└─────────────────────────────────────────────────────────────┘
    ↓
┌─────────────────────────────────────────────────────────────┐
│ [2] Application Layer (アプリケーション層)                   │
│ internal/features/sake/application/usecase/list_sakes.go    │
│ - ListSakesUsecase.Execute()                                │
│ - ビジネスルール適用 (limit: 1-100, default: 20)            │
│ - ListPublicSakesFilterに変換                               │
└─────────────────────────────────────────────────────────────┘
    ↓
┌─────────────────────────────────────────────────────────────┐
│ [3] Query Interface (クエリインターフェース - CQRS)          │
│ internal/features/sake/application/query/sake_query.go      │
│ - SakeQuery.ListPublic() インターフェース                   │
└─────────────────────────────────────────────────────────────┘
    ↓
┌─────────────────────────────────────────────────────────────┐
│ [4] Infrastructure Layer (インフラストラクチャ層)            │
│ internal/features/sake/infrastructure/query/sake_query_impl.go│
│ - sakeQueryImpl.ListPublic()                                │
│ - COUNT クエリで総件数取得                                   │
│ - SELECT + JOIN で酒データ取得                              │
│   (sakes ⋈ sake_kinds ⋈ breweries)                         │
└─────────────────────────────────────────────────────────────┘
    ↓
┌─────────────────────────────────────────────────────────────┐
│ [5] Database Layer (データベース層)                          │
│ db/queries/sakes.sql (SQLC定義)                             │
│ PostgreSQL + PostGIS                                        │
└─────────────────────────────────────────────────────────────┘
    ↓ entity.SakeDetail に変換
┌─────────────────────────────────────────────────────────────┐
│ [6] Domain Layer (ドメイン層)                                │
│ internal/features/sake/domain/entity/sake.go                │
│ - SakeDetail: 酒の詳細情報エンティティ                       │
│ - Pagination: ページネーション情報                           │
└─────────────────────────────────────────────────────────────┘
    ↓ ListSakesOutput (Sakes + Total)
┌─────────────────────────────────────────────────────────────┐
│ [7] Presentation Layer (プレゼンテーション層)                │
│ internal/features/sake/presentation/converter.go            │
│ - toListSakesResponse(): DTO変換                            │
│ - entity.SakeDetail → generated.SakeDetail                  │
└─────────────────────────────────────────────────────────────┘
    ↓
HTTPレスポンス (JSON)
```

#### レイヤー別の責務

| レイヤー           | 責務                                               | 依存方向                 |
| ------------------ | -------------------------------------------------- | ------------------------ |
| **Presentation**   | HTTP入出力、バリデーション、DTO変換                | → Application            |
| **Application**    | ビジネスロジック、ユースケース調整                 | → Domain, Query IF       |
| **Domain**         | ビジネスルール、エンティティ、バリューオブジェクト | 依存なし(中心)           |
| **Infrastructure** | DB操作、外部API呼び出し                            | → Domain, Application IF |

#### CQRS パターンの採用

読み取り(Query)と書き込み(Command)を分離することで、以下のメリットを実現:

- **Query**: `application/query/` + `infrastructure/query/`
  - 複雑な検索条件に最適化
  - JOIN、集計、キャッシュなどに特化
  - 例: `ListPublic()`, `GetDetail()`

- **Command**: `domain/repository/` + `infrastructure/repository/`
  - トランザクション管理
  - ビジネスルール保証
  - 例: `Create()`, `Update()`, `Delete()`

この分離により、読み取りと書き込みで異なる最適化戦略を採用できます。

## ローカル環境セットアップ

### 前提条件

- Go 1.26.0
- Docker & Docker Compose
- make

### クイックスタート（初回セットアップ）

```bash
# 1. リポジトリをクローン
git clone git@github.com:sake-kasu/sake-hack-backend.git
cd sake-hack-backend

# 2. サブモジュール初期化
# → docs/ディレクトリに要件定義書・議事録をダウンロード
make submodule-init

# 3. 依存関係インストール
# → Go modulesの依存パッケージを取得
make deps

# 4. Docker起動
# → PostgreSQL + Valkey + pgweb を起動
# → config.yamlから自動的にdocker/.env.localを生成
make docker-up

# 5. マイグレーション実行
# → DBスキーマを作成（初回のみmigrate-installが必要）
make migrate-install
make migrate-up

# 6. テストデータ投入（オプション）
# → 酒の種類、メーカー、酒データ、在庫データを投入
docker exec -i sake-hack-db psql -U postgres -d sake_hack_app < test_data.sql

# 7. コード生成
# → OpenAPI仕様からHTTPハンドラ、SQLからDB操作コードを自動生成
make generate

# 8. APIサーバー起動
# → http://localhost:8080 で起動
make run
```

**動作確認（別ターミナルで実行）：**

```bash
# ヘルスチェック
curl http://localhost:8080/health

# 酒一覧取得
curl "http://localhost:8080/api/sakes?limit=10" | jq

# 酒詳細取得（IDは test_data.sql のデータを参照）
curl "http://localhost:8080/api/sakes/{sake_id}" | jq
```

### トラブルシューティング

#### ポートが使用中の場合

```bash
# config.yamlでポート変更
vim config/config.yml
# database.port: 5432 → 5433 に変更

# Docker再起動
make docker-down
make docker-up

# マイグレーション（新しいポートで）
DB_PORT=5433 make migrate-up
```

#### Docker環境のクリーンアップ

```bash
# コンテナ停止・削除
make docker-down

# ボリュームも含めて完全削除
cd docker && docker compose down -v
cd ..
```

#### マイグレーションのやり直し

```bash
# 全ロールバック
make migrate-down-all

# 再適用
make migrate-up
```

### 2回目以降の起動

```bash
# Docker起動
make docker-up

# APIサーバー起動
make run
```

設定ファイルやマイグレーションは保持されているため、すぐに開発を開始できます。

## 開発コマンド

### ビルド・実行

```bash
make build              # ビルド
make run                # 実行
make dev                # 開発サーバー起動(ホットリロード)
make clean              # ビルド成果物削除
```

### テスト・品質

```bash
make test               # 全テスト実行
make test-unit          # 単体テストのみ
make test-integration   # 統合テストのみ
make cover              # カバレッジ測定
make lint               # リンター実行
make gosec-scan         # セキュリティスキャン
```

### サブモジュール

```bash
make submodule-init     # サブモジュールを初期化
make submodule-update   # サブモジュールを最新に更新
make submodule-status   # サブモジュールの状態確認
```

> **Note**: `docs/` ディレクトリは [sake-docs](https://github.com/sake-kasu/sake-docs) リポジトリからサブモジュールとして管理されています。
> プロジェクトドキュメント（要件定義書、議事録など）は `docs/` ディレクトリで確認できます。

### API 開発

```bash
make api-validate       # OpenAPI仕様検証
make api-generate       # コード生成
make api-bundle         # OpenAPI仕様バンドル
make api-gendoc         # APIドキュメント生成
make api-watch          # APIドキュメント生成（リアルタイム反映）
```

### データベース

```bash
make db-migrate-up              # マイグレーション実行
make db-migrate-down            # マイグレーションロールバック
make db-migrate-create NAME=xxx # 新規マイグレーション作成
make sqlc-generate              # sqlcコード生成
```

## CI/CD

### GitHub Actions ワークフロー

#### developへのPR時 (`.github/workflows/pr-develop.yml`)

PR作成時に自動実行されるパイプライン:

- **Validate**: OpenAPI仕様検証、Go modules検証、生成コード検証
- **Build**: バイナリビルド、Dockerイメージビルド(変更時のみ)
- **Test**: 単体テスト、統合テスト + カバレッジ測定
- **Quality**: golangci-lint によるコード品質チェック
- **Security**:
  - gosec: Goセキュリティスキャン
  - gitleaks: シークレット検出
  - govulncheck: 依存関係の脆弱性スキャン
  - license-check: ライセンス適合性チェック
  - hadolint: Dockerfileセキュリティチェック

#### developブランチへのpush時 (`.github/workflows/develop.yml`)

developブランチへのマージ時に自動実行されるパイプライン:

- **Build & Push**: Dockerイメージをビルドし、GitHub Container Registryにプッシュ
  - タグ: `develop`, `daily`, `sha-<commit-hash>`
- **Test**: 単体テスト、統合テスト + カバレッジ測定
- **Documentation**: APIドキュメント生成(Redocly)
- **Deploy**: GitHub Pagesへドキュメントをデプロイ
  - API仕様書、テストカバレッジレポート

#### mainブランチへのpush時 (`.github/workflows/main.yml`)

mainブランチへのマージ時に自動実行されるパイプライン:

- **Build & Push**: Dockerイメージをビルドし、GitHub Container Registryにプッシュ
  - タグ: `main`, `latest`, `sha-<commit-hash>`
- **Test**: 単体テスト、統合テスト + カバレッジ測定

### 必要な設定

GitHub Actionsを使用するために、以下の設定が必要です:

1. **GitHub Container Registryの有効化**:
   - リポジトリ設定 → Packages → Container registry を有効化
   - ワークフローに `packages: write` 権限を付与済み

2. **GitHub Pagesの有効化** (develop用):
   - リポジトリ設定 → Pages → Source を「GitHub Actions」に設定
   - developブランチへのpush時に自動デプロイ
   - 公開URL: `https://<username>.github.io/<repository>/`

3. **カバレッジ閾値の設定** (オプション):
   - デフォルト: 80%
   - 変更する場合は、各ワークフローの `COVERAGE_THRESHOLD` 環境変数を編集

## プロジェクト構造

```
sake-hack-backend/
├── cmd/
│   └── server/
│       └── main.go            # エントリーポイント
├── internal/
│   ├── features/              # 機能別パッケージ
│   ├── middleware/            # HTTPミドルウェア
│   ├── logger/                # ロギング
│   ├── server/                # サーバー設定
│   ├── apperror/              # エラー定義
│   ├── utils/                 # ユーティリティ
│   └── database/              # DB接続管理
├── api/
│   ├── openapi.yaml           # OpenAPI定義(メインファイル)
│   ├── openapi.bundled.yaml   # バンドル済み仕様(自動生成)
│   ├── paths/                 # エンドポイント定義
│   │   └── health.yaml
│   ├── components/            # 再利用可能なコンポーネント
│   │   ├── schemas/           # データモデル定義
│   │   ├── responses/         # レスポンス定義
│   │   ├── parameters/        # パラメータ定義
│   │   └── securitySchemes/   # 認証スキーム定義
│   ├── redocly.yaml           # バリデーション設定
│   ├── oapi-codegen.yaml      # コード生成設定
│   └── generated/             # 自動生成コード(コミット対象)
├── db/
│   ├── migrations/            # マイグレーションファイル
│   └── queries/               # sqlc用SQLファイル
├── config/
│   └── config.yml.sample      # 設定テンプレート
├── docker/
│   └── compose.yaml           # Docker Compose設定
├── Makefile
├── sqlc.yaml
└── README.md
```

## 開発ワークフロー

### 新機能追加

1. **API 仕様定義**:
   - エンドポイント: `api/paths/<endpoint_name>.yaml` に追加
   - スキーマ: `api/components/schemas/<schema_name>.yaml` に追加
   - レスポンス: `api/components/responses/` で共通レスポンスを再利用
   - メインファイル: `api/openapi.yaml` に `$ref` で参照を追加
2. **バリデーション**: `make api-validate`
3. **コード生成**: `make api-generate` (自動的にバンドル → 生成)
4. **パッケージ作成**: `internal/features/<feature_name>/`
5. **実装**: Domain → Application → Infrastructure → Presentation
6. **SQL 作成**: `db/queries/` に追加
7. **sqlc 生成**: `make sqlc-generate`
8. **テスト作成**: `*_test.go`
9. **検証**: `make test` → `make lint` → `make build`

### OpenAPI 仕様の構成

OpenAPI 仕様はモジュール化されており、以下のように分割されています：

- **`api/openapi.yaml`**: メインファイル(各ファイルへの参照のみ)
- **`api/paths/`**: エンドポイントごとの定義
- **`api/components/schemas/`**: データモデル定義(Sake, common等)
- **`api/components/responses/`**: 共通レスポンス定義
- **`api/components/parameters/`**: 共通パラメータ定義
- **`api/components/securitySchemes/`**: 認証スキーム定義

新しいエンドポイントを追加する際は、`paths/` に新しいファイルを作成し、
`openapi.yaml` から `$ref` で参照してください。

## 環境変数

| 変数名              | デフォルト値     | 説明                       |
| ------------------- | ---------------- | -------------------------- |
| `SERVER_PORT`       | 8080             | サーバーポート             |
| `GIN_MODE`          | debug            | Gin モード (debug/release) |
| `POSTGRES_HOST`     | localhost        | PostgreSQL ホスト          |
| `POSTGRES_DB`       | sake_hack_app    | データベース名             |
| `POSTGRES_USER`     | postgres         | ユーザー名                 |
| `POSTGRES_PASSWORD` | sakehacksakehack | パスワード                 |
| `VALKEY_HOST`       | localhost        | Valkey ホスト              |
| `VALKEY_PASSWORD`   | sakehacksakehack | Valkey パスワード          |

## コーディング規約

### 型変換

CWE-190 Integer Overflow 対策のため、必ず安全な型変換を使用:

```go
// ✅ 正しい
difficulty, err := utils.IntToInt32(int(req.Difficulty))
if err != nil {
    return nil, apperror.BadRequestError("難易度の値が不正です")
}

// ❌ 禁止
return int32(value)  // gosec G115/G109エラー
```

### エラーハンドリング

```go
// アプリケーションエラー
if err != nil {
    return apperror.BadRequestError("不正なリクエストです")
}

// データベースエラー
logger.LogDatabaseError(ctx, "CREATE", "users", err)
return apperror.DatabaseError("ユーザー作成に失敗しました")
```

### ロギング

```go
// メソッドトレース
defer logger.TraceMethodAuto(ctx, params)()

// エラーログ
logger.LogDatabaseError(ctx, "SELECT", "users", err)
logger.LogBusinessError(ctx, "DuplicateEmail", err)
logger.LogValidationError(ctx, "email", email, "invalid format")
```

## Git コミット規約

フォーマット: `{type}:{emoji}{対象の説明}(#チケット番号)`

**Type**:

- `add` (新機能)
- `fix` (バグ修正)
- `update` (改善)
- `refactor` (リファクタリング)
- `docs` (ドキュメント)
- `test` (テスト)
- `chore` (雑務)

**Emoji**:

- ✨ (新機能)
- 🐛 (バグ修正)
- ⚡ (改善)
- ♻️ (リファクタリング)
- 📝 (ドキュメント)
- 🧪 (テスト)
- 🔧 (設定)

例: `add:✨ヘルスチェックエンドポイント実装(#123)`

## ライセンス

MIT License
