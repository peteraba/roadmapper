CREATE TABLE "roadmaps" (
    "id" bigint NOT NULL,
    "prev_id" bigint NULL,
    "txt" text NOT NULL,
    "created_at" timestamp NOT NULL DEFAULT NOW(),
    "updated_at" timestamp NOT NULL DEFAULT NOW()
);

CREATE OR REPLACE FUNCTION update_updated_at()
    RETURNS TRIGGER AS $$
BEGIN
    IF row(NEW.*) IS DISTINCT FROM row(OLD.*) THEN
        NEW.updated_at = now();
        RETURN NEW;
    ELSE
        RETURN OLD;
    END IF;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_roadmap_updated_at BEFORE UPDATE ON roadmaps FOR EACH ROW EXECUTE PROCEDURE update_updated_at();
