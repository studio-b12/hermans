-- +goose Up
ALTER TABLE "OrderList" ADD COLUMN "Deadline" DATETIME NULL;

-- +goose Down
ALTER TABLE "OrderList" DROP COLUMN "Deadline";