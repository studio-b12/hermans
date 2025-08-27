-- +goose Up
CREATE TABLE "OrderItems" (
    "OrderId" TEXT NOT NULL,
    "StoreItemId" TEXT NOT NULL,
    FOREIGN KEY(OrderId) REFERENCES "Order"(Id) ON DELETE CASCADE
);

ALTER TABLE "Order" DROP COLUMN "StoreItemId";

-- +goose Down
ALTER TABLE "Order" ADD COLUMN "StoreItemId" TEXT;
DROP TABLE "OrderItems";