-- +goose Up
CREATE TABLE IF NOT EXISTS "Drink" (
    "Id" TEXT PRIMARY KEY,
    "Name" TEXT NOT NULL,
    "Size" INTEGER NOT NULL
);

-- +goose Down
DROP TABLE IF EXISTS "Drink";
