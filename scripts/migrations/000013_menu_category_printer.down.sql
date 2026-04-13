ALTER TABLE "hueat_menu_category" DROP CONSTRAINT "fk_hueat_menu_category_printer_inside";
ALTER TABLE "hueat_menu_category" DROP CONSTRAINT "fk_hueat_menu_category_printer_outside";
ALTER TABLE "hueat_menu_category" DROP COLUMN "printer_outside_id";
ALTER TABLE "hueat_menu_category" RENAME "printer_inside_id" TO "printer_id";


ALTER TABLE "hueat_menu_category"
ADD CONSTRAINT "fk_hueat_menu_category_printer"
FOREIGN KEY ("printer_id")
REFERENCES "hueat_printer" ("id")
ON DELETE SET NULL
ON UPDATE CASCADE;