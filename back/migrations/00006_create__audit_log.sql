-- +goose up
CREATE TABLE
    IF NOT EXISTS "audit_log" (
        "id" INTEGER PRIMARY KEY AUTOINCREMENT,
        "action_at" INTEGER NOT NULL,
        "action_type" TEXT NOT NULL,
        "user_id" TEXT NOT NULL,
        "table_name" TEXT NOT NULL,
        "record_id" TEXT NOT NULL,
        "old_value" TEXT,
        "new_value" TEXT
    );

CREATE INDEX IF NOT EXISTS "idx_audit_log__action_at" ON "audit_log" ("action_at");
CREATE INDEX IF NOT EXISTS "idx_audit_log__user_id" ON "audit_log" ("user_id");
CREATE INDEX IF NOT EXISTS "idx_audit_log__record_id" ON "audit_log" ("record_id");

-- +goose down
DROP TABLE "audit_log";
