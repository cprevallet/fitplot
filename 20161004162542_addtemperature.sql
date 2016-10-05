
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
ALTER TABLE "runfiles" ADD COLUMN "temperature" REAL DEFAULT 0.0;

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
	ALTER TABLE "runfiles" RENAME TO "temp_runfiles";
	
	CREATE TABLE IF NOT EXISTS "runfiles" (
		"id" integer not null primary key, 
		"filename" text, 
		"filetype" text,
		"filecontent" blob, 
		"timestamp" datetime );
	
	INSERT INTO "runfiles" 
	SELECT
	"id", "filename", "filetype", "filecontent", "timestamp"
	FROM
	"temp_runfiles";
	
	DROP TABLE "temp_runfiles";
	