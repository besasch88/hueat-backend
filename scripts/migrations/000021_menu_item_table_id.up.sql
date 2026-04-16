ALTER TABLE "hueat_menu_item" ADD COLUMN "table_id" VARCHAR(36);

ALTER TABLE "hueat_menu_item"
ADD CONSTRAINT "fk_hueat_menu_item_table"
FOREIGN KEY ("table_id")
REFERENCES "hueat_table" ("id")
ON DELETE CASCADE
ON UPDATE CASCADE;
