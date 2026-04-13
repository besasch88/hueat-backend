ALTER TABLE "hueat_menu_category" DROP CONSTRAINT "idx_hueat_menu_category_title";
ALTER TABLE "hueat_menu_category" DROP CONSTRAINT "fk_hueat_menu_category_printer";

DROP TABLE IF EXISTS "hueat_menu_category";
