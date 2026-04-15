ALTER TABLE "hueat_menu_item" ADD COLUMN "printer_inside_id" VARCHAR(36);
ALTER TABLE "hueat_menu_item" ADD COLUMN "printer_outside_id" VARCHAR(36);


ALTER TABLE "hueat_menu_item"
ADD CONSTRAINT "fk_hueat_menu_item_printer_inside"
FOREIGN KEY ("printer_inside_id")
REFERENCES "hueat_printer" ("id")
ON DELETE SET NULL
ON UPDATE CASCADE;

ALTER TABLE "hueat_menu_item"
ADD CONSTRAINT "fk_hueat_menu_item_printer_outside"
FOREIGN KEY ("printer_outside_id")
REFERENCES "hueat_printer" ("id")
ON DELETE SET NULL
ON UPDATE CASCADE;