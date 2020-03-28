-- +migrate Up

CREATE TABLE "roadmaps" (
    "id" bigint NOT NULL,
    "prev_id" bigint NULL,
    "txt" text NOT NULL,
    "date_format" text NOT NULL,
    "base_url" text NOT NULL,
    "created_at" timestamp NOT NULL DEFAULT NOW(),
    "updated_at" timestamp NOT NULL DEFAULT NOW()
);

-- +migrate Down
