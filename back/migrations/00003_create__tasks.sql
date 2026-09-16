-- +goose up
CREATE TABLE
    IF NOT EXISTS "tasks" (
        "id" TEXT NOT NULL PRIMARY KEY,
        "user_id" TEXT NOT NULL,
        "project_id" TEXT NOT NULL,
        "name" TEXT NOT NULL,
        "comment" TEXT NULL,
        "created_at" INTEGER NOT NULL,
        "updated_at" INTEGER NULL,
        "archived_at" INTEGER NULL,
        "completed_at" INTEGER NULL,
        CONSTRAINT "tasks_project_id_name" UNIQUE ("project_id", "name")
    );

CREATE INDEX IF NOT EXISTS "project_id" ON "tasks" ("project_id");

-- +goose down
DROP TABLE "tasks";
