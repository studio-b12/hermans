-- +goose Up
CREATE TABLE IF NOT EXISTS "OrderItems" (
    "OrderId" TEXT NOT NULL,
    "StoreItemId" TEXT NOT NULL,
    PRIMARY KEY("OrderId", "StoreItemId"),
    FOREIGN KEY("OrderId") REFERENCES "Order"("Id")
);

-- +goose Down
DROP TABLE IF EXISTS "OrderItems";
