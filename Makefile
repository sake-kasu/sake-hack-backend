.PHONY: help build run air-install dev clean test test-unit test-integration cover lint deps submodule-init submodule-update submodule-status openapi-generate openapi-gendoc openapi-watch migrate-install migrate-create migrate-up migrate-up-one migrate-down migrate-reset migrate-force migrate-version migrate-status sqlc-generate

# .env fileが存在すれば読み込み
-include .env
export

# 変数定義
BINARY_NAME=sake-hack-server
MIGRATE_VERSION=v4.18.1
MAIN_PATH=./cmd/server
BUILD_DIR=./bin
OAPI_CODEGEN_VERSION=v2.6.0

# DB_URLを環境変数から動的生成
DB_HOST ?= localhost
DB_PORT ?= 5432
DB_USER ?= postgres
DB_PASSWORD ?= sakehacksakehack
DB_NAME ?= sake_hack_app
DB_SSL_MODE ?= disable
DB_URL=postgresql://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSL_MODE)

help: ## このヘルプメッセージを表示
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-20s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# ビルド・実行
build: ## アプリケーションをビルド
	@echo "🔨 $(BINARY_NAME)をビルドしています..."
	@go build -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)

run: ## アプリケーションを実行
	@echo "🚀 $(BINARY_NAME)を実行しています..."
	@go run $(MAIN_PATH)/main.go

air-install: ## airのインストール
	@echo "Installing air..."
	@go install github.com/air-verse/air@latest

dev: air-install ## ホットリロードで開発サーバーを起動(Air使用)
	@echo "🔥 開発サーバーを起動しています(ホットリロード有効)..."
	@air

clean: ## ビルド成果物をクリーンアップ
	@echo "🧹 クリーンアップ中..."
	@rm -rf $(BUILD_DIR)
	@rm -f coverage.out coverage.html

# テスト・品質
test: ## 全てのテストを実行
	@echo "🧪 テストを実行しています..."
	@go test -v -race ./...

test-unit: ## 単体テストのみ実行
	@echo "⚡ 単体テストを実行しています..."
	@go test -v -short ./...

test-integration: ## 統合テストのみ実行
	@echo "🔗 統合テストを実行しています..."
	@go test -v -run Integration ./...

cover: ## カバレッジ測定付きでテストを実行(自動生成コード除外)
	@echo "📊 カバレッジを計測しています..."
	@go test -v -race -coverprofile=coverage.out -covermode=atomic $$(go list ./... | grep -v "/generated$$")
	@echo ""
	@echo "📈 全体カバレッジ:"
	@go tool cover -func=coverage.out | grep total | awk '{print "   Total Coverage: " $$3}'
	@echo ""
	@go tool cover -html=coverage.out -o coverage.html
	@echo "✅ カバレッジレポートを生成しました: coverage.html"

lint: ## リンターを実行
	@echo "🔍 リンターを実行しています..."
	@golangci-lint run --timeout=5m ./...

