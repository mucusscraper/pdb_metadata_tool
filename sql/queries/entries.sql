-- name: CreateEntry :one
INSERT INTO entries (id, created_at, updated_at, rcsb_id,deposit_date,doi,paper_title,method,user_group)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7,
    $8,
    $9
)
RETURNING *;

-- name: GetEntry :one
SELECT * FROM entries WHERE rcsb_id=$1;
