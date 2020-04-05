-- +migrate Up

CREATE TABLE "roadmaps" (
    "id" bigint NOT NULL UNIQUE CONSTRAINT positive_id CHECK (id > 0),
    "prev_id" bigint NULL REFERENCES roadmaps (id),
    "date_format" text NOT NULL,
    "base_url" text NOT NULL,
    "projects" jsonb,
    "milestones" jsonb,
    "created_at" timestamp NOT NULL DEFAULT NOW(),
    "updated_at" timestamp NOT NULL DEFAULT NOW(),
    "accessed_at" timestamp NOT NULL DEFAULT NOW()
);

CREATE INDEX roadmaps_prev_id ON roadmaps (prev_id ASC);
CREATE INDEX roadmaps_accessed_at ON roadmaps (accessed_at DESC);

-- +migrate Down
