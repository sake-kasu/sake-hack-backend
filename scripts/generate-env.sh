#!/bin/bash
# config.yamlから.envを自動生成するスクリプト

set -e

CONFIG_FILE="config/config.yml"
OUTPUT_FILE=".env"

if [ ! -f "$CONFIG_FILE" ]; then
    echo "Error: $CONFIG_FILE not found"
    exit 1
fi

echo "# Auto-generated from config/config.yml - DO NOT EDIT MANUALLY" > "$OUTPUT_FILE"
echo "# Run 'make env' to regenerate this file" >> "$OUTPUT_FILE"
echo "" >> "$OUTPUT_FILE"

# =============================================================================
# Server設定
# =============================================================================
echo "# =============================================================================" >> "$OUTPUT_FILE"
echo "# Server設定" >> "$OUTPUT_FILE"
echo "# =============================================================================" >> "$OUTPUT_FILE"

SERVER_PORT=$(awk '/^server:/{flag=1} flag && /^  port:/{print $2; flag=0}' "$CONFIG_FILE")
SERVER_MODE=$(awk '/^server:/{flag=1} flag && /^  mode:/{print $2; flag=0}' "$CONFIG_FILE")
SERVER_GRACEFUL_SHUTDOWN_TIMEOUT=$(awk '/^server:/{flag=1} flag && /^  gracefulShutdownTimeout:/{print $2; flag=0}' "$CONFIG_FILE")

echo "SERVER_PORT=${SERVER_PORT:-8080}" >> "$OUTPUT_FILE"
echo "SERVER_MODE=${SERVER_MODE:-debug}" >> "$OUTPUT_FILE"
echo "SERVER_GRACEFUL_SHUTDOWN_TIMEOUT=${SERVER_GRACEFUL_SHUTDOWN_TIMEOUT:-30s}" >> "$OUTPUT_FILE"
echo "" >> "$OUTPUT_FILE"

# =============================================================================
# Database設定 (PostgreSQL)
# =============================================================================
echo "# =============================================================================" >> "$OUTPUT_FILE"
echo "# Database設定 (PostgreSQL)" >> "$OUTPUT_FILE"
echo "# =============================================================================" >> "$OUTPUT_FILE"

DB_HOST=$(awk '/^database:/{flag=1} flag && /^  host:/{print $2; flag=0}' "$CONFIG_FILE")
DB_PORT=$(awk '/^database:/{flag=1} flag && /^  port:/{print $2; flag=0}' "$CONFIG_FILE")
DB_NAME=$(awk '/^database:/{flag=1} flag && /^  database:/{print $2; flag=0}' "$CONFIG_FILE")
DB_USER=$(awk '/^database:/{flag=1} flag && /^  user:/{print $2; flag=0}' "$CONFIG_FILE")
DB_PASSWORD=$(awk '/^database:/{flag=1} flag && /^  password:/{print $2; flag=0}' "$CONFIG_FILE")
DB_SSL_MODE=$(awk '/^database:/{flag=1} flag && /^  sslmode:/{print $2; flag=0}' "$CONFIG_FILE")
DB_MAX_OPEN_CONNS=$(awk '/^database:/{flag=1} flag && /^  maxOpenConns:/{print $2; flag=0}' "$CONFIG_FILE")
DB_MAX_IDLE_CONNS=$(awk '/^database:/{flag=1} flag && /^  maxIdleConns:/{print $2; flag=0}' "$CONFIG_FILE")
DB_CONN_MAX_LIFETIME=$(awk '/^database:/{flag=1} flag && /^  connMaxLifetime:/{print $2; flag=0}' "$CONFIG_FILE")

echo "DB_HOST=${DB_HOST:-localhost}" >> "$OUTPUT_FILE"
echo "DB_PORT=${DB_PORT:-5432}" >> "$OUTPUT_FILE"
echo "DB_NAME=${DB_NAME:-sake_hack_app}" >> "$OUTPUT_FILE"
echo "DB_USER=${DB_USER:-postgres}" >> "$OUTPUT_FILE"
echo "DB_PASSWORD=${DB_PASSWORD:-sakehacksakehack}" >> "$OUTPUT_FILE"
echo "DB_SSL_MODE=${DB_SSL_MODE:-disable}" >> "$OUTPUT_FILE"
echo "DB_MAX_OPEN_CONNS=${DB_MAX_OPEN_CONNS:-25}" >> "$OUTPUT_FILE"
echo "DB_MAX_IDLE_CONNS=${DB_MAX_IDLE_CONNS:-5}" >> "$OUTPUT_FILE"
echo "DB_CONN_MAX_LIFETIME=${DB_CONN_MAX_LIFETIME:-5m}" >> "$OUTPUT_FILE"
echo "" >> "$OUTPUT_FILE"

