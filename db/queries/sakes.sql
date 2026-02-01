-- name: ListSakes :many
SELECT
    id,
    name,
    phonetic,
    image_id,
    category,
    description,
    alcohol_percentage,
    volume_max,
    volume_remain,
    region,
    price,
    memo,
    created_at,
    updated_at
FROM sakes
WHERE
    (sqlc.narg('type_id')::INTEGER IS NULL OR s.type_id = sqlc.narg('type_id'))
    AND (sqlc.narg('brewery_id')::INTEGER IS NULL OR s.brewery_id = sqlc.narg('brewery_id'))
ORDER BY s.created_at DESC
LIMIT $1 OFFSET $2;

-- name: CountSakes :one
SELECT COUNT(*) AS total
FROM sakes s
WHERE
    (sqlc.narg('type_id')::INTEGER IS NULL OR s.type_id = sqlc.narg('type_id'))
    AND (sqlc.narg('brewery_id')::INTEGER IS NULL OR s.brewery_id = sqlc.narg('brewery_id'));

-- name: GetSakeByID :one
SELECT
    id,
    name,
    phonetic,
    image_id,
    category,
    description,
    alcohol_percentage,
    volume_max,
    volume_remain,
    region,
    price,
    memo,
    created_at,
    updated_at
FROM sakes
WHERE id = $1;

-- name: CreateSake :one
INSERT INTO sakes (name, phonetic, image_id, category, description, alcohol_percentage, volume_max, volume_remain, region, price, memo)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
RETURNING id, created_at, updated_at;

-- name: UpdateSake :exec
UPDATE sakes
SET name = $2, phonetic = $3, image_id = $4, category = $5, description = $6, alcohol_percentage = $7, volume_max = $8, volume_remain = $9, region = $10, price = $11, memo = $12, updated_at = CURRENT_TIMESTAMP
WHERE id = $1;

-- name: GetDrinkStylesBySakeID :many
SELECT
    ds.id,
    ds.name,
    ds.description,
    ds.created_at,
    ds.updated_at
FROM drink_styles ds
INNER JOIN sake_drink_styles sds ON ds.id = sds.drink_style_id
WHERE sds.sake_id = $1
ORDER BY ds.id;
