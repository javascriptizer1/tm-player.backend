-- +goose Up
-- +goose StatementBegin
CREATE TABLE "cities" (
    "id" uuid PRIMARY KEY,
    "name" varchar NOT NULL,
    "created_at" timestamp with time zone NOT NULL DEFAULT now(),
    "updated_at" timestamp with time zone NOT NULL DEFAULT now()
);

CREATE TABLE "positions" (
    "id" uuid PRIMARY KEY,
    "name" varchar NOT NULL,
    "created_at" timestamp with time zone NOT NULL DEFAULT now(),
    "updated_at" timestamp with time zone NOT NULL DEFAULT now()
);

CREATE TABLE "players" (
    "id" uuid PRIMARY KEY,
    "first_name" varchar NOT NULL,
    "last_name" varchar NOT NULL,
    "middle_name" varchar,
    "birthday" timestamp with time zone NOT NULL,
    "photo" varchar,
    "height" integer NOT NULL DEFAULT 0,
    "impact_leg" varchar NOT NULL,
    "market_value" integer NOT NULL DEFAULT 0,
    "city_id" uuid NOT NULL,
    "created_at" timestamp with time zone NOT NULL DEFAULT now(),
    "updated_at" timestamp with time zone NOT NULL DEFAULT now(),
    FOREIGN KEY ("city_id") REFERENCES "cities" ("id") ON DELETE SET NULL
);

CREATE TABLE "player_positions" (
    "player_id" uuid NOT NULL,
    "position_id" uuid NOT NULL,
    "main" boolean NOT NULL,
    "created_at" timestamp with time zone NOT NULL DEFAULT now(),
    "updated_at" timestamp with time zone NOT NULL DEFAULT now(),
    PRIMARY KEY ("player_id", "position_id"),
    FOREIGN KEY ("player_id") REFERENCES "players" ("id") ON DELETE CASCADE,
    FOREIGN KEY ("position_id") REFERENCES "positions" ("id") ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE "player_positions";
DROP TABLE "players";
DROP TABLE "positions";
DROP TABLE "cities";
-- +goose StatementEnd
