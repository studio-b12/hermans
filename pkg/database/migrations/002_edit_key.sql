-- +goose Up
ALTER TABLE "Order" ADD COLUMN "EditKey" TEXT;

-- +goose Down
ALTER TABLE "Order" DROP COLUMN "EditKey";
