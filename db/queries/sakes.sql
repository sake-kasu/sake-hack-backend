-- name: ListSakes :many
SELECT
    s.id,
    s.category,
    s.name,
    si.image_key
FROM sakes s
LEFT JOIN sake_images si ON si.sake_id = s.id AND si.sort_order = 0
WHERE
    (sqlc.narg('category')::sake_category IS NULL OR s.category = sqlc.narg('category'))
ORDER BY s.created_at DESC
LIMIT $1 OFFSET $2;

-- name: CountSakes :one
SELECT COUNT(*) AS total
FROM sakes s
WHERE
    (sqlc.narg('category')::sake_category IS NULL OR s.category = sqlc.narg('category'));

-- name: GetSakeDetailByID :one
SELECT
    s.id,
    s.category,
    s.name,
    s.phonetic,
    s.alcohol_percentage,
    s.volume_max,
    s.volume_remain,
    s.region,
    s.memo,
    s.price,
    s.created_at,
    s.updated_at
FROM sakes s
WHERE s.id = $1;

-- name: CreateSake :one
INSERT INTO sakes (category, name, phonetic, alcohol_percentage, volume_max, volume_remain, region, price, memo)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING id, category, name, created_at, updated_at;

-- name: UpdateSake :one
UPDATE sakes
SET category = $2,
    name = $3,
    phonetic = $4,
    alcohol_percentage = $5,
    volume_max = $6,
    volume_remain = $7,
    region = $8,
    price = $9,
    memo = $10,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING id, category, name, created_at, updated_at;

-- name: DeleteSake :execrows
DELETE FROM sakes WHERE id = $1;
