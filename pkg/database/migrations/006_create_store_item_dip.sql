-- +goose Up
CREATE TABLE IF NOT EXISTS "StoreItemDip" (
    "OrderId" TEXT NOT NULL,
    "StoreItemId" TEXT NOT NULL,
    "Dip" TEXT NOT NULL,
    PRIMARY KEY("OrderId", "StoreItemId", "Dip"),
    FOREIGN KEY("OrderId") REFERENCES "Order"("Id")
);

-- +goose Down
DROP TABLE IF EXISTS "StoreItemDip";
