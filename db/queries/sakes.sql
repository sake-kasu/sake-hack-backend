-- name: ListSakes :many
SELECT id, name, category, subcategory, alcohol_percentage, volume, origin, price, stock, created_at, updated_at
FROM sakes
WHERE
    (sqlc.narg('category')::sake_category IS NULL OR category = sqlc.narg('category'))
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: CountSakes :one
SELECT COUNT(*) AS total
FROM sakes
WHERE
    (sqlc.narg('category')::sake_category IS NULL OR category = sqlc.narg('category'));

-- name: GetSakeByID :one
SELECT id, name, category, subcategory, alcohol_percentage, volume, origin, price, stock, created_at, updated_at
FROM sakes
WHERE id = $1;

-- name: CreateSake :one
INSERT INTO sakes (name, category, subcategory, alcohol_percentage, volume, origin, price, stock)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING id, created_at, updated_at;

-- name: UpdateSake :exec
UPDATE sakes
SET name = $2, category = $3, subcategory = $4, alcohol_percentage = $5,
    volume = $6, origin = $7, price = $8, stock = $9, updated_at = CURRENT_TIMESTAMP
WHERE id = $1;

-- name: DeleteSake :exec
DELETE FROM sakes WHERE id = $1;
