-- +goose up
CREATE TABLE
    IF NOT EXISTS "clients" (
        "id" TEXT NOT NULL PRIMARY KEY,
        "user_id" TEXT NOT NULL,
        "name" TEXT NOT NULL,
        "billable_rate" INTEGER NOT NULL,
        "comment" TEXT NULL,
        "created_at" INTEGER NOT NULL,
        "updated_at" INTEGER NULL,
        "archived_at" INTEGER NULL,
        CONSTRAINT "clients_user_id_name" UNIQUE ("user_id", "name")
    );

CREATE TABLE
    IF NOT EXISTS "projects" (
        "id" TEXT NOT NULL PRIMARY KEY,
        "user_id" TEXT NOT NULL,
        "client_id" TEXT NOT NULL,
        "name" TEXT NOT NULL,
        "billable_rate" INTEGER NOT NULL,
        "comment" TEXT NULL,
        "created_at" INTEGER NOT NULL,
        "updated_at" INTEGER NULL,
        "archived_at" INTEGER NULL,
        CONSTRAINT "projects_client_id_name" UNIQUE ("client_id", "name")
    );

CREATE INDEX IF NOT EXISTS "client_id" ON "projects" ("client_id");

CREATE TABLE
    IF NOT EXISTS "timelogs" (
        "id" TEXT NOT NULL PRIMARY KEY,
        "user_id" TEXT NOT NULL,
        "project_id" TEXT NOT NULL,
        "date" TEXT NOT NULL,
        "time_start" TEXT NOT NULL,
        "time_end" TEXT NULL,
        "duration_seconds" INTEGER NOT NULL,
        "billable_rate" INTEGER NOT NULL,
        "billable_amount" INTEGER NOT NULL,
        "comment" TEXT NULL,
        "created_at" INTEGER NOT NULL,
        "updated_at" INTEGER NULL
    );

CREATE INDEX IF NOT EXISTS "project_id" ON "timelogs" ("project_id");
CREATE INDEX IF NOT EXISTS "timelogs_user_id_date" ON "timelogs" ("user_id", "date");

CREATE TABLE
    IF NOT EXISTS "settings" (
        "user_id" TEXT NOT NULL,
        "key" TEXT NOT NULL,
        "value" TEXT NOT NULL,
        CONSTRAINT "settings_user_id_key" UNIQUE ("user_id", "key")
    );

CREATE TABLE
    IF NOT EXISTS "sessions" (
        "id" TEXT NOT NULL PRIMARY KEY,
        "cid" TEXT NOT NULL UNIQUE,
        "uid" TEXT NOT NULL,
        "checked_at" INTEGER NOT NULL
    );

CREATE TABLE
    IF NOT EXISTS "users" (
        "id" TEXT NOT NULL PRIMARY KEY,
        "email" TEXT NOT NULL,
        "created_at" INTEGER NOT NULL,
        "updated_at" INTEGER NULL,
        CONSTRAINT "users_email" UNIQUE ("email")
    );

CREATE TABLE
    IF NOT EXISTS "auth_codes" (
        "id" TEXT NOT NULL PRIMARY KEY,
        "user_id" TEXT NOT NULL,
        "user_email" TEXT NOT NULL,
        "code" TEXT NOT NULL,
        "created_at" INTEGER NOT NULL,
        CONSTRAINT "auth_codes_code" UNIQUE ("code")
    );

-- +goose down
DROP TABLE "clients";

DROP TABLE "projects";

DROP TABLE "timelogs";

DROP TABLE "settings";

DROP TABLE "sessions";

DROP TABLE "users";

DROP TABLE "auth_codes";
