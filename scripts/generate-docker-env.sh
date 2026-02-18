#!/bin/bash
# config.yamlからdocker/.env.localを自動生成するスクリプト

set -e

CONFIG_FILE="config/config.yml"
OUTPUT_FILE="docker/.env.local"

if [ ! -f "$CONFIG_FILE" ]; then
    echo "Error: $CONFIG_FILE not found"
    exit 1
fi

echo "# Auto-generated from config/config.yml" > "$OUTPUT_FILE"
echo "# 変更する場合はconfig/config.ymlを編集してください" >> "$OUTPUT_FILE"
echo "" >> "$OUTPUT_FILE"

# PostgreSQL settings
POSTGRES_DB=$(awk '/^database:/{flag=1} flag && /^  database:/{print $2; flag=0}' "$CONFIG_FILE")
POSTGRES_USER=$(awk '/^database:/{flag=1} flag && /^  user:/{print $2; flag=0}' "$CONFIG_FILE")
POSTGRES_PASSWORD=$(awk '/^database:/{flag=1} flag && /^  password:/{print $2; flag=0}' "$CONFIG_FILE")
POSTGRES_HOST_PORT=$(awk '/^database:/{flag=1} flag && /^  port:/{print $2; flag=0}' "$CONFIG_FILE")

echo "POSTGRES_DB=$POSTGRES_DB" >> "$OUTPUT_FILE"
echo "POSTGRES_USER=$POSTGRES_USER" >> "$OUTPUT_FILE"
echo "POSTGRES_PASSWORD=$POSTGRES_PASSWORD" >> "$OUTPUT_FILE"
echo "POSTGRES_HOST_PORT=$POSTGRES_HOST_PORT" >> "$OUTPUT_FILE"
echo "TZ=UTC" >> "$OUTPUT_FILE"
echo "" >> "$OUTPUT_FILE"

# Valkey settings
VALKEY_PASSWORD=$(awk '/^valkey:/{flag=1} flag && /^  password:/{print $2; flag=0}' "$CONFIG_FILE")
VALKEY_HOST_PORT=$(awk '/^valkey:/{flag=1} flag && /^  port:/{print $2; flag=0}' "$CONFIG_FILE")

echo "VALKEY_PASSWORD=$VALKEY_PASSWORD" >> "$OUTPUT_FILE"
echo "VALKEY_HOST_PORT=$VALKEY_HOST_PORT" >> "$OUTPUT_FILE"
