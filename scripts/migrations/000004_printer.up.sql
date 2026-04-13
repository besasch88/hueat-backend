CREATE TABLE "hueat_printer" (
    "id" VARCHAR(36) PRIMARY KEY NOT NULL,
    "title" VARCHAR(255) NOT NULL,
    "url" VARCHAR(255) NOT NULL,
    "active" BOOLEAN NOT NULL,
    "created_at" TIMESTAMP NOT NULL,
    "updated_at" TIMESTAMP NOT NULL
);

ALTER TABLE "hueat_printer" ADD CONSTRAINT "idx_hueat_printer_title" UNIQUE ("title");