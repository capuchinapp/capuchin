-- +goose up
ALTER TABLE "users"
ADD COLUMN "activate_code" TEXT;

-- +goose down
ALTER TABLE "users"
DROP COLUMN "activate_code";
