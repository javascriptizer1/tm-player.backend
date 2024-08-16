-- name: UpsertCity :one
INSERT INTO "cities" (
    "id", "name"
) VALUES ($1, $2) ON CONFLICT("id") DO UPDATE SET
    "name" = EXCLUDED."name",
    "updated_at" = DEFAULT
RETURNING "id";

-- name: ListCities :many
SELECT * FROM "cities"
LIMIT sqlc.arg('limit')::bigint OFFSET sqlc.arg('offset')::bigint;

-- name: CountCities :one
SELECT COUNT(*) AS count FROM "cities";

-- name: GetCityByID :one
SELECT * FROM "cities" WHERE "id" = $1;

-- name: GetCityByNameExist :one
SELECT COUNT(*) > 0 FROM "cities" WHERE "name" = $1;