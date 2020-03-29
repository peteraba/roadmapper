-- +migrate Up

CREATE TABLE "roadmaps" (
    "id" bigint NOT NULL,
    "prev_id" bigint NULL,
    "txt" text NOT NULL DEFAULT '',
    "date_format" text NOT NULL DEFAULT '',
    "base_url" text NOT NULL DEFAULT '',
    "created_at" timestamp NOT NULL DEFAULT NOW(),
    "updated_at" timestamp NOT NULL DEFAULT NOW(),
    "accessed_at" timestamp NOT NULL DEFAULT NOW()
);

-- +migrate Down