# =============================================================================
# Cache設定 (Valkey)
# =============================================================================
echo "# =============================================================================" >> "$OUTPUT_FILE"
echo "# Cache設定 (Valkey)" >> "$OUTPUT_FILE"
echo "# =============================================================================" >> "$OUTPUT_FILE"

CACHE_HOST=$(awk '/^valkey:/{flag=1} flag && /^  host:/{print $2; flag=0}' "$CONFIG_FILE")
CACHE_PORT=$(awk '/^valkey:/{flag=1} flag && /^  port:/{print $2; flag=0}' "$CONFIG_FILE")
CACHE_PASSWORD=$(awk '/^valkey:/{flag=1} flag && /^  password:/{print $2; flag=0}' "$CONFIG_FILE")
CACHE_DATABASE=$(awk '/^valkey:/{flag=1} flag && /^  database:/{print $2; flag=0}' "$CONFIG_FILE")
CACHE_POOL_SIZE=$(awk '/^valkey:/{flag=1} flag && /^  poolSize:/{print $2; flag=0}' "$CONFIG_FILE")
CACHE_MIN_IDLE_CONNS=$(awk '/^valkey:/{flag=1} flag && /^  minIdleConns:/{print $2; flag=0}' "$CONFIG_FILE")
CACHE_MAX_RETRIES=$(awk '/^valkey:/{flag=1} flag && /^  maxRetries:/{print $2; flag=0}' "$CONFIG_FILE")
CACHE_DIAL_TIMEOUT=$(awk '/^valkey:/{flag=1} flag && /^  dialTimeout:/{print $2; flag=0}' "$CONFIG_FILE")
CACHE_READ_TIMEOUT=$(awk '/^valkey:/{flag=1} flag && /^  readTimeout:/{print $2; flag=0}' "$CONFIG_FILE")
CACHE_WRITE_TIMEOUT=$(awk '/^valkey:/{flag=1} flag && /^  writeTimeout:/{print $2; flag=0}' "$CONFIG_FILE")

echo "CACHE_HOST=${CACHE_HOST:-localhost}" >> "$OUTPUT_FILE"
echo "CACHE_PORT=${CACHE_PORT:-6379}" >> "$OUTPUT_FILE"
echo "CACHE_PASSWORD=${CACHE_PASSWORD:-sakehacksakehack}" >> "$OUTPUT_FILE"
echo "CACHE_DATABASE=${CACHE_DATABASE:-0}" >> "$OUTPUT_FILE"
echo "CACHE_POOL_SIZE=${CACHE_POOL_SIZE:-10}" >> "$OUTPUT_FILE"
echo "CACHE_MIN_IDLE_CONNS=${CACHE_MIN_IDLE_CONNS:-5}" >> "$OUTPUT_FILE"
echo "CACHE_MAX_RETRIES=${CACHE_MAX_RETRIES:-3}" >> "$OUTPUT_FILE"
echo "CACHE_DIAL_TIMEOUT=${CACHE_DIAL_TIMEOUT:-5s}" >> "$OUTPUT_FILE"
echo "CACHE_READ_TIMEOUT=${CACHE_READ_TIMEOUT:-3s}" >> "$OUTPUT_FILE"
echo "CACHE_WRITE_TIMEOUT=${CACHE_WRITE_TIMEOUT:-3s}" >> "$OUTPUT_FILE"
echo "" >> "$OUTPUT_FILE"

# =============================================================================
# JWT設定
# =============================================================================
echo "# =============================================================================" >> "$OUTPUT_FILE"
echo "# JWT設定" >> "$OUTPUT_FILE"
echo "# =============================================================================" >> "$OUTPUT_FILE"

JWT_SECRET=$(awk '/^jwt:/{flag=1} flag && /^  secret:/{print $2; flag=0}' "$CONFIG_FILE" | tr -d '"')
JWT_EXPIRATION=$(awk '/^jwt:/{flag=1} flag && /^  expiration:/{print $2; flag=0}' "$CONFIG_FILE")
JWT_COOKIE_SECURE=$(awk '/^jwt:/{flag=1} flag && /^  cookieSecure:/{print $2; flag=0}' "$CONFIG_FILE")
JWT_COOKIE_NAME=$(awk '/^jwt:/{flag=1} flag && /^  cookieName:/{print $2; flag=0}' "$CONFIG_FILE" | tr -d '"')
JWT_COOKIE_PATH=$(awk '/^jwt:/{flag=1} flag && /^  cookiePath:/{print $2; flag=0}' "$CONFIG_FILE" | tr -d '"')
JWT_COOKIE_DOMAIN=$(awk '/^jwt:/{flag=1} flag && /^  cookieDomain:/{print $2; flag=0}' "$CONFIG_FILE" | tr -d '"')

