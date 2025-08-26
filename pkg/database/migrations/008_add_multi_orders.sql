-- +goose Up

-- 1. Neue Order-Tabelle ohne StoreItemId erstellen
CREATE TABLE "Order_new" (
    "Id" TEXT PRIMARY KEY,
    "Created" DATETIME NOT NULL,
    "Creator" TEXT NOT NULL,
    "OrderListId" TEXT NOT NULL,
    "DrinkId" TEXT,
    "EditKey" TEXT,
    FOREIGN KEY("OrderListId") REFERENCES "OrderList"("Id") ON DELETE CASCADE,
    FOREIGN KEY("DrinkId") REFERENCES "Drink"("Id") ON DELETE CASCADE
);

-- 2. Daten von alter Tabelle übertragen
INSERT INTO "Order_new" ("Id", "Created", "Creator", "OrderListId", "DrinkId", "EditKey")
SELECT "Id", "Created", "Creator", "OrderListId", "DrinkId", "EditKey"
FROM "Order";

-- 3. Alte Tabelle löschen
DROP TABLE "Order";

-- 4. Neue Tabelle umbenennen
ALTER TABLE "Order_new" RENAME TO "Order";

-- 5. Neue Tabellen für Multi-Orders erstellen
CREATE TABLE IF NOT EXISTS "OrderItems" (
    "OrderId" TEXT NOT NULL,
    "StoreItemId" TEXT NOT NULL,
    PRIMARY KEY ("OrderId", "StoreItemId"),
    FOREIGN KEY("OrderId") REFERENCES "Order"("Id") ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS "StoreItemVariant" (
    "OrderId" TEXT NOT NULL,
    "StoreItemId" TEXT NOT NULL,
    "Variant" TEXT NOT NULL,
    PRIMARY KEY ("OrderId", "StoreItemId", "Variant"),
    FOREIGN KEY("OrderId") REFERENCES "Order"("Id") ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS "StoreItemDip" (
    "OrderId" TEXT NOT NULL,
    "StoreItemId" TEXT NOT NULL,
    "Dip" TEXT NOT NULL,
    PRIMARY KEY ("OrderId", "StoreItemId", "Dip"),
    FOREIGN KEY("OrderId") REFERENCES "Order"("Id") ON DELETE CASCADE
);

-- +goose Down

DROP TABLE IF EXISTS "StoreItemDip";
DROP TABLE IF EXISTS "StoreItemVariant";
DROP TABLE IF EXISTS "OrderItems";
DROP TABLE IF EXISTS "Order";
