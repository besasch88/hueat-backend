ALTER TABLE "ceng_menu_category" DROP CONSTRAINT "fk_ceng_menu_category_printer_inside";
ALTER TABLE "ceng_menu_category" DROP CONSTRAINT "fk_ceng_menu_category_printer_outside";
ALTER TABLE "ceng_menu_category" DROP COLUMN "printer_outside_id";
ALTER TABLE "ceng_menu_category" RENAME "printer_inside_id" TO "printer_id";


ALTER TABLE "ceng_menu_category"
ADD CONSTRAINT "fk_ceng_menu_category_printer"
FOREIGN KEY ("printer_id")
REFERENCES "ceng_printer" ("id")
ON DELETE SET NULL
ON UPDATE CASCADE;