arch-check: ## アーキテクチャ検証(Feature間依存チェック)
	@echo "🏗️ Feature間依存チェックを実行中..."
	@FEATURES=$$(ls -d internal/features/*/ 2>/dev/null | xargs -n1 basename | grep -v shared); \
	VIOLATIONS=0; \
	for feature in $$FEATURES; do \
		for other in $$FEATURES; do \
			if [ "$$feature" != "$$other" ]; then \
				FOUND=$$(grep -r "features/$$other" internal/features/$$feature --include="*.go" 2>/dev/null | grep -v "_test.go" || true); \
				if [ -n "$$FOUND" ]; then \
					echo ""; \
					echo "❌ $$feature -> $$other への直接依存を検出:"; \
					echo "$$FOUND" | head -5; \
					VIOLATIONS=$$((VIOLATIONS + 1)); \
				fi; \
			fi; \
		done; \
	done; \
	if [ $$VIOLATIONS -gt 0 ]; then \
		echo ""; \
		echo "⚠️  Feature間の直接依存が検出されました。"; \
		echo "   Consumer側でインターフェースを定義し、Application層で型変換を行ってください。"; \
		echo "   (shared機能への依存は許可されます)"; \
		exit 1; \
	else \
		echo "✅ Feature間の直接依存はありません"; \
	fi

gosec-install: ## Gosecのインストール
	@echo "Installing gosec..."
	@go install github.com/securego/gosec/v2/cmd/gosec@latest

gosec-scan: gosec-install ## セキュリティスキャナーを実行
	@echo "🔍 Gosec セキュリティスキャンを実行中..."
	@rm -f gosec-report.json
	@gosec -fmt json -out gosec-report.json \
		-exclude-dir=.git \
		-exclude-dir=.go \
		-exclude-dir=vendor \
		-exclude-dir=generated \
		-exclude-generated \
		-tests=false \
		-concurrency=4 \
		-severity=high \
		./...; \
	GOSEC_EXIT_CODE=$$?; \
	if [ -f gosec-report.json ]; then \
		if command -v jq >/dev/null 2>&1; then \
			ISSUE_COUNT=$$(jq '.Stats.found // 0' gosec-report.json); \
		else \
			ISSUE_COUNT=$$(grep -o '"found": [0-9]*' gosec-report.json | grep -o '[0-9]*' || echo "0"); \
		fi; \
		if [ "$$ISSUE_COUNT" -gt 0 ]; then \
			echo ""; \
			echo "❌ セキュリティ上の問題が $$ISSUE_COUNT 件検出されました"; \
			echo ""; \
			echo "📋 検出された問題:"; \
			if command -v jq >/dev/null 2>&1; then \
				jq -r '.Issues[] | "  [\(.severity)] \(.file):\(.line) - \(.details)"' gosec-report.json; \
			else \
				cat gosec-report.json; \
			fi; \
			echo ""; \
			echo "📄 詳細レポート: gosec-report.json"; \
			exit 1; \
		else \
			echo "✅ セキュリティ上の問題は検出されませんでした"; \
		fi \
	else \
		echo "✅ セキュリティ上の問題は検出されませんでした"; \
		exit $$GOSEC_EXIT_CODE; \
	fi

# 依存関係
deps: ## 依存関係を整理
	@echo "📦 依存関係を整理しています..."
	@go mod tidy
	@go mod download

# サブモジュール
submodule-init: ## サブモジュールを初期化
	@echo "🔧 サブモジュールを初期化しています..."
	@git submodule update --init --recursive

submodule-update: ## サブモジュールを最新に更新
	@echo "🔄 サブモジュールを最新に更新しています..."
	@git submodule update --remote --merge
	@echo "✅ サブモジュールが最新になりました"

submodule-status: ## サブモジュールの状態を確認
	@echo "📋 サブモジュールの状態:"
	@git submodule status

# API開発(OpenAPI仕様から自動生成)
openapi-generate: ## OpenAPI仕様からコードを自動生成
	@echo "🤖 OpenAPI仕様からコードを生成しています..."
	@echo "📦 Step 1: OpenAPI仕様をバンドルしています..."
	@npx @redocly/cli bundle openapi/openapi.yaml -o openapi/openapi.bundled.yaml
	@echo "⚙️  Step 2: Goコードを生成しています..."
	@mkdir -p api/generated
	@go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@$(OAPI_CODEGEN_VERSION) \
		-config openapi/oapi-codegen.yaml openapi/openapi.bundled.yaml

openapi-gendoc: ## APIドキュメントを生成
	@echo "📚 APIドキュメントを生成しています..."
	@npx @redocly/cli build-docs openapi/openapi.yaml -o openapi/docs/index.html
openapi-watch: ## APIドキュメントを監視して自動更新
	@echo "👀 APIファイルを監視しています..."
	@echo "📝 変更を検知すると自動的にドキュメントを再生成します"
	@echo "🌐 ドキュメント: http://localhost:8080"
	@npx concurrently -n "watch,serve" -c "blue,green" \
		"npx nodemon --watch openapi --ext yaml,json --exec 'make openapi-gendoc'" \
		"cd openapi/docs && python3 -m http.server 8080"

migrate-install: ## golang-migrateのインストール
	@echo "golang-migrate をインストールしています..."
	@go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@$(MIGRATE_VERSION)

migrate-create: ## 新規マイグレーション作成 (NAME=xxx)
	@if [ -z "$(NAME)" ]; then echo "Error: NAME is required. Usage: make migrate-create NAME=xxx"; exit 1; fi
	@echo "✨ マイグレーションを作成しています: $(NAME)"
	@migrate create -ext sql -dir db/migrations -seq $(NAME)

migrate-up: ## 全マイグレーション適用
	@echo "⬆️  マイグレーションを実行しています(up)..."
	@migrate -path db/migrations -database "$(DB_URL)" up

migrate-up-one: ## 1つ次のマイグレーション適用
	@echo "⬆️  マイグレーションを実行しています(up-one)..."
	@migrate -path db/migrations -database "$(DB_URL)" up 1

migrate-down: ## 1つ前にロールバック
	@echo "⬇️  マイグレーションをロールバックしています(down)..."
	@migrate -path db/migrations -database "$(DB_URL)" down 1

migrate-reset: ## 全リセット(down -all -> up)
	@echo "♻️  マイグレーションを全リセットしています..."
	@migrate -path db/migrations -database "$(DB_URL)" down -all
	@migrate -path db/migrations -database "$(DB_URL)" up

migrate-force: ## バージョン強制設定 (VERSION=xxx) ※障害復旧用
	@if [ -z "$(VERSION)" ]; then echo "Error: VERSION is required. Usage: make migrate-force VERSION=xxx"; exit 1; fi
	@echo "Forcing version: $(VERSION)..."
	@migrate -path db/migrations -database "$(DB_URL)" force $(VERSION)

migrate-version: ## 現在のバージョン確認
	@migrate -path db/migrations -database "$(DB_URL)" version

migrate-status: ## マイグレーション状態確認
	@echo "Migration status:"
	@migrate -path db/migrations -database "$(DB_URL)" version 2>&1 || true

# sqlc
sqlc-generate: ## SQLからGoコードを生成
	@echo "🔧 SQLからGoコードを生成しています..."
	@sqlc generate

generate: openapi-generate sqlc-generate ## 全コード生成(OpenAPI + SQLC)
