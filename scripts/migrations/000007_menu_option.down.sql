ALTER TABLE "hueat_menu_option" DROP CONSTRAINT "idx_hueat_menu_option_title";
ALTER TABLE "hueat_menu_option" DROP CONSTRAINT "fk_hueat_menu_option_menu_item";

DROP TABLE IF EXISTS "hueat_menu_option";
