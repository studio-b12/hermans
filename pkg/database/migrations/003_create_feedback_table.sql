-- +goose Up
CREATE TABLE "Feedback" (
    "Id"        TEXT NOT NULL PRIMARY KEY,
    "Timestamp" DATETIME NOT NULL,
    "Type"      TEXT NOT NULL,
    "Message"   TEXT NOT NULL,
    "Page"      TEXT NOT NULL
);

-- +goose Down
DROP TABLE "Feedback";