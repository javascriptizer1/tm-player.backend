-- name: UpsertPosition :one
INSERT INTO "positions" (
    "id", "name"
) VALUES ($1, $2) ON CONFLICT("id") DO UPDATE SET
    "name" = EXCLUDED."name",
    "updated_at" = DEFAULT
RETURNING "id";

-- name: ListPositions :many
SELECT * FROM "positions"
WHERE (@ids::uuid[] = '{}' OR id = ANY(@ids::uuid[]))
LIMIT sqlc.arg('limit')::bigint OFFSET sqlc.arg('offset')::bigint;

-- name: CountPositions :one
SELECT COUNT(*) AS count FROM "positions"
WHERE (@ids::uuid[] = '{}' OR id = ANY(@ids::uuid[]));

-- name: GetPositionByID :one
SELECT * FROM "positions" WHERE "id" = $1;

-- name: GetPositionByNameExist :one
SELECT COUNT(*) > 0 FROM "positions" WHERE "name" = $1;