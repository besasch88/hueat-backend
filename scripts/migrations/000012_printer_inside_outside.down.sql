ALTER TABLE "ceng_printer"
DROP COLUMN "inside",
DROP COLUMN "outside";

ALTER TABLE "ceng_printer" ADD CONSTRAINT "idx_ceng_printer_title" UNIQUE ("title");