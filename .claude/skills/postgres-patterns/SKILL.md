---
name: postgres-patterns
description: PostgreSQL + SQLC + PostGIS + H3インデックスのデータベースパターン。クエリ最適化、スキーマ設計、マイグレーション管理。
---

# PostgreSQLパターン

PostgreSQL + SQLC + PostGIS/go-geom + H3インデックスを使用したデータベースパターン集。

## 活用タイミング

- SQLクエリやマイグレーションの作成
- データベーススキーマの設計
- 遅いクエリのトラブルシューティング
- SQLC クエリファイルの作成
- H3/PostGISインデックスの設計

## インデックス チートシート

| クエリパターン | インデックス型 | 例 |
|--------------|------------|---------|
| `WHERE col = value` | B-tree(デフォルト) | `CREATE INDEX idx ON t (col)` |
| `WHERE col > value` | B-tree | `CREATE INDEX idx ON t (col)` |
| `WHERE a = x AND b > y` | 複合 | `CREATE INDEX idx ON t (a, b)` |
| `WHERE jsonb @> '{}'` | GIN | `CREATE INDEX idx ON t USING gin (col)` |
| `WHERE tsv @@ query` | GIN | `CREATE INDEX idx ON t USING gin (col)` |
| 時系列の範囲検索 | BRIN | `CREATE INDEX idx ON t USING brin (col)` |
| 空間検索(PostGIS) | GiST | `CREATE INDEX idx ON t USING gist (geom)` |

## データ型 クイックリファレンス

| 用途 | 正しい型 | 避けるべき型 |
|------|---------|-------------|
| ID | `uuid` (UUIDv7推奨) | `int`, ランダムUUID |
| 文字列 | `text` | `varchar(255)` |
| タイムスタンプ | `timestamptz` | `timestamp` |
| 金額 | `numeric(10,2)` | `float` |
| フラグ | `boolean` | `varchar`, `int` |
| H3インデックス | `varchar(15)` | `text`, `bigint` |
| ジオメトリ | `geometry(Polygon, 4326)` | `text`, `jsonb` |

## H3インデックスパターン

### テーブル設計

```sql
CREATE TABLE fields (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    name text NOT NULL,
    geometry geometry(Polygon, 4326) NOT NULL,
    -- H3インデックス: 4解像度をVARCHAR(15)で保存
    h3_index_res3 varchar(15),
    h3_index_res5 varchar(15),
    h3_index_res7 varchar(15),
    h3_index_res9 varchar(15),
    -- 監査カラム
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    created_by uuid NOT NULL,
    updated_by uuid NOT NULL
);

-- H3インデックスの検索用インデックス
CREATE INDEX fields_h3_res3_idx ON fields (h3_index_res3);
CREATE INDEX fields_h3_res5_idx ON fields (h3_index_res5);
CREATE INDEX fields_h3_res7_idx ON fields (h3_index_res7);
CREATE INDEX fields_h3_res9_idx ON fields (h3_index_res9);

-- 空間インデックス
CREATE INDEX fields_geometry_idx ON fields USING gist (geometry);
```

### Go側でのH3計算(uber/h3-go)

```go
import "github.com/uber/h3-go/v4"

// centroidからH3インデックスを計算
func calculateH3Indexes(lat, lng float64) H3Indexes {
    latLng := h3.NewLatLng(lat, lng)
    return H3Indexes{
        Res3: h3.LatLngToCell(latLng, 3).String(),
        Res5: h3.LatLngToCell(latLng, 5).String(),
        Res7: h3.LatLngToCell(latLng, 7).String(),
        Res9: h3.LatLngToCell(latLng, 9).String(),
    }
}
```

## PostGIS/go-geom パターン

### ジオメトリの保存と取得

```go
import (
    "github.com/twpayne/go-geom"
    "github.com/twpayne/go-geom/encoding/ewkb"
)

// Polygonの作成
func createPolygon(coords [][]float64) *geom.Polygon {
    flatCoords := make([]float64, 0, len(coords)*2)
    for _, c := range coords {
        flatCoords = append(flatCoords, c[0], c[1])
    }
    return geom.NewPolygon(geom.XY).MustSetCoords([][]geom.Coord{
        // 座標リスト(最初と最後が同じ点で閉じる)
    })
}
```

