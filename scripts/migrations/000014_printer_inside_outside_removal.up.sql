ALTER TABLE "hueat_printer"
DROP COLUMN "inside",
DROP COLUMN "outside";

ALTER TABLE "hueat_printer" ADD CONSTRAINT "idx_hueat_printer_title" UNIQUE ("title");