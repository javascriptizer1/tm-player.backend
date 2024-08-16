-- name: UpsertPlayer :one
INSERT INTO "players" (
    "id",
    "first_name",
    "last_name",
    "middle_name",
    "birthday",
    "photo",
    "city_id",
    "height",
    "impact_leg",
    "market_value"
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10
) ON CONFLICT("id") DO UPDATE SET
    "first_name" = EXCLUDED."first_name",
    "last_name" = EXCLUDED."last_name",
    "middle_name" = EXCLUDED."middle_name",
    "birthday" = EXCLUDED."birthday",
    "photo" = EXCLUDED."photo",
    "city_id" = EXCLUDED."city_id",
    "height" = EXCLUDED."height",
    "impact_leg" = EXCLUDED."impact_leg",
    "market_value" = EXCLUDED."market_value",
    "updated_at" = DEFAULT
RETURNING "id";

-- name: ListPlayers :many
SELECT p.*
FROM "players" p
LEFT JOIN 
    "player_positions" pp ON p.id = pp.player_id
WHERE 
    (sqlc.narg('city_id')::uuid IS NULL OR p.city_id = sqlc.narg('city_id'))
    AND (sqlc.narg('position_id')::uuid IS NULL OR pp.position_id = sqlc.narg('position_id'))
    AND (sqlc.narg('max_age')::bigint IS NULL OR p.birthday <= NOW() - INTERVAL '1 year' * sqlc.narg('max_age'))
    AND (sqlc.narg('min_age')::bigint IS NULL OR p.birthday >= NOW() - INTERVAL '1 year' * sqlc.narg('min_age'))
    AND (COALESCE(@ids::uuid[], '{}') = '{}' OR p.id = ANY(@ids::uuid[]))
LIMIT sqlc.arg('limit')::bigint OFFSET sqlc.arg('offset')::bigint;


-- name: CountPlayers :one
SELECT COUNT(p.*)
FROM "players" p
LEFT JOIN 
    "player_positions" pp ON p.id = pp.player_id
WHERE 
    (sqlc.narg('city_id')::uuid IS NULL OR p.city_id = sqlc.narg('city_id'))
    AND (sqlc.narg('position_id')::uuid IS NULL OR pp.position_id = sqlc.narg('position_id'))
    AND (sqlc.narg('max_age')::bigint IS NULL OR p.birthday <= NOW() - INTERVAL '1 year' * sqlc.narg('max_age'))
    AND (sqlc.narg('min_age')::bigint IS NULL OR p.birthday >= NOW() - INTERVAL '1 year' * sqlc.narg('min_age'))
    AND (COALESCE(@ids::uuid[], '{}') = '{}' OR p.id = ANY(@ids::uuid[]));

-- name: GetPlayerByID :one
SELECT p.*
FROM 
    "players" p
LEFT JOIN 
    "player_positions" pp ON p.id = pp.player_id
WHERE 
    p.id = $1;

-- name: UpsertPlayerPositions :exec
INSERT INTO "player_positions" ("player_id", "position_id", "main") 
SELECT $1, unnest(@position_ids::uuid[]), unnest(@mains::boolean[])
ON CONFLICT ("player_id", "position_id") DO NOTHING;

-- name: TrimNotExistingPlayerPositions :many
DELETE FROM "player_positions"
WHERE "player_id" = $1 AND "position_id" NOT IN (SELECT unnest(@existing_position_ids::uuid[]))
RETURNING *;

-- name: GetPositionsForPlayers :many
SELECT *
FROM "player_positions"
WHERE "player_id" = ANY(@player_ids::uuid[]);