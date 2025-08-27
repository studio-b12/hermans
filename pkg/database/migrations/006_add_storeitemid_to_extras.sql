-- +goose Up
ALTER TABLE "StoreItemVariant" ADD COLUMN "StoreItemId" TEXT NOT NULL DEFAULT '';
ALTER TABLE "StoreItemDip" ADD COLUMN "StoreItemId" TEXT NOT NULL DEFAULT '';

-- +goose Down
ALTER TABLE "StoreItemVariant" DROP COLUMN "StoreItemId";
ALTER TABLE "StoreItemDip" DROP COLUMN "StoreItemId";