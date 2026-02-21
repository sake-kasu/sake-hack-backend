#!/usr/bin/env bash

set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT_DIR"

APP_ENV_FILE=".env"
APP_ENV_SAMPLE=".env.sample"
DOCKER_ENV_FILE="docker/.env"
DOCKER_ENV_SAMPLE="docker/.env.sample"
DOCKER_COMPOSE_FILE="docker/compose.yaml"
TOOL_VERSIONS_FILE=".tool-versions"

require_cmd() {
  local cmd="$1"
  if ! command -v "$cmd" >/dev/null 2>&1; then
    echo "❌ Error: '$cmd' が見つかりません。インストールしてください。" >&2
    exit 1
  fi
}

copy_if_missing() {
  local src="$1"
  local dst="$2"
  if [[ -f "$dst" ]]; then
    echo "⏭️  skip: $dst は既に存在します"
    return
  fi
  if [[ ! -f "$src" ]]; then
    echo "❌ Error: $src が存在しません" >&2
    exit 1
  fi
  cp "$src" "$dst"
  echo "✅ created: $dst"
}

load_env_var() {
  local file="$1"
  local key="$2"
  grep -E "^${key}=" "$file" | tail -n 1 | cut -d '=' -f2- || true
}

ensure_asdf_tools() {
  local file="$1"
  if [[ ! -f "$file" ]]; then
    echo "❌ Error: $file が存在しません" >&2
    exit 1
  fi

  while read -r tool version _; do
    if [[ -z "${tool:-}" || "${tool:0:1}" == "#" ]]; then
      continue
    fi
    if [[ -z "${version:-}" ]]; then
      echo "❌ Error: $file の '${tool}' にバージョン指定がありません" >&2
      exit 1
    fi

    if ! asdf plugin list | tr -d ' ' | grep -Fx "$tool" >/dev/null 2>&1; then
      echo "🔌 add asdf plugin: $tool"
      asdf plugin add "$tool"
    fi

    if ! asdf where "$tool" "$version" >/dev/null 2>&1; then
      echo "📦 install tool: $tool $version"
      asdf install "$tool" "$version"
    else
      echo "⏭️  skip: tool '$tool $version' はインストール済みです"
    fi
  done < "$file"
}

run_aws() {
  if command -v asdf >/dev/null 2>&1; then
    asdf exec aws "$@"
  else
    aws "$@"
  fi
}

wait_for_rustfs() {
  local endpoint="$1"
  local access_key="$2"
  local secret_key="$3"
  local max_retries=30
  local i=1

  while (( i <= max_retries )); do
    if AWS_ACCESS_KEY_ID="$access_key" AWS_SECRET_ACCESS_KEY="$secret_key" \
      run_aws --endpoint-url "$endpoint" s3api list-buckets >/dev/null 2>&1; then
      return 0
    fi
    sleep 1
    i=$((i + 1))
  done

  echo "❌ Error: RustFS($endpoint) への接続を確認できませんでした" >&2
  return 1
}

require_cmd docker
require_cmd make
require_cmd asdf

echo "🔎 ensure asdf plugins/tools..."
ensure_asdf_tools "$TOOL_VERSIONS_FILE"
asdf reshim || true

if ! asdf exec aws --version >/dev/null 2>&1; then
  echo "❌ Error: aws CLI のセットアップに失敗しました" >&2
  exit 1
fi

copy_if_missing "$APP_ENV_SAMPLE" "$APP_ENV_FILE"
copy_if_missing "$DOCKER_ENV_SAMPLE" "$DOCKER_ENV_FILE"

echo "🐳 docker compose up..."
docker compose -f "$DOCKER_COMPOSE_FILE" up -d

if [[ ! -f "$APP_ENV_FILE" ]]; then
  echo "❌ Error: $APP_ENV_FILE が存在しません" >&2
  exit 1
fi

bucket_name="$(load_env_var "$APP_ENV_FILE" "STORAGE_BUCKET_NAME")"
access_key="$(load_env_var "$APP_ENV_FILE" "STORAGE_ACCESS_KEY")"
secret_key="$(load_env_var "$APP_ENV_FILE" "STORAGE_SECRET_KEY")"
endpoint="$(load_env_var "$APP_ENV_FILE" "STORAGE_ENDPOINT")"

bucket_name="${bucket_name:-sake-hack-bucket}"
endpoint="${endpoint:-http://localhost:9000}"
access_key="${access_key:-rustfsadmin}"
secret_key="${secret_key:-rustfsadmin}"

echo "⏳ wait for RustFS..."
wait_for_rustfs "$endpoint" "$access_key" "$secret_key"

echo "🪣 ensure bucket: $bucket_name"
if AWS_ACCESS_KEY_ID="$access_key" AWS_SECRET_ACCESS_KEY="$secret_key" \
  run_aws --endpoint-url "$endpoint" s3api head-bucket --bucket "$bucket_name" >/dev/null 2>&1; then
  echo "⏭️  skip: bucket '$bucket_name' は既に存在します"
else
  AWS_ACCESS_KEY_ID="$access_key" AWS_SECRET_ACCESS_KEY="$secret_key" \
    run_aws --endpoint-url "$endpoint" s3api create-bucket --bucket "$bucket_name" >/dev/null
  echo "✅ created: bucket '$bucket_name'"
fi

if ! command -v migrate >/dev/null 2>&1; then
  echo "📥 migrate コマンドが見つからないためインストールします..."
  make migrate-install
fi

echo "🗄️  run migrations..."
make migrate-up

echo "🎉 completed: 環境構築が完了しました"
echo "🚀 next: make run"
