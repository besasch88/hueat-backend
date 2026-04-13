CREATE TABLE "hueat_menu_category" (
    "id" VARCHAR(36) PRIMARY KEY NOT NULL,
    "title" VARCHAR(255) NOT NULL,
    "position" BIGINT NOT NULL,
    "active" BOOLEAN NOT NULL,
    "printer_id" VARCHAR(36),
    "created_at" TIMESTAMP NOT NULL,
    "updated_at" TIMESTAMP NOT NULL
);

ALTER TABLE "hueat_menu_category"
ADD CONSTRAINT "fk_hueat_menu_category_printer"
FOREIGN KEY ("printer_id")
REFERENCES "hueat_printer" ("id")
ON DELETE SET NULL
ON UPDATE CASCADE;

ALTER TABLE "hueat_menu_category" ADD CONSTRAINT "idx_hueat_menu_category_title" UNIQUE ("title");