### wagriデータ変換

```go
// LinearPolygon → Polygon に変換
// 座標が閉じていない場合は最初の点を末尾に追加
func convertLinearRingToPolygon(coords [][]float64) [][]float64 {
    if len(coords) == 0 {
        return coords
    }
    first := coords[0]
    last := coords[len(coords)-1]
    if first[0] != last[0] || first[1] != last[1] {
        coords = append(coords, first)
    }
    return coords
}
```

## SQLC統合パターン

### クエリファイル設計(`db/queries/`)

```sql
-- db/queries/fields.sql

-- name: GetField :one
-- 指定されたIDの圃場を取得する
SELECT id, name, geometry, h3_index_res3, h3_index_res5,
       h3_index_res7, h3_index_res9,
       created_at, updated_at, created_by, updated_by
FROM fields
WHERE id = $1;

-- name: ListFieldsByH3 :many
-- H3インデックス(解像度5)で圃場を検索する
SELECT id, name, geometry, h3_index_res5
FROM fields
WHERE h3_index_res5 = $1
ORDER BY name;

-- name: CreateField :one
-- 新規圃場を作成する
INSERT INTO fields (name, geometry, h3_index_res3, h3_index_res5,
                    h3_index_res7, h3_index_res9, created_by, updated_by)
VALUES ($1, $2, $3, $4, $5, $6, $7, $7)
RETURNING *;

-- name: UpsertFieldBatch :exec
-- 圃場を一括でUpsertする
INSERT INTO fields (id, name, geometry, h3_index_res3, h3_index_res5,
                    h3_index_res7, h3_index_res9, created_by, updated_by)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $8)
ON CONFLICT (id) DO UPDATE SET
    name = EXCLUDED.name,
    geometry = EXCLUDED.geometry,
    h3_index_res3 = EXCLUDED.h3_index_res3,
    h3_index_res5 = EXCLUDED.h3_index_res5,
    h3_index_res7 = EXCLUDED.h3_index_res7,
    h3_index_res9 = EXCLUDED.h3_index_res9,
    updated_at = now(),
    updated_by = EXCLUDED.updated_by;
```

### SQLC生成と注意事項

```bash
# SQLC生成
make sqlc-generate

# 生成先: internal/database/sqlc/ (編集禁止)
# クエリファイル: db/queries/<テーブル名>.sql
```

**重要**: `internal/database/sqlc/` のコードは直接編集禁止。SQLを変更する場合は `db/queries/` を編集して再生成する。

## マイグレーション

### 命名規則

```
NNNNNN_動詞_対象.sql
```

例:
- `000001_create_fields.sql`
- `000002_add_h3_indexes_to_fields.sql`
- `000003_create_import_jobs.sql`

### 必須ルール

- **up.sql と down.sql は必ずペアで作成**
- **down.sql は完全にロールバック可能であること**
- CLI運用: サーバー起動時ではなく `make migrate-*` コマンドで管理

### マイグレーションコマンド

```bash
# 新規マイグレーション作成
make migrate-create NAME=create_fields

# 全マイグレーション適用
make migrate-up

# 1つ前にロールバック
make migrate-down

# 現在のバージョン確認
make migrate-version
```

### マイグレーション例

```sql
-- 000001_create_fields.up.sql
CREATE EXTENSION IF NOT EXISTS postgis;

CREATE TABLE fields (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    name text NOT NULL,
    geometry geometry(Polygon, 4326) NOT NULL,
    h3_index_res3 varchar(15),
    h3_index_res5 varchar(15),
    h3_index_res7 varchar(15),
    h3_index_res9 varchar(15),
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    created_by uuid NOT NULL,
    updated_by uuid NOT NULL
);

CREATE INDEX fields_geometry_idx ON fields USING gist (geometry);
CREATE INDEX fields_h3_res3_idx ON fields (h3_index_res3);
CREATE INDEX fields_h3_res5_idx ON fields (h3_index_res5);
CREATE INDEX fields_h3_res7_idx ON fields (h3_index_res7);
CREATE INDEX fields_h3_res9_idx ON fields (h3_index_res9);
CREATE INDEX fields_created_by_idx ON fields (created_by);

-- 000001_create_fields.down.sql
DROP TABLE IF EXISTS fields;
```

