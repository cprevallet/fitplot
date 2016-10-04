
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE IF NOT EXISTS "runfiles" ("id" integer not null primary key, "filename" text, "filetype" text,"filecontent" blob, "timestamp" datetime );

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE "runfiles";

