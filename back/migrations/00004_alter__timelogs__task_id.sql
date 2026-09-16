-- +goose up
ALTER TABLE "timelogs"
ADD COLUMN "task_id" TEXT;

-- +goose down
ALTER TABLE "timelogs"
DROP COLUMN "task_id";
