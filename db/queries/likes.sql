-- name: CreateLike :exec
INSERT INTO likes (sake_id, token) VALUES ($1, $2)
ON CONFLICT (sake_id, token) DO NOTHING;

-- name: DeleteLike :execrows
DELETE FROM likes WHERE sake_id = $1 AND token = $2;

-- name: GetLikeCount :one
SELECT COUNT(*) AS like_count FROM likes WHERE sake_id = $1;

-- name: ExistsLike :one
SELECT EXISTS(SELECT 1 FROM likes WHERE sake_id = $1 AND token = $2) AS is_liked;
