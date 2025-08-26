-- +goose Up
CREATE TABLE IF NOT EXISTS "OrderList" (
    "Id" TEXT PRIMARY KEY,
    "Created" DATETIME NOT NULL,
    "Deadline" DATETIME
);

-- +goose Down
DROP TABLE IF EXISTS "OrderList";
