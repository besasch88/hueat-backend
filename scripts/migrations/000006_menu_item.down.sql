ALTER TABLE "hueat_menu_item" DROP CONSTRAINT "idx_hueat_menu_item_title";
ALTER TABLE "hueat_menu_item" DROP CONSTRAINT "fk_hueat_menu_item_menu_category";

DROP TABLE IF EXISTS "hueat_menu_item";