## 共通パターン

### 複合インデックスの順序

```sql
-- 等値カラムを先、範囲カラムを後
CREATE INDEX idx ON orders (status, created_at);
-- WHERE status = 'pending' AND created_at > '2024-01-01' で有効
```

### カバリングインデックス

```sql
CREATE INDEX idx ON users (email) INCLUDE (name, created_at);
-- テーブル参照なしでSELECT email, name, created_at が可能
```

### 部分インデックス

```sql
CREATE INDEX idx ON users (email) WHERE deleted_at IS NULL;
-- 小さなインデックス、アクティブユーザーのみ
```

### UPSERT

```sql
INSERT INTO settings (user_id, key, value)
VALUES (123, 'theme', 'dark')
ON CONFLICT (user_id, key)
DO UPDATE SET value = EXCLUDED.value, updated_at = now();
```

### カーソルベースページネーション

```sql
-- O(1)のパフォーマンス(OFFSETはO(n))
SELECT * FROM products WHERE id > $last_id ORDER BY id LIMIT 20;
```

### キュー処理(SKIP LOCKED)

```sql
UPDATE jobs SET status = 'processing'
WHERE id = (
  SELECT id FROM jobs WHERE status = 'pending'
  ORDER BY created_at LIMIT 1
  FOR UPDATE SKIP LOCKED
) RETURNING *;
```

### バッチInsert

```sql
-- 個別INSERTの10-50倍高速
INSERT INTO events (user_id, action) VALUES
  (1, 'click'),
  (2, 'view'),
  (3, 'click');
```

### N+1クエリの排除

```sql
-- 悪い: N+1パターン
SELECT id FROM users WHERE active = true;
-- 100回の個別クエリ...

-- 良い: ANYで一括取得
SELECT * FROM orders WHERE user_id = ANY(ARRAY[1, 2, 3]);

-- 良い: JOINで取得
SELECT u.id, u.name, o.*
FROM users u
LEFT JOIN orders o ON o.user_id = u.id
WHERE u.active = true;
```

## 監査カラムパターン

全テーブルに以下の監査カラムを含める:

```sql
created_at timestamptz NOT NULL DEFAULT now(),
updated_at timestamptz NOT NULL DEFAULT now(),
created_by uuid NOT NULL,
updated_by uuid NOT NULL
```

`updated_at` の自動更新トリガー:

```sql
CREATE OR REPLACE FUNCTION update_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = now();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER fields_updated_at
    BEFORE UPDATE ON fields
    FOR EACH ROW EXECUTE FUNCTION update_updated_at();
```

## アンチパターン検出

```sql
-- 未インデックスの外部キーを検出
SELECT conrelid::regclass, a.attname
FROM pg_constraint c
JOIN pg_attribute a ON a.attrelid = c.conrelid AND a.attnum = ANY(c.conkey)
WHERE c.contype = 'f'
  AND NOT EXISTS (
    SELECT 1 FROM pg_index i
    WHERE i.indrelid = c.conrelid AND a.attnum = ANY(i.indkey)
  );

-- 遅いクエリを検出
SELECT query, mean_exec_time, calls
FROM pg_stat_statements
WHERE mean_exec_time > 100
ORDER BY mean_exec_time DESC;

-- テーブル肥大化を検出
SELECT relname, n_dead_tup, last_vacuum
FROM pg_stat_user_tables
WHERE n_dead_tup > 1000
ORDER BY n_dead_tup DESC;
```

## 関連

- エージェント: `database-reviewer` — データベースの包括的レビュー
- SQLC生成: `make sqlc-generate`
- マイグレーション: `make migrate-create`, `make migrate-up`