echo "JWT_SECRET=${JWT_SECRET:-your-secret-key-change-me-in-production}" >> "$OUTPUT_FILE"
echo "JWT_EXPIRATION=${JWT_EXPIRATION:-86400}" >> "$OUTPUT_FILE"
echo "JWT_COOKIE_SECURE=${JWT_COOKIE_SECURE:-false}" >> "$OUTPUT_FILE"
echo "JWT_COOKIE_NAME=${JWT_COOKIE_NAME:-sake_hack_token}" >> "$OUTPUT_FILE"
echo "JWT_COOKIE_PATH=${JWT_COOKIE_PATH:-/}" >> "$OUTPUT_FILE"
echo "JWT_COOKIE_DOMAIN=${JWT_COOKIE_DOMAIN:-}" >> "$OUTPUT_FILE"
echo "" >> "$OUTPUT_FILE"

# =============================================================================
# CORS設定
# =============================================================================
echo "# =============================================================================" >> "$OUTPUT_FILE"
echo "# CORS設定" >> "$OUTPUT_FILE"
echo "# =============================================================================" >> "$OUTPUT_FILE"

# CORS設定は配列なので特別な処理が必要
CORS_ALLOWED_ORIGINS=$(awk '/^cors:/{flag=1} flag && /^  allowedOrigins:/{flag=2; next} flag==2 && /^  [a-zA-Z]/{flag=0} flag==2 && /^ *-/{gsub(/^ *- *"?|"?$/, ""); printf "%s,", $0}' "$CONFIG_FILE" | sed 's/,$//')
CORS_ALLOWED_METHODS=$(awk '/^cors:/{flag=1} flag && /^  allowedMethods:/{flag=2; next} flag==2 && /^  [a-zA-Z]/{flag=0} flag==2 && /^ *-/{gsub(/^ *- *"?|"?$/, ""); printf "%s,", $0}' "$CONFIG_FILE" | sed 's/,$//')
CORS_ALLOWED_HEADERS=$(awk '/^cors:/{flag=1} flag && /^  allowedHeaders:/{flag=2; next} flag==2 && /^  [a-zA-Z]/{flag=0} flag==2 && /^ *-/{gsub(/^ *- *"?|"?$/, ""); printf "%s,", $0}' "$CONFIG_FILE" | sed 's/,$//')
CORS_EXPOSED_HEADERS=$(awk '/^cors:/{flag=1} flag && /^  exposedHeaders:/{flag=2; next} flag==2 && /^  [a-zA-Z]/{flag=0} flag==2 && /^ *-/{gsub(/^ *- *"?|"?$/, ""); printf "%s,", $0}' "$CONFIG_FILE" | sed 's/,$//')
CORS_ALLOW_CREDENTIALS=$(awk '/^cors:/{flag=1} flag && /^  allowCredentials:/{print $2; flag=0}' "$CONFIG_FILE")
CORS_MAX_AGE=$(awk '/^cors:/{flag=1} flag && /^  maxAge:/{print $2; flag=0}' "$CONFIG_FILE")

echo "CORS_ALLOWED_ORIGINS=${CORS_ALLOWED_ORIGINS:-http://localhost:3000,http://localhost:8080}" >> "$OUTPUT_FILE"
echo "CORS_ALLOWED_METHODS=${CORS_ALLOWED_METHODS:-GET,POST,PUT,DELETE,PATCH,OPTIONS}" >> "$OUTPUT_FILE"
echo "CORS_ALLOWED_HEADERS=${CORS_ALLOWED_HEADERS:-Origin,Content-Type,Accept,Authorization}" >> "$OUTPUT_FILE"
echo "CORS_EXPOSED_HEADERS=${CORS_EXPOSED_HEADERS:-Content-Length}" >> "$OUTPUT_FILE"
echo "CORS_ALLOW_CREDENTIALS=${CORS_ALLOW_CREDENTIALS:-true}" >> "$OUTPUT_FILE"
echo "CORS_MAX_AGE=${CORS_MAX_AGE:-43200}" >> "$OUTPUT_FILE"
echo "" >> "$OUTPUT_FILE"

# =============================================================================
# Logging設定
# =============================================================================
echo "# =============================================================================" >> "$OUTPUT_FILE"
echo "# Logging設定" >> "$OUTPUT_FILE"
echo "# =============================================================================" >> "$OUTPUT_FILE"

LOG_LEVEL=$(awk '/^logging:/{flag=1} flag && /^  level:/{print $2; flag=0}' "$CONFIG_FILE")
LOG_FORMAT=$(awk '/^logging:/{flag=1} flag && /^  format:/{print $2; flag=0}' "$CONFIG_FILE")

echo "LOG_LEVEL=${LOG_LEVEL:-debug}" >> "$OUTPUT_FILE"
echo "LOG_FORMAT=${LOG_FORMAT:-console}" >> "$OUTPUT_FILE"
