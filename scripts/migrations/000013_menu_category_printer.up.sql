ALTER TABLE "hueat_menu_category" DROP CONSTRAINT "fk_hueat_menu_category_printer";

ALTER TABLE "hueat_menu_category" RENAME "printer_id" TO "printer_inside_id";

ALTER TABLE "hueat_menu_category" ADD COLUMN "printer_outside_id" VARCHAR(36);

ALTER TABLE "hueat_menu_category"
ADD CONSTRAINT "fk_hueat_menu_category_printer_inside"
FOREIGN KEY ("printer_inside_id")
REFERENCES "hueat_printer" ("id")
ON DELETE SET NULL
ON UPDATE CASCADE;

ALTER TABLE "hueat_menu_category"
ADD CONSTRAINT "fk_hueat_menu_category_printer_outside"
FOREIGN KEY ("printer_outside_id")
REFERENCES "hueat_printer" ("id")
ON DELETE SET NULL
ON UPDATE CASCADE;