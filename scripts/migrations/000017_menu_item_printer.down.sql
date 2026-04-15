ALTER TABLE "hueat_menu_item" DROP CONSTRAINT "fk_hueat_menu_item_printer_inside";
ALTER TABLE "hueat_menu_item" DROP CONSTRAINT "fk_hueat_menu_item_printer_outside";


ALTER TABLE "hueat_menu_item"
DROP COLUMN "printer_inside_id",
DROP COLUMN "printer_outside_id";