-- name: ListKinds :many
SELECT id, name FROM sake_kinds ORDER BY id;

-- name: ListBreweries :many
SELECT id, name, origin_country, origin_region, latitude, longitude
FROM breweries
WHERE (sqlc.narg('keyword')::TEXT IS NULL OR name LIKE '%' || sqlc.narg('keyword') || '%')
ORDER BY id
LIMIT $1;

-- name: ListDrinkStyles :many
SELECT id, name, description FROM drink_styles ORDER BY id;
