CREATE TABLE "hueat_config" (
    "config_key" VARCHAR(255) PRIMARY KEY NOT NULL,
    "config_value" TEXT NOT NULL,
    "created_at" TIMESTAMP NOT NULL,
    "updated_at" TIMESTAMP NOT NULL
);

ALTER TABLE "hueat_config" ADD CONSTRAINT "idx_hueat_config_config_key" UNIQUE ("config_key");

INSERT INTO "hueat_config" ("config_key", "config_value", "created_at", "updated_at")
VALUES ('pricedOrderPrinterInsideID', '', NOW(), NOW());

INSERT INTO "hueat_config" ("config_key", "config_value", "created_at", "updated_at")
VALUES ('pricedOrderPrinterOutsideID', '', NOW(), NOW());

INSERT INTO "hueat_config" ("config_key", "config_value", "created_at", "updated_at")
VALUES ('pricedOrderPrinterTitle', '', NOW(), NOW());

INSERT INTO "hueat_config" ("config_key", "config_value", "created_at", "updated_at")
VALUES ('totalPricePaymentPrinterInsideID', '', NOW(), NOW());

INSERT INTO "hueat_config" ("config_key", "config_value", "created_at", "updated_at")
VALUES ('totalPricePaymentPrinterOutsideID', '', NOW(), NOW());

INSERT INTO "hueat_config" ("config_key", "config_value", "created_at", "updated_at")
VALUES ('totalPricePaymentPrinterTitle', '', NOW(), NOW());