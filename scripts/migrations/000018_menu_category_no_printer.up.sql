ALTER TABLE "hueat_menu_category" DROP CONSTRAINT "fk_hueat_menu_category_printer_inside";
ALTER TABLE "hueat_menu_category" DROP CONSTRAINT "fk_hueat_menu_category_printer_outside";
ALTER TABLE "hueat_menu_category" DROP COLUMN "printer_inside_id";
ALTER TABLE "hueat_menu_category" DROP COLUMN "printer_outside_id";