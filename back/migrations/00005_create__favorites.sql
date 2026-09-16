-- +goose up
CREATE TABLE
    IF NOT EXISTS "favorites" (
        "id" TEXT NOT NULL PRIMARY KEY,
        "user_id" TEXT NOT NULL,
        "name" TEXT NOT NULL UNIQUE,
        "project_id" TEXT NOT NULL,
        "task_id" TEXT,
        "billable_rate" INTEGER NOT NULL,
        "comment" TEXT NULL
    );

CREATE INDEX IF NOT EXISTS "project_id" ON "favorites" ("project_id");

-- +goose down
DROP TABLE "favorites";
