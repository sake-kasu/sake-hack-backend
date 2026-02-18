-- name: ListSakes :many
SELECT
    s.id,
    s.category,
    s.name,
    s.image_url
FROM sakes s
WHERE
    (sqlc.narg('kind_id')::INTEGER IS NULL OR s.kind_id = sqlc.narg('kind_id'))
    AND (sqlc.narg('brewery_id')::INTEGER IS NULL OR s.brewery_id = sqlc.narg('brewery_id'))
ORDER BY s.created_at DESC
LIMIT $1 OFFSET $2;

-- name: CountSakes :one
SELECT COUNT(*) AS total
FROM sakes s
WHERE
    (sqlc.narg('kind_id')::INTEGER IS NULL OR s.kind_id = sqlc.narg('kind_id'))
    AND (sqlc.narg('brewery_id')::INTEGER IS NULL OR s.brewery_id = sqlc.narg('brewery_id'));

-- name: ListPublicSakes :many
SELECT
    s.id,
    s.category,
    s.name,
    s.image_url,
    s.abv,
    s.memo,
    s.created_at,
    s.updated_at,
    sk.name AS kind_name,
    b.origin_region AS brewery_region
FROM sakes s
INNER JOIN sake_kinds sk ON s.kind_id = sk.id
INNER JOIN breweries b ON s.brewery_id = b.id
WHERE
    (sqlc.narg('category')::sake_category IS NULL OR s.category = sqlc.narg('category'))
    AND (sqlc.narg('search_query')::TEXT IS NULL OR s.name ILIKE '%' || sqlc.narg('search_query') || '%')
ORDER BY s.created_at DESC
LIMIT $1 OFFSET $2;

-- name: CountPublicSakes :one
SELECT COUNT(*) AS total
FROM sakes s
WHERE
    (sqlc.narg('category')::sake_category IS NULL OR s.category = sqlc.narg('category'))
    AND (sqlc.narg('search_query')::TEXT IS NULL OR s.name ILIKE '%' || sqlc.narg('search_query') || '%');

-- name: GetSakeDetailByID :one
SELECT
    s.id,
    s.category,
    s.name,
    s.phonetic,
    s.abv,
    s.purchase_volume,
    s.remaining_volume,
    s.memo,
    s.price,
    s.image_url,
    s.created_at,
    s.updated_at,
    sk.id AS kind_id,
    sk.name AS kind_name,
    b.id AS brewery_id,
    b.name AS brewery_name,
    b.origin_country AS brewery_origin_country,
    b.origin_region AS brewery_origin_region,
    b.latitude AS brewery_latitude,
    b.longitude AS brewery_longitude
FROM sakes s
INNER JOIN sake_kinds sk ON s.kind_id = sk.id
INNER JOIN breweries b ON s.brewery_id = b.id
WHERE s.id = $1;

-- name: CreateSake :one
INSERT INTO sakes (category, kind_id, brewery_id, name, phonetic, abv, purchase_volume, remaining_volume, memo, price, image_url)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
RETURNING id, category, name, image_url, created_at, updated_at;

-- name: UpdateSake :one
UPDATE sakes
SET category = $2,
    kind_id = $3,
    brewery_id = $4,
    name = $5,
    phonetic = $6,
    abv = $7,
    purchase_volume = $8,
    remaining_volume = $9,
    memo = $10,
    price = $11,
    image_url = $12,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING id, category, name, image_url, created_at, updated_at;

-- name: DeleteSake :execrows
DELETE FROM sakes WHERE id = $1;

-- name: GetDrinkStylesBySakeID :many
SELECT
    ds.id,
    ds.name,
    ds.description
FROM drink_styles ds
INNER JOIN sake_drink_styles sds ON ds.id = sds.drink_style_id
WHERE sds.sake_id = $1
ORDER BY ds.id;

-- name: UpsertBrewery :one
INSERT INTO breweries (name, origin_country, origin_region, latitude, longitude)
VALUES ($1, $2, $3, $4, $5)
ON CONFLICT (name, origin_country) DO UPDATE SET
    origin_region = EXCLUDED.origin_region,
    latitude = EXCLUDED.latitude,
    longitude = EXCLUDED.longitude,
    updated_at = CURRENT_TIMESTAMP
RETURNING id;

-- name: GetSakeKindByName :one
SELECT id, name FROM sake_kinds WHERE name = $1;

-- name: UpsertSakeKind :one
INSERT INTO sake_kinds (name)
VALUES ($1)
ON CONFLICT (name) DO UPDATE SET
    updated_at = CURRENT_TIMESTAMP
RETURNING id;

-- name: UpsertDrinkStyle :one
INSERT INTO drink_styles (name, description)
VALUES ($1, $2)
ON CONFLICT (name) DO UPDATE SET
    description = EXCLUDED.description,
    updated_at = CURRENT_TIMESTAMP
RETURNING id;

-- name: InsertSakeDrinkStyle :exec
INSERT INTO sake_drink_styles (sake_id, drink_style_id)
VALUES ($1, $2)
ON CONFLICT DO NOTHING;

-- name: DeleteSakeDrinkStyles :exec
DELETE FROM sake_drink_styles WHERE sake_id = $1;
