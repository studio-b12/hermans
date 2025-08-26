-- +goose Up
CREATE TABLE IF NOT EXISTS "Order" (
    "Id" TEXT PRIMARY KEY,
    "Created" DATETIME NOT NULL,
    "Creator" TEXT NOT NULL,
    "OrderListId" TEXT NOT NULL,
    "DrinkId" TEXT,
    "EditKey" TEXT,
    FOREIGN KEY("OrderListId") REFERENCES "OrderList"("Id"),
    FOREIGN KEY("DrinkId") REFERENCES "Drink"("Id")
);

-- +goose Down
DROP TABLE IF EXISTS "Order";
