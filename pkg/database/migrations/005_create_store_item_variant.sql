-- +goose Up
CREATE TABLE IF NOT EXISTS "StoreItemVariant" (
    "OrderId" TEXT NOT NULL,
    "StoreItemId" TEXT NOT NULL,
    "Variant" TEXT NOT NULL,
    PRIMARY KEY("OrderId", "StoreItemId", "Variant"),
    FOREIGN KEY("OrderId") REFERENCES "Order"("Id")
);

-- +goose Down
DROP TABLE IF EXISTS "StoreItemVariant";